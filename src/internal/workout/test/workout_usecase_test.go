package test

import (
	"context"
	"errors"
	"testing"

	"github.com/alecthomas/assert/v2"
	"github.com/alielmi98/go-hexa-workout/internal/workout/core/models"
	"github.com/alielmi98/go-hexa-workout/internal/workout/core/usecase/dto"
	"github.com/alielmi98/go-hexa-workout/internal/workout/port/filter"
	"github.com/alielmi98/go-hexa-workout/pkg/service_errors"
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

// ==================== GET BY FILTER USECASE TESTS ====================

func TestGetByFilterWorkout_Success(t *testing.T) {
	workoutRepo := &MockWorkoutRepository{
		GetByFilterFn: func(ctx context.Context, req filter.PaginationInputWithFilter) (int64, *[]models.Workout, error) {
			// Verify that user filter was automatically injected
			userFilter, exists := req.DynamicFilter.Filter["UserId"]
			assert.Equal(t, true, exists)
			assert.Equal(t, "equals", userFilter.Type)
			assert.Equal(t, "1", userFilter.From)
			assert.Equal(t, "number", userFilter.FilterType)

			workouts := []models.Workout{
				{
					Id:          1,
					UserId:      1,
					Name:        "Test Workout 1",
					Description: "Test Description 1",
					Comments:    "Test Comments 1",
				},
				{
					Id:          2,
					UserId:      1,
					Name:        "Test Workout 2",
					Description: "Test Description 2",
					Comments:    "Test Comments 2",
				},
			}
			return 2, &workouts, nil
		},
	}
	useCase := setupWorkoutUsecase(workoutRepo)

	ctx := createContextWithUserId(1)
	req := filter.PaginationInputWithFilter{
		PaginationInput: filter.PaginationInput{
			PageNumber: 1,
			PageSize:   10,
		},
		DynamicFilter: filter.DynamicFilter{
			Filter: map[string]filter.Filter{
				"name": {
					Type:       "contains",
					From:       "Test",
					FilterType: "string",
				},
			},
		},
	}

	response, err := useCase.GetByFilter(ctx, req)

	assert.NoError(t, err)
	if response == nil {
		t.Fatal("Expected response to not be nil")
	}
	assert.Equal(t, int64(2), response.TotalRows)
	assert.Equal(t, 2, len(*response.Items))
	assert.Equal(t, "Test Workout 1", (*response.Items)[0].Name)
	assert.Equal(t, "Test Workout 2", (*response.Items)[1].Name)
}

func TestGetByFilterWorkout_EmptyResult(t *testing.T) {
	workoutRepo := &MockWorkoutRepository{
		GetByFilterFn: func(ctx context.Context, req filter.PaginationInputWithFilter) (int64, *[]models.Workout, error) {
			// Verify user filter injection
			userFilter, exists := req.DynamicFilter.Filter["UserId"]
			assert.Equal(t, true, exists)
			assert.Equal(t, "1", userFilter.From)

			emptyWorkouts := []models.Workout{}
			return 0, &emptyWorkouts, nil
		},
	}
	useCase := setupWorkoutUsecase(workoutRepo)

	ctx := createContextWithUserId(1)
	req := filter.PaginationInputWithFilter{
		PaginationInput: filter.PaginationInput{
			PageNumber: 1,
			PageSize:   10,
		},
		DynamicFilter: filter.DynamicFilter{
			Filter: map[string]filter.Filter{
				"name": {
					Type:       "contains",
					From:       "NonExistent",
					FilterType: "string",
				},
			},
		},
	}

	response, err := useCase.GetByFilter(ctx, req)

	assert.NoError(t, err)
	if response == nil {
		t.Fatal("Expected response to not be nil")
	}
	assert.Equal(t, int64(0), response.TotalRows)
	assert.Equal(t, 0, len(*response.Items))
}

func TestGetByFilterWorkout_RepositoryError(t *testing.T) {
	workoutRepo := &MockWorkoutRepository{
		GetByFilterFn: func(ctx context.Context, req filter.PaginationInputWithFilter) (int64, *[]models.Workout, error) {
			return 0, nil, &service_errors.ServiceError{EndUserMessage: "database connection failed"}
		},
	}
	useCase := setupWorkoutUsecase(workoutRepo)

	ctx := createContextWithUserId(1)
	req := filter.PaginationInputWithFilter{
		PaginationInput: filter.PaginationInput{
			PageNumber: 1,
			PageSize:   10,
		},
	}

	_, err := useCase.GetByFilter(ctx, req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "database connection failed")
}

func TestGetByFilterWorkout_UserFilterInjection(t *testing.T) {
	workoutRepo := &MockWorkoutRepository{
		GetByFilterFn: func(ctx context.Context, req filter.PaginationInputWithFilter) (int64, *[]models.Workout, error) {
			// Verify that user filter was automatically injected even with empty filter map
			userFilter, exists := req.DynamicFilter.Filter["UserId"]
			assert.Equal(t, true, exists)
			assert.Equal(t, "equals", userFilter.Type)
			assert.Equal(t, "123", userFilter.From) // User ID 123
			assert.Equal(t, "number", userFilter.FilterType)

			workouts := []models.Workout{
				{
					Id:          1,
					UserId:      123,
					Name:        "User 123 Workout",
					Description: "Test Description",
				},
			}
			return 1, &workouts, nil
		},
	}
	useCase := setupWorkoutUsecase(workoutRepo)

	ctx := createContextWithUserId(123)
	req := filter.PaginationInputWithFilter{
		PaginationInput: filter.PaginationInput{
			PageNumber: 1,
			PageSize:   10,
		},
		// Empty DynamicFilter - user filter should still be injected
	}

	response, err := useCase.GetByFilter(ctx, req)

	assert.NoError(t, err)
	if response == nil {
		t.Fatal("Expected response to not be nil")
	}
	assert.Equal(t, int64(1), response.TotalRows)
	assert.Equal(t, 1, len(*response.Items))
	assert.Equal(t, "User 123 Workout", (*response.Items)[0].Name)
	assert.Equal(t, 123, (*response.Items)[0].UserId)
}

func TestGetByFilterWorkout_PaginationTest(t *testing.T) {
	workoutRepo := &MockWorkoutRepository{
		GetByFilterFn: func(ctx context.Context, req filter.PaginationInputWithFilter) (int64, *[]models.Workout, error) {
			// Verify pagination parameters
			assert.Equal(t, 2, req.GetPageNumber())
			assert.Equal(t, 5, req.GetPageSize())
			assert.Equal(t, 5, req.GetOffset()) // (2-1) * 5 = 5

			// Return page 2 data (items 6-10 out of 15 total)
			workouts := []models.Workout{
				{Id: 6, UserId: 1, Name: "Workout 6"},
				{Id: 7, UserId: 1, Name: "Workout 7"},
				{Id: 8, UserId: 1, Name: "Workout 8"},
				{Id: 9, UserId: 1, Name: "Workout 9"},
				{Id: 10, UserId: 1, Name: "Workout 10"},
			}
			return 15, &workouts, nil // Total 15 items
		},
	}
	useCase := setupWorkoutUsecase(workoutRepo)

	ctx := createContextWithUserId(1)
	req := filter.PaginationInputWithFilter{
		PaginationInput: filter.PaginationInput{
			PageNumber: 2,
			PageSize:   5,
		},
	}

	response, err := useCase.GetByFilter(ctx, req)

	assert.NoError(t, err)
	if response == nil {
		t.Fatal("Expected response to not be nil")
	}
	assert.Equal(t, int64(15), response.TotalRows)
	assert.Equal(t, 3, response.TotalPages) // 15/5 = 3 pages
	assert.Equal(t, 2, response.PageNumber)
	assert.Equal(t, int64(5), response.PageSize)
	assert.Equal(t, true, response.HasPreviousPage)
	assert.Equal(t, true, response.HasNextPage)
	assert.Equal(t, 5, len(*response.Items))
	assert.Equal(t, "Workout 6", (*response.Items)[0].Name)
	assert.Equal(t, "Workout 10", (*response.Items)[4].Name)
}
