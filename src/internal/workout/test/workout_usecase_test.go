package test

import (
	"context"
	"errors"
	"testing"

	"github.com/alecthomas/assert/v2"
	"github.com/alielmi98/go-hexa-workout/internal/workout/core/models"
	"github.com/alielmi98/go-hexa-workout/internal/workout/core/usecase/dto"
)

// ==================== WORKOUT USECASE TESTS ====================

func TestCreateWorkout_Success(t *testing.T) {
	workoutRepo := &MockWorkoutRepository{}
	useCase := setupWorkoutUsecase(workoutRepo)

	ctx := createContextWithUserId(1)
	req := dto.CreateWorkoutRequest{
		Name:        "Test Workout",
		Description: "Test Description",
		Comments:    "Test Comments",
		UserId:      1,
	}

	response, err := useCase.Create(ctx, req)

	assert.NoError(t, err)
	assert.Equal(t, 1, response.Id)
	assert.Equal(t, "Test Workout", response.Name)
	assert.Equal(t, "Test Description", response.Description)
}

func TestCreateWorkout_RepositoryError(t *testing.T) {
	workoutRepo := &MockWorkoutRepository{
		CreateFn: func(ctx context.Context, entity models.Workout) (models.Workout, error) {
			return models.Workout{}, errors.New("database error")
		},
	}
	useCase := setupWorkoutUsecase(workoutRepo)

	ctx := createContextWithUserId(1)
	req := dto.CreateWorkoutRequest{
		Name:        "Test Workout",
		Description: "Test Description",
		Comments:    "Test Comments",
		UserId:      1,
	}

	_, err := useCase.Create(ctx, req)

	assert.Error(t, err)
	assert.Equal(t, "database error", err.Error())
}

func TestUpdateWorkout_Success(t *testing.T) {
	workoutRepo := &MockWorkoutRepository{}
	useCase := setupWorkoutUsecase(workoutRepo)

	ctx := createContextWithUserId(1)
	req := dto.UpdateWorkoutRequest{
		Name:        "Updated Workout",
		Description: "Updated Description",
		Comments:    "Updated Comments",
	}

	response, err := useCase.Update(ctx, 1, req)

	assert.NoError(t, err)
	assert.Equal(t, 1, response.Id)
	assert.Equal(t, "Updated Workout", response.Name)
}

func TestUpdateWorkout_UnauthorizedUser(t *testing.T) {
	workoutRepo := &MockWorkoutRepository{
		GetByIdFn: func(ctx context.Context, id int) (models.Workout, error) {
			return models.Workout{
				Id:     1,
				UserId: 2, // Different user ID
				Name:   "Test Workout",
			}, nil
		},
	}
	useCase := setupWorkoutUsecase(workoutRepo)

	ctx := createContextWithUserId(1) // User ID 1 trying to update workout owned by user ID 2
	req := dto.UpdateWorkoutRequest{
		Name:        "Updated Workout",
		Description: "Updated Description",
		Comments:    "Updated Comments",
	}

	_, err := useCase.Update(ctx, 1, req)

	assert.Error(t, err)
}

func TestDeleteWorkout_Success(t *testing.T) {
	workoutRepo := &MockWorkoutRepository{}
	useCase := setupWorkoutUsecase(workoutRepo)

	ctx := createContextWithUserId(1)

	err := useCase.Delete(ctx, 1)

	assert.NoError(t, err)
}

func TestDeleteWorkout_UnauthorizedUser(t *testing.T) {
	workoutRepo := &MockWorkoutRepository{
		GetByIdFn: func(ctx context.Context, id int) (models.Workout, error) {
			return models.Workout{
				Id:     1,
				UserId: 2, // Different user ID
				Name:   "Test Workout",
			}, nil
		},
	}
	useCase := setupWorkoutUsecase(workoutRepo)

	ctx := createContextWithUserId(1) // User ID 1 trying to delete workout owned by user ID 2

	err := useCase.Delete(ctx, 1)

	assert.Error(t, err)
}

func TestGetWorkoutById_Success(t *testing.T) {
	workoutRepo := &MockWorkoutRepository{}
	useCase := setupWorkoutUsecase(workoutRepo)

	ctx := createContextWithUserId(1)

	response, err := useCase.GetById(ctx, 1)

	assert.NoError(t, err)
	assert.Equal(t, 1, response.Id)
	assert.Equal(t, "Test Workout", response.Name)
}

func TestGetWorkoutById_UnauthorizedUser(t *testing.T) {
	workoutRepo := &MockWorkoutRepository{
		GetByIdFn: func(ctx context.Context, id int) (models.Workout, error) {
			return models.Workout{
				Id:     1,
				UserId: 2, // Different user ID
				Name:   "Test Workout",
			}, nil
		},
	}
	useCase := setupWorkoutUsecase(workoutRepo)

	ctx := createContextWithUserId(1) // User ID 1 trying to access workout owned by user ID 2

	_, err := useCase.GetById(ctx, 1)

	assert.Error(t, err)
}
