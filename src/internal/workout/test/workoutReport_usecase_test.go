package test

import (
	"context"
	"errors"
	"testing"

	"github.com/alecthomas/assert/v2"
	"github.com/alielmi98/go-hexa-workout/internal/workout/core/models"
	"github.com/alielmi98/go-hexa-workout/internal/workout/core/usecase/dto"
)

// ==================== WORKOUT REPORT USECASE TESTS ====================

func TestCreateWorkoutReport_Success(t *testing.T) {
	reportRepo := &MockWorkoutReportRepository{}
	workoutRepo := &MockWorkoutRepository{}
	useCase := setupWorkoutReportUsecase(reportRepo, workoutRepo)

	ctx := createContextWithUserId(1)
	req := dto.CreateWorkoutReportRequest{
		WorkoutId: 1,
		Details:   "Great workout session today!",
		UserID:    1,
	}

	response, err := useCase.Create(ctx, req)

	assert.NoError(t, err)
	assert.Equal(t, 1, response.Id)
	assert.Equal(t, 1, response.WorkoutId)
	assert.Equal(t, "Great workout session today!", response.Details)
}

func TestCreateWorkoutReport_UnauthorizedUser(t *testing.T) {
	reportRepo := &MockWorkoutReportRepository{}
	workoutRepo := &MockWorkoutRepository{
		GetByIdFn: func(ctx context.Context, id int) (models.Workout, error) {
			return models.Workout{
				Id:     1,
				UserId: 2, // Different user ID
				Name:   "Test Workout",
			}, nil
		},
	}
	useCase := setupWorkoutReportUsecase(reportRepo, workoutRepo)

	ctx := createContextWithUserId(1) // User ID 1 trying to create report for workout owned by user ID 2
	req := dto.CreateWorkoutReportRequest{
		WorkoutId: 1,
		Details:   "Great workout session today!",
		UserID:    1,
	}

	_, err := useCase.Create(ctx, req)

	assert.Error(t, err)
}

func TestCreateWorkoutReport_RepositoryError(t *testing.T) {
	reportRepo := &MockWorkoutReportRepository{
		CreateFn: func(ctx context.Context, entity models.WorkoutReport) (models.WorkoutReport, error) {
			return models.WorkoutReport{}, errors.New("database error")
		},
	}
	workoutRepo := &MockWorkoutRepository{}
	useCase := setupWorkoutReportUsecase(reportRepo, workoutRepo)

	ctx := createContextWithUserId(1)
	req := dto.CreateWorkoutReportRequest{
		WorkoutId: 1,
		Details:   "Great workout session today!",
		UserID:    1,
	}

	_, err := useCase.Create(ctx, req)

	assert.Error(t, err)
	assert.Equal(t, "database error", err.Error())
}

func TestUpdateWorkoutReport_Success(t *testing.T) {
	reportRepo := &MockWorkoutReportRepository{}
	workoutRepo := &MockWorkoutRepository{}
	useCase := setupWorkoutReportUsecase(reportRepo, workoutRepo)

	ctx := createContextWithUserId(1)
	req := dto.UpdateWorkoutReportRequest{
		WorkoutId: 1,
		Details:   "Updated workout report details",
	}

	response, err := useCase.Update(ctx, 1, req)

	assert.NoError(t, err)
	assert.Equal(t, 1, response.Id)
	assert.Equal(t, "Updated workout report details", response.Details)
}

func TestUpdateWorkoutReport_UnauthorizedUser(t *testing.T) {
	reportRepo := &MockWorkoutReportRepository{}
	workoutRepo := &MockWorkoutRepository{
		GetByIdFn: func(ctx context.Context, id int) (models.Workout, error) {
			return models.Workout{
				Id:     1,
				UserId: 2, // Different user ID
				Name:   "Test Workout",
			}, nil
		},
	}
	useCase := setupWorkoutReportUsecase(reportRepo, workoutRepo)

	ctx := createContextWithUserId(1) // User ID 1 trying to update report for workout owned by user ID 2
	req := dto.UpdateWorkoutReportRequest{
		WorkoutId: 1,
		Details:   "Updated workout report details",
	}

	_, err := useCase.Update(ctx, 1, req)

	assert.Error(t, err)
}

func TestDeleteWorkoutReport_Success(t *testing.T) {
	reportRepo := &MockWorkoutReportRepository{}
	workoutRepo := &MockWorkoutRepository{}
	useCase := setupWorkoutReportUsecase(reportRepo, workoutRepo)

	ctx := createContextWithUserId(1)

	err := useCase.Delete(ctx, 1)

	assert.NoError(t, err)
}

func TestDeleteWorkoutReport_UnauthorizedUser(t *testing.T) {
	reportRepo := &MockWorkoutReportRepository{}
	workoutRepo := &MockWorkoutRepository{
		GetByIdFn: func(ctx context.Context, id int) (models.Workout, error) {
			return models.Workout{
				Id:     1,
				UserId: 2, // Different user ID
				Name:   "Test Workout",
			}, nil
		},
	}
	useCase := setupWorkoutReportUsecase(reportRepo, workoutRepo)

	ctx := createContextWithUserId(1) // User ID 1 trying to delete report for workout owned by user ID 2

	err := useCase.Delete(ctx, 1)

	assert.Error(t, err)
}

func TestGetWorkoutReportById_Success(t *testing.T) {
	reportRepo := &MockWorkoutReportRepository{}
	workoutRepo := &MockWorkoutRepository{}
	useCase := setupWorkoutReportUsecase(reportRepo, workoutRepo)

	ctx := createContextWithUserId(1)

	response, err := useCase.GetById(ctx, 1)

	assert.NoError(t, err)
	assert.Equal(t, 1, response.Id)
	assert.Equal(t, 1, response.WorkoutId)
	assert.Equal(t, "Test Report Details", response.Details)
}

func TestGetWorkoutReportById_UnauthorizedUser(t *testing.T) {
	reportRepo := &MockWorkoutReportRepository{}
	workoutRepo := &MockWorkoutRepository{
		GetByIdFn: func(ctx context.Context, id int) (models.Workout, error) {
			return models.Workout{
				Id:     1,
				UserId: 2, // Different user ID
				Name:   "Test Workout",
			}, nil
		},
	}
	useCase := setupWorkoutReportUsecase(reportRepo, workoutRepo)

	ctx := createContextWithUserId(1) // User ID 1 trying to access report for workout owned by user ID 2

	_, err := useCase.GetById(ctx, 1)

	assert.Error(t, err)
}

func TestGetWorkoutReportById_NotFound(t *testing.T) {
	reportRepo := &MockWorkoutReportRepository{
		GetByIdFn: func(ctx context.Context, id int) (models.WorkoutReport, error) {
			return models.WorkoutReport{}, errors.New("report not found")
		},
	}
	workoutRepo := &MockWorkoutRepository{}
	useCase := setupWorkoutReportUsecase(reportRepo, workoutRepo)

	ctx := createContextWithUserId(1)

	_, err := useCase.GetById(ctx, 999)

	assert.Error(t, err)
}
