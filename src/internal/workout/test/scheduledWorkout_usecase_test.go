package test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/alecthomas/assert/v2"
	"github.com/alielmi98/go-hexa-workout/internal/workout/core/models"
	"github.com/alielmi98/go-hexa-workout/internal/workout/core/usecase/dto"
	"github.com/alielmi98/go-hexa-workout/pkg/service_errors"
)

func TestCreateScheduledWorkout_Success(t *testing.T) {
	scheduledRepo := &MockScheduledWorkoutsRepository{}
	workoutRepo := &MockWorkoutRepository{}
	useCase := setupScheduledWorkoutUsecase(scheduledRepo, workoutRepo)

	ctx := createContextWithUserId(1)
	req := dto.CreateScheduledWorkoutsRequest{
		WorkoutId:     1,
		ScheduledTime: time.Now(),
		Status:        "active",
	}

	response, err := useCase.Create(ctx, req)

	assert.NoError(t, err)
	assert.Equal(t, 1, response.Id)
	assert.Equal(t, 1, response.WorkoutId)
	assert.Equal(t, "active", response.Status)
}

func TestCreateScheduledWorkout_InvalidStatus(t *testing.T) {
	scheduledRepo := &MockScheduledWorkoutsRepository{}
	workoutRepo := &MockWorkoutRepository{}
	useCase := setupScheduledWorkoutUsecase(scheduledRepo, workoutRepo)

	ctx := createContextWithUserId(1)
	req := dto.CreateScheduledWorkoutsRequest{
		WorkoutId:     1,
		ScheduledTime: time.Now(),
		Status:        "invalid_status",
	}

	_, err := useCase.Create(ctx, req)

	assert.Error(t, err)
	serviceErr, ok := err.(*service_errors.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, service_errors.InvalidStatus, serviceErr.EndUserMessage)
}

func TestCreateScheduledWorkout_WorkoutNotFound(t *testing.T) {
	scheduledRepo := &MockScheduledWorkoutsRepository{}
	workoutRepo := &MockWorkoutRepository{
		GetByIdFn: func(ctx context.Context, id int) (models.Workout, error) {
			return models.Workout{}, errors.New("workout not found")
		},
	}
	useCase := setupScheduledWorkoutUsecase(scheduledRepo, workoutRepo)

	ctx := createContextWithUserId(1)
	req := dto.CreateScheduledWorkoutsRequest{
		WorkoutId:     999,
		ScheduledTime: time.Now(),
		Status:        "active",
	}

	_, err := useCase.Create(ctx, req)

	assert.Error(t, err)
}

func TestCreateScheduledWorkout_UnauthorizedUser(t *testing.T) {
	scheduledRepo := &MockScheduledWorkoutsRepository{}
	workoutRepo := &MockWorkoutRepository{
		GetByIdFn: func(ctx context.Context, id int) (models.Workout, error) {
			return models.Workout{
				Id:     1,
				UserId: 2, // Different user ID
				Name:   "Test Workout",
			}, nil
		},
	}
	useCase := setupScheduledWorkoutUsecase(scheduledRepo, workoutRepo)

	ctx := createContextWithUserId(1) // User ID 1 trying to access workout owned by user ID 2
	req := dto.CreateScheduledWorkoutsRequest{
		WorkoutId:     1,
		ScheduledTime: time.Now(),
		Status:        "active",
	}

	_, err := useCase.Create(ctx, req)

	assert.Error(t, err)
}

func TestCreateScheduledWorkout_RepositoryError(t *testing.T) {
	scheduledRepo := &MockScheduledWorkoutsRepository{
		CreateFn: func(ctx context.Context, entity models.ScheduledWorkouts) (models.ScheduledWorkouts, error) {
			return models.ScheduledWorkouts{}, errors.New("database error")
		},
	}
	workoutRepo := &MockWorkoutRepository{}
	useCase := setupScheduledWorkoutUsecase(scheduledRepo, workoutRepo)

	ctx := createContextWithUserId(1)
	req := dto.CreateScheduledWorkoutsRequest{
		WorkoutId:     1,
		ScheduledTime: time.Now(),
		Status:        "active",
	}

	_, err := useCase.Create(ctx, req)

	assert.Error(t, err)
	assert.Equal(t, "database error", err.Error())
}

func TestUpdateScheduledWorkout_Success(t *testing.T) {
	scheduledRepo := &MockScheduledWorkoutsRepository{}
	workoutRepo := &MockWorkoutRepository{}
	useCase := setupScheduledWorkoutUsecase(scheduledRepo, workoutRepo)

	ctx := createContextWithUserId(1)
	req := dto.UpdateScheduledWorkoutsRequest{
		ScheduledTime: time.Now(),
		Status:        "cancelled",
	}

	response, err := useCase.Update(ctx, 1, req)

	assert.NoError(t, err)
	assert.Equal(t, 1, response.Id)
	assert.Equal(t, "cancelled", response.Status)
}

func TestUpdateScheduledWorkout_InvalidStatus(t *testing.T) {
	scheduledRepo := &MockScheduledWorkoutsRepository{}
	workoutRepo := &MockWorkoutRepository{}
	useCase := setupScheduledWorkoutUsecase(scheduledRepo, workoutRepo)

	ctx := createContextWithUserId(1)
	req := dto.UpdateScheduledWorkoutsRequest{
		ScheduledTime: time.Now(),
		Status:        "invalid_status",
	}

	_, err := useCase.Update(ctx, 1, req)

	assert.Error(t, err)
	serviceErr, ok := err.(*service_errors.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, service_errors.InvalidStatus, serviceErr.EndUserMessage)
}

func TestUpdateScheduledWorkout_ScheduledWorkoutNotFound(t *testing.T) {
	scheduledRepo := &MockScheduledWorkoutsRepository{
		GetByIdFn: func(ctx context.Context, id int) (models.ScheduledWorkouts, error) {
			return models.ScheduledWorkouts{}, errors.New("scheduled workout not found")
		},
	}
	workoutRepo := &MockWorkoutRepository{}
	useCase := setupScheduledWorkoutUsecase(scheduledRepo, workoutRepo)

	ctx := createContextWithUserId(1)
	req := dto.UpdateScheduledWorkoutsRequest{
		ScheduledTime: time.Now(),
		Status:        "completed",
	}

	_, err := useCase.Update(ctx, 999, req)

	assert.Error(t, err)
}

func TestUpdateScheduledWorkout_UnauthorizedUser(t *testing.T) {
	scheduledRepo := &MockScheduledWorkoutsRepository{}
	workoutRepo := &MockWorkoutRepository{
		GetByIdFn: func(ctx context.Context, id int) (models.Workout, error) {
			return models.Workout{
				Id:     1,
				UserId: 2, // Different user ID
				Name:   "Test Workout",
			}, nil
		},
	}
	useCase := setupScheduledWorkoutUsecase(scheduledRepo, workoutRepo)

	ctx := createContextWithUserId(1) // User ID 1 trying to update workout owned by user ID 2
	req := dto.UpdateScheduledWorkoutsRequest{
		ScheduledTime: time.Now(),
		Status:        "completed",
	}

	_, err := useCase.Update(ctx, 1, req)

	assert.Error(t, err)
}

func TestDeleteScheduledWorkout_Success(t *testing.T) {
	scheduledRepo := &MockScheduledWorkoutsRepository{}
	workoutRepo := &MockWorkoutRepository{}
	useCase := setupScheduledWorkoutUsecase(scheduledRepo, workoutRepo)

	ctx := createContextWithUserId(1)

	err := useCase.Delete(ctx, 1)

	assert.NoError(t, err)
}

func TestDeleteScheduledWorkout_ScheduledWorkoutNotFound(t *testing.T) {
	scheduledRepo := &MockScheduledWorkoutsRepository{
		GetByIdFn: func(ctx context.Context, id int) (models.ScheduledWorkouts, error) {
			return models.ScheduledWorkouts{}, errors.New("scheduled workout not found")
		},
	}
	workoutRepo := &MockWorkoutRepository{}
	useCase := setupScheduledWorkoutUsecase(scheduledRepo, workoutRepo)

	ctx := createContextWithUserId(1)

	err := useCase.Delete(ctx, 999)

	assert.Error(t, err)
}

func TestDeleteScheduledWorkout_UnauthorizedUser(t *testing.T) {
	scheduledRepo := &MockScheduledWorkoutsRepository{}
	workoutRepo := &MockWorkoutRepository{
		GetByIdFn: func(ctx context.Context, id int) (models.Workout, error) {
			return models.Workout{
				Id:     1,
				UserId: 2, // Different user ID
				Name:   "Test Workout",
			}, nil
		},
	}
	useCase := setupScheduledWorkoutUsecase(scheduledRepo, workoutRepo)

	ctx := createContextWithUserId(1) // User ID 1 trying to delete workout owned by user ID 2

	err := useCase.Delete(ctx, 1)

	assert.Error(t, err)
}

func TestGetScheduledWorkoutById_Success(t *testing.T) {
	scheduledRepo := &MockScheduledWorkoutsRepository{}
	workoutRepo := &MockWorkoutRepository{}
	useCase := setupScheduledWorkoutUsecase(scheduledRepo, workoutRepo)

	ctx := createContextWithUserId(1)

	response, err := useCase.GetById(ctx, 1)

	assert.NoError(t, err)
	assert.Equal(t, 1, response.Id)
	assert.Equal(t, 1, response.WorkoutId)
	assert.Equal(t, "active", response.Status)
}

func TestGetScheduledWorkoutById_NotFound(t *testing.T) {
	scheduledRepo := &MockScheduledWorkoutsRepository{
		GetByIdFn: func(ctx context.Context, id int) (models.ScheduledWorkouts, error) {
			return models.ScheduledWorkouts{}, errors.New("scheduled workout not found")
		},
	}
	workoutRepo := &MockWorkoutRepository{}
	useCase := setupScheduledWorkoutUsecase(scheduledRepo, workoutRepo)

	ctx := createContextWithUserId(1)

	_, err := useCase.GetById(ctx, 999)

	assert.Error(t, err)
}

func TestGetScheduledWorkoutById_UnauthorizedUser(t *testing.T) {
	scheduledRepo := &MockScheduledWorkoutsRepository{}
	workoutRepo := &MockWorkoutRepository{
		GetByIdFn: func(ctx context.Context, id int) (models.Workout, error) {
			return models.Workout{
				Id:     1,
				UserId: 2, // Different user ID
				Name:   "Test Workout",
			}, nil
		},
	}
	useCase := setupScheduledWorkoutUsecase(scheduledRepo, workoutRepo)

	ctx := createContextWithUserId(1) // User ID 1 trying to access workout owned by user ID 2

	_, err := useCase.GetById(ctx, 1)

	assert.Error(t, err)
	assert.Equal(t, "user is not the owner of this workout", err.Error())

}
