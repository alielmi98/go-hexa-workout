package test

import (
	"context"
	"errors"
	"testing"

	"github.com/alecthomas/assert/v2"
	"github.com/alielmi98/go-hexa-workout/internal/workout/core/models"
	"github.com/alielmi98/go-hexa-workout/internal/workout/core/usecase/dto"
)

// ==================== WORKOUT EXERCISE USECASE TESTS ====================

func TestCreateWorkoutExercise_Success(t *testing.T) {
	exerciseRepo := &MockWorkoutExerciseRepository{}
	workoutRepo := &MockWorkoutRepository{}
	useCase := setupWorkoutExerciseUsecase(exerciseRepo, workoutRepo)

	ctx := createContextWithUserId(1)
	req := dto.CreateWorkoutExerciseRequest{
		WorkoutId:   1,
		Name:        "Push Up",
		Description: "Standard push up exercise",
		Repetitions: 15,
		Sets:        3,
		Weight:      0.0,
	}

	response, err := useCase.Create(ctx, req)

	assert.NoError(t, err)
	assert.Equal(t, 1, response.Id)
	assert.Equal(t, "Push Up", response.Name)
	assert.Equal(t, 15, response.Repetitions)
}

func TestCreateWorkoutExercise_UnauthorizedUser(t *testing.T) {
	exerciseRepo := &MockWorkoutExerciseRepository{}
	workoutRepo := &MockWorkoutRepository{
		GetByIdFn: func(ctx context.Context, id int) (models.Workout, error) {
			return models.Workout{
				Id:     1,
				UserId: 2, // Different user ID
				Name:   "Test Workout",
			}, nil
		},
	}
	useCase := setupWorkoutExerciseUsecase(exerciseRepo, workoutRepo)

	ctx := createContextWithUserId(1) // User ID 1 trying to access workout owned by user ID 2
	req := dto.CreateWorkoutExerciseRequest{
		WorkoutId:   1,
		Name:        "Push Up",
		Description: "Standard push up exercise",
		Repetitions: 15,
		Sets:        3,
		Weight:      0.0,
	}

	_, err := useCase.Create(ctx, req)

	assert.Error(t, err)
}

func TestCreateWorkoutExercise_RepositoryError(t *testing.T) {
	exerciseRepo := &MockWorkoutExerciseRepository{
		CreateFn: func(ctx context.Context, entity models.WorkoutExercise) (models.WorkoutExercise, error) {
			return models.WorkoutExercise{}, errors.New("database error")
		},
	}
	workoutRepo := &MockWorkoutRepository{}
	useCase := setupWorkoutExerciseUsecase(exerciseRepo, workoutRepo)

	ctx := createContextWithUserId(1)
	req := dto.CreateWorkoutExerciseRequest{
		WorkoutId:   1,
		Name:        "Push Up",
		Description: "Standard push up exercise",
		Repetitions: 15,
		Sets:        3,
		Weight:      0.0,
	}

	_, err := useCase.Create(ctx, req)

	assert.Error(t, err)
	assert.Equal(t, "database error", err.Error())
}

func TestUpdateWorkoutExercise_Success(t *testing.T) {
	exerciseRepo := &MockWorkoutExerciseRepository{}
	workoutRepo := &MockWorkoutRepository{}
	useCase := setupWorkoutExerciseUsecase(exerciseRepo, workoutRepo)

	ctx := createContextWithUserId(1)
	req := dto.UpdateWorkoutExerciseRequest{
		Name:        "Updated Push Up",
		WorkoutId:   1,
		Description: "Updated description",
		Repetitions: 20,
		Sets:        4,
		Weight:      5.0,
	}

	response, err := useCase.Update(ctx, 1, req)

	assert.NoError(t, err)
	assert.Equal(t, 1, response.Id)
	assert.Equal(t, "Updated Push Up", response.Name)
	assert.Equal(t, 20, response.Repetitions)
}

func TestUpdateWorkoutExercise_UnauthorizedUser(t *testing.T) {
	exerciseRepo := &MockWorkoutExerciseRepository{}
	workoutRepo := &MockWorkoutRepository{
		GetByIdFn: func(ctx context.Context, id int) (models.Workout, error) {
			return models.Workout{
				Id:     1,
				UserId: 2, // Different user ID
				Name:   "Test Workout",
			}, nil
		},
	}
	useCase := setupWorkoutExerciseUsecase(exerciseRepo, workoutRepo)

	ctx := createContextWithUserId(1) // User ID 1 trying to update exercise in workout owned by user ID 2
	req := dto.UpdateWorkoutExerciseRequest{
		Name:        "Updated Push Up",
		WorkoutId:   1,
		Description: "Updated description",
		Repetitions: 20,
		Sets:        4,
		Weight:      5.0,
	}

	_, err := useCase.Update(ctx, 1, req)

	assert.Error(t, err)
}

func TestDeleteWorkoutExercise_Success(t *testing.T) {
	exerciseRepo := &MockWorkoutExerciseRepository{}
	workoutRepo := &MockWorkoutRepository{}
	useCase := setupWorkoutExerciseUsecase(exerciseRepo, workoutRepo)

	ctx := createContextWithUserId(1)

	err := useCase.Delete(ctx, 1)

	assert.NoError(t, err)
}

func TestDeleteWorkoutExercise_UnauthorizedUser(t *testing.T) {
	exerciseRepo := &MockWorkoutExerciseRepository{}
	workoutRepo := &MockWorkoutRepository{
		GetByIdFn: func(ctx context.Context, id int) (models.Workout, error) {
			return models.Workout{
				Id:     1,
				UserId: 2, // Different user ID
				Name:   "Test Workout",
			}, nil
		},
	}
	useCase := setupWorkoutExerciseUsecase(exerciseRepo, workoutRepo)

	ctx := createContextWithUserId(1) // User ID 1 trying to delete exercise in workout owned by user ID 2

	err := useCase.Delete(ctx, 1)

	assert.Error(t, err)
}

func TestGetWorkoutExerciseById_Success(t *testing.T) {
	exerciseRepo := &MockWorkoutExerciseRepository{}
	workoutRepo := &MockWorkoutRepository{}
	useCase := setupWorkoutExerciseUsecase(exerciseRepo, workoutRepo)

	ctx := createContextWithUserId(1)

	response, err := useCase.GetById(ctx, 1)

	assert.NoError(t, err)
	assert.Equal(t, 1, response.Id)
	assert.Equal(t, "Test Exercise", response.Name)
	assert.Equal(t, 10, response.Repetitions)
}

func TestGetWorkoutExerciseById_UnauthorizedUser(t *testing.T) {
	exerciseRepo := &MockWorkoutExerciseRepository{}
	workoutRepo := &MockWorkoutRepository{
		GetByIdFn: func(ctx context.Context, id int) (models.Workout, error) {
			return models.Workout{
				Id:     1,
				UserId: 2, // Different user ID
				Name:   "Test Workout",
			}, nil
		},
	}
	useCase := setupWorkoutExerciseUsecase(exerciseRepo, workoutRepo)

	ctx := createContextWithUserId(1) // User ID 1 trying to access exercise in workout owned by user ID 2

	_, err := useCase.GetById(ctx, 1)

	assert.Error(t, err)
}

func TestGetWorkoutExerciseById_NotFound(t *testing.T) {
	exerciseRepo := &MockWorkoutExerciseRepository{
		GetByIdFn: func(ctx context.Context, id int) (models.WorkoutExercise, error) {
			return models.WorkoutExercise{}, errors.New("exercise not found")
		},
	}
	workoutRepo := &MockWorkoutRepository{}
	useCase := setupWorkoutExerciseUsecase(exerciseRepo, workoutRepo)

	ctx := createContextWithUserId(1)

	_, err := useCase.GetById(ctx, 999)

	assert.Error(t, err)
}
