package test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/alecthomas/assert/v2"
	"github.com/alielmi98/go-hexa-workout/internal/workout/adapter/http/dto"
	"github.com/alielmi98/go-hexa-workout/internal/workout/adapter/http/handler"
	"github.com/alielmi98/go-hexa-workout/internal/workout/core/models"
	"github.com/alielmi98/go-hexa-workout/internal/workout/core/usecase"
	"github.com/alielmi98/go-hexa-workout/internal/workout/port/filter"
	"github.com/alielmi98/go-hexa-workout/pkg/config"
	"github.com/alielmi98/go-hexa-workout/pkg/helper"
	"github.com/alielmi98/go-hexa-workout/pkg/service_errors"
	"github.com/gin-gonic/gin"
)

func setupWorkoutHandler(workoutRepo *MockWorkoutRepository) (*handler.WorkoutHandler, *MockTokenProvider, *config.Config) {
	cfg := &config.Config{}
	tokenProvider := &MockTokenProvider{}
	useCase := usecase.NewWorkoutUsecase(cfg, workoutRepo)
	return &handler.WorkoutHandler{
		Usecase: useCase,
	}, tokenProvider, cfg
}

func TestCreateWorkout_Handler_Success(t *testing.T) {
	workoutRepo := &MockWorkoutRepository{}
	handler, tokenProvider, cfg := setupWorkoutHandler(workoutRepo)

	requestBody := dto.CreateWorkoutRequest{
		Name:        "Test Workout",
		Description: "Test Description",
		Comments:    "Test Comments",
	}

	jsonBody, _ := json.Marshal(requestBody)
	c, w := createAuthenticatedGinContext("POST", "/v1/workouts/workout/", jsonBody, tokenProvider, cfg)

	// Set up the route and call the handler
	c.Request.URL.Path = "/v1/workouts/workout/"
	handler.Create(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	var response helper.BaseHttpResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, true, response.Success)
}
func TestCreateWorkout_Handler_ValidationError(t *testing.T) {
	workoutRepo := &MockWorkoutRepository{}
	handler, tokenProvider, cfg := setupWorkoutHandler(workoutRepo)

	// Invalid request body - missing required fields
	requestBody := dto.CreateWorkoutRequest{
		Description: "Test Description",
		Comments:    "Test Comments",
		// Missing Name
	}

	jsonBody, _ := json.Marshal(requestBody)
	c, w := createAuthenticatedGinContext("POST", "/v1/workouts/workout/", jsonBody, tokenProvider, cfg)

	// Set up the route and call the handler
	c.Request.URL.Path = "/v1/workouts/workout/"
	handler.Create(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var response helper.BaseHttpResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, false, response.Success)
	assert.Equal(t, helper.ValidationError, response.ResultCode)

}

// GetbyId
func TestGetWorkoutById_Handler_Success(t *testing.T) {
	workoutRepo := &MockWorkoutRepository{}
	handler, tokenProvider, cfg := setupWorkoutHandler(workoutRepo)

	// Use the helper function with parameters
	params := gin.Params{{Key: "id", Value: "1"}}
	c, w := createAuthenticatedGinContextWithParams("GET", "/v1/workouts/workout/1", nil, params, tokenProvider, cfg)

	// Set up the route and call the handler
	c.Request.URL.Path = "/v1/workouts/workout/1"
	handler.GetById(c)
	assert.Equal(t, http.StatusOK, w.Code)
	var response helper.BaseHttpResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, true, response.Success)
}

func TestGetWorkoutById_Handler_InvalidId(t *testing.T) {
	workoutRepo := &MockWorkoutRepository{}
	handler, tokenProvider, cfg := setupWorkoutHandler(workoutRepo)

	// Use the helper function with parameters
	params := gin.Params{{Key: "id", Value: "invalid"}}
	c, w := createAuthenticatedGinContextWithParams("GET", "/v1/workouts/workout/invalid", nil, params, tokenProvider, cfg)

	// Set up the route and call the handler
	c.Request.URL.Path = "/v1/workouts/workout/invalid"
	handler.GetById(c)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	var response helper.BaseHttpResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, false, response.Success)
	assert.Equal(t, helper.ValidationError, response.ResultCode)
	assert.Equal(t, "strconv.Atoi: parsing \"invalid\": invalid syntax", response.Error)

}

func TestGetWorkoutById_Handler_NotFound(t *testing.T) {
	workoutRepo := &MockWorkoutRepository{
		GetByIdFn: func(ctx context.Context, id int) (models.Workout, error) {
			return models.Workout{}, &service_errors.ServiceError{EndUserMessage: service_errors.RecordNotFound}
		},
	}
	handler, tokenProvider, cfg := setupWorkoutHandler(workoutRepo)

	// Use the helper function with parameters
	params := gin.Params{{Key: "id", Value: "1"}}
	c, w := createAuthenticatedGinContextWithParams("GET", "/v1/workouts/workout/1", nil, params, tokenProvider, cfg)

	// Set up the route and call the handler
	c.Request.URL.Path = "/v1/workouts/workout/1"
	handler.GetById(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
	var response helper.BaseHttpResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, false, response.Success)
	assert.Equal(t, "record not found", response.Error)
}

func TestUpdateWorkout_Handler_Success(t *testing.T) {
	workoutRepo := &MockWorkoutRepository{}
	handler, tokenProvider, cfg := setupWorkoutHandler(workoutRepo)
	// Use the helper function with parameters
	params := gin.Params{{Key: "id", Value: "1"}}
	requestBody := dto.UpdateWorkoutRequest{
		Name:        "Test Workout",
		Description: "Test Description",
		Comments:    "Test Comments",
	}
	jsonBody, _ := json.Marshal(requestBody)
	c, w := createAuthenticatedGinContextWithParams("PUT", "/v1/workouts/workout/1", jsonBody, params, tokenProvider, cfg)

	// Set up the route and call the handler
	c.Request.URL.Path = "/v1/workouts/workout/1"
	handler.Update(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var response helper.BaseHttpResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, true, response.Success)
}

func TestUpdateWorkout_Handler_InvalidId(t *testing.T) {
	workoutRepo := &MockWorkoutRepository{}
	handler, tokenProvider, cfg := setupWorkoutHandler(workoutRepo)
	// Use the helper function with parameters
	params := gin.Params{{Key: "id", Value: "invalid"}}
	requestBody := dto.UpdateWorkoutRequest{
		Name:        "Test Workout",
		Description: "Test Description",
		Comments:    "Test Comments",
	}
	jsonBody, _ := json.Marshal(requestBody)
	c, w := createAuthenticatedGinContextWithParams("PUT", "/v1/workouts/workout/invalid", jsonBody, params, tokenProvider, cfg)

	// Set up the route and call the handler
	c.Request.URL.Path = "/v1/workouts/workout/invalid"
	handler.Update(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response helper.BaseHttpResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, false, response.Success)
	assert.Equal(t, "strconv.Atoi: parsing \"invalid\": invalid syntax", response.Error)

}

func TestDeleteWorkout_Handler_InvalidId(t *testing.T) {
	workoutRepo := &MockWorkoutRepository{}
	handler, tokenProvider, cfg := setupWorkoutHandler(workoutRepo)

	// Use the helper function with parameters
	params := gin.Params{{Key: "id", Value: "invalid"}}
	c, w := createAuthenticatedGinContextWithParams("DELETE", "/v1/workouts/workout/invalid", nil, params, tokenProvider, cfg)

	// Set up the route and call the handler
	c.Request.URL.Path = "/v1/workouts/workout/invalid"
	handler.Delete(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response helper.BaseHttpResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, false, response.Success)
	assert.Equal(t, "strconv.Atoi: parsing \"invalid\": invalid syntax", response.Error)

}

func TestDeleteWorkout_Handler_Success(t *testing.T) {
	workoutRepo := &MockWorkoutRepository{}
	handler, tokenProvider, cfg := setupWorkoutHandler(workoutRepo)

	// Use the helper function with parameters
	params := gin.Params{{Key: "id", Value: "1"}}
	c, w := createAuthenticatedGinContextWithParams("DELETE", "/v1/workouts/workout/1", nil, params, tokenProvider, cfg)

	// Set up the route and call the handler
	c.Request.URL.Path = "/v1/workouts/workout/1"
	handler.Delete(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var response helper.BaseHttpResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, true, response.Success)
}

func TestDeleteWorkout_Handler_NotFound(t *testing.T) {
	workoutRepo := &MockWorkoutRepository{
		DeleteFn: func(ctx context.Context, id int) error {
			return &service_errors.ServiceError{EndUserMessage: service_errors.RecordNotFound}
		},
	}
	handler, tokenProvider, cfg := setupWorkoutHandler(workoutRepo)

	// Use the helper function with parameters
	params := gin.Params{{Key: "id", Value: "7"}}
	c, w := createAuthenticatedGinContextWithParams("DELETE", "/v1/workouts/workout/7", nil, params, tokenProvider, cfg)

	// Set up the route and call the handler
	c.Request.URL.Path = "/v1/workouts/workout/7"
	handler.Delete(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
	var response helper.BaseHttpResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, false, response.Success)
	assert.Equal(t, service_errors.RecordNotFound, response.Error)
}

func TestGetByFilterWorkout_Handler_Success(t *testing.T) {
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
	handler, tokenProvider, cfg := setupWorkoutHandler(workoutRepo)

	requestBody := filter.PaginationInputWithFilter{
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
	jsonBody, _ := json.Marshal(requestBody)

	c, w := createAuthenticatedGinContextWithParams("POST", "/v1/workouts/workout/get-by-filter", jsonBody, nil, tokenProvider, cfg)

	// Set up the route and call the handler
	c.Request.URL.Path = "/v1/workouts/workout/get-by-filter"
	handler.GetByFilter(c)
	assert.Equal(t, http.StatusOK, w.Code)

	var response helper.BaseHttpResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, true, response.Success)
	if response.Result != nil {
		t.Logf("Response result is not nil")
	}
}

func TestGetByFilterWorkout_Handler_InvalidJSON(t *testing.T) {
	workoutRepo := &MockWorkoutRepository{}
	handler, tokenProvider, cfg := setupWorkoutHandler(workoutRepo)

	// Invalid JSON body
	invalidJSON := []byte(`{"page": "invalid", "pageSize": "invalid"}`)

	c, w := createAuthenticatedGinContextWithParams("POST", "/v1/workouts/workout/get-by-filter", invalidJSON, nil, tokenProvider, cfg)

	c.Request.URL.Path = "/v1/workouts/workout/get-by-filter"
	handler.GetByFilter(c)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response helper.BaseHttpResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, false, response.Success)
}

func TestGetByFilterWorkout_Handler_RepositoryError(t *testing.T) {
	workoutRepo := &MockWorkoutRepository{
		GetByFilterFn: func(ctx context.Context, req filter.PaginationInputWithFilter) (int64, *[]models.Workout, error) {
			return 0, nil, &service_errors.ServiceError{EndUserMessage: "database connection failed"}
		},
	}
	handler, tokenProvider, cfg := setupWorkoutHandler(workoutRepo)

	requestBody := filter.PaginationInputWithFilter{
		PaginationInput: filter.PaginationInput{
			PageNumber: 1,
			PageSize:   10,
		},
	}
	jsonBody, _ := json.Marshal(requestBody)

	c, w := createAuthenticatedGinContextWithParams("POST", "/v1/workouts/workout/get-by-filter", jsonBody, nil, tokenProvider, cfg)

	c.Request.URL.Path = "/v1/workouts/workout/get-by-filter"
	handler.GetByFilter(c)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response helper.BaseHttpResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, false, response.Success)
}

func TestGetByFilterWorkout_Handler_EmptyResult(t *testing.T) {
	workoutRepo := &MockWorkoutRepository{
		GetByFilterFn: func(ctx context.Context, req filter.PaginationInputWithFilter) (int64, *[]models.Workout, error) {
			emptyWorkouts := []models.Workout{}
			return 0, &emptyWorkouts, nil
		},
	}
	handler, tokenProvider, cfg := setupWorkoutHandler(workoutRepo)

	requestBody := filter.PaginationInputWithFilter{
		PaginationInput: filter.PaginationInput{
			PageNumber: 1,
			PageSize:   10,
		},
		DynamicFilter: filter.DynamicFilter{
			Filter: map[string]filter.Filter{
				"name": {
					Type:       "contains",
					From:       "nonexistent",
					FilterType: "string",
				},
			},
		},
	}
	jsonBody, _ := json.Marshal(requestBody)

	c, w := createAuthenticatedGinContextWithParams("POST", "/v1/workouts/workout/get-by-filter", jsonBody, nil, tokenProvider, cfg)

	c.Request.URL.Path = "/v1/workouts/workout/get-by-filter"
	handler.GetByFilter(c)
	assert.Equal(t, http.StatusOK, w.Code)

	var response helper.BaseHttpResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, true, response.Success)
}
