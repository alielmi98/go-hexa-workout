package test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/alecthomas/assert/v2"
	"github.com/alielmi98/go-hexa-workout/internal/workout/adapter/http/dto"
	"github.com/alielmi98/go-hexa-workout/internal/workout/adapter/http/handler"
	"github.com/alielmi98/go-hexa-workout/internal/workout/core/models"
	"github.com/alielmi98/go-hexa-workout/internal/workout/core/usecase"
	"github.com/alielmi98/go-hexa-workout/pkg/config"
	"github.com/alielmi98/go-hexa-workout/pkg/helper"
	"github.com/alielmi98/go-hexa-workout/pkg/service_errors"
	"github.com/gin-gonic/gin"
)

func setupScheduledWorkoutHandler(scheduledRepo *MockScheduledWorkoutsRepository, workoutRepo *MockWorkoutRepository) (*handler.ScheduledWorkoutsHandler, *MockTokenProvider, *config.Config) {
	cfg := &config.Config{}
	tokenProvider := &MockTokenProvider{}
	useCase := usecase.NewScheduledWorkoutsUsecase(cfg, scheduledRepo, workoutRepo)
	return &handler.ScheduledWorkoutsHandler{
		Usecase: useCase,
	}, tokenProvider, cfg
}

func TestCreateScheduledWorkout_Handler_Success(t *testing.T) {
	scheduledRepo := &MockScheduledWorkoutsRepository{}
	workoutRepo := &MockWorkoutRepository{}
	handler, tokenProvider, cfg := setupScheduledWorkoutHandler(scheduledRepo, workoutRepo)

	requestBody := dto.CreateScheduledWorkoutsRequest{
		WorkoutId:     1,
		ScheduledTime: time.Now(),
		Status:        "active",
	}

	jsonBody, _ := json.Marshal(requestBody)
	c, w := createAuthenticatedGinContext("POST", "/v1/workouts/scheduled-workouts/", jsonBody, tokenProvider, cfg)

	// Set up the route and call the handler
	c.Request.URL.Path = "/v1/workouts/scheduled-workouts/"
	handler.Create(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	var response helper.BaseHttpResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, true, response.Success)
}

func TestCreateScheduledWorkout_Handler_ValidationError(t *testing.T) {
	scheduledRepo := &MockScheduledWorkoutsRepository{}
	workoutRepo := &MockWorkoutRepository{}
	handler, tokenProvider, cfg := setupScheduledWorkoutHandler(scheduledRepo, workoutRepo)

	// Invalid request body - missing required fields
	requestBody := dto.CreateScheduledWorkoutsRequest{
		Status: "active",
		// Missing WorkoutId and ScheduledTime
	}

	jsonBody, _ := json.Marshal(requestBody)
	c, w := createAuthenticatedGinContext("POST", "/v1/workouts/scheduled-workouts/", jsonBody, tokenProvider, cfg)

	// Set up the route and call the handler
	c.Request.URL.Path = "/v1/workouts/scheduled-workouts/"
	handler.Create(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var response helper.BaseHttpResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, false, response.Success)
	assert.Equal(t, helper.ValidationError, response.ResultCode)

}

func TestCreateScheduledWorkout_Handler_InvalidStatus(t *testing.T) {
	scheduledRepo := &MockScheduledWorkoutsRepository{}
	workoutRepo := &MockWorkoutRepository{}
	handler, tokenProvider, cfg := setupScheduledWorkoutHandler(scheduledRepo, workoutRepo)

	requestBody := dto.CreateScheduledWorkoutsRequest{
		WorkoutId:     1,
		ScheduledTime: time.Now(),
		Status:        "invalid_status",
	}

	jsonBody, _ := json.Marshal(requestBody)
	c, w := createAuthenticatedGinContext("POST", "/v1/workouts/scheduled-workouts/", jsonBody, tokenProvider, cfg)

	// Set up the route and call the handler
	c.Request.URL.Path = "/v1/workouts/scheduled-workouts/"
	handler.Create(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	var response helper.BaseHttpResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, false, response.Success)
	assert.Equal(t, helper.InternalError, response.ResultCode)
	assert.Equal(t, "invalid status. Status must be 'active' or 'completed' or 'canceled'", response.Error)

}

func TestCreateScheduledWorkout_Handler_UnauthorizedUser(t *testing.T) {
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
	handler, tokenProvider, cfg := setupScheduledWorkoutHandler(scheduledRepo, workoutRepo)

	requestBody := dto.CreateScheduledWorkoutsRequest{
		WorkoutId:     1,
		ScheduledTime: time.Now(),
		Status:        "active",
	}

	jsonBody, _ := json.Marshal(requestBody)
	c, w := createAuthenticatedGinContext("POST", "/v1/workouts/scheduled-workouts/", jsonBody, tokenProvider, cfg)

	// Set up the route and call the handler
	c.Request.URL.Path = "/v1/workouts/scheduled-workouts/"
	handler.Create(c)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	var response helper.BaseHttpResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, false, response.Success)
	assert.Equal(t, helper.InternalError, response.ResultCode)
	assert.Equal(t, "user is not the owner of this workout", response.Error)
}

func TestGetScheduledWorkoutById_Handler_Success(t *testing.T) {
	scheduledRepo := &MockScheduledWorkoutsRepository{}
	workoutRepo := &MockWorkoutRepository{}
	handler, tokenProvider, cfg := setupScheduledWorkoutHandler(scheduledRepo, workoutRepo)

	// Use the helper function with parameters
	params := gin.Params{{Key: "id", Value: "1"}}
	c, w := createAuthenticatedGinContextWithParams("GET", "/v1/workouts/scheduled-workouts/1", nil, params, tokenProvider, cfg)

	// Set up the route and call the handler
	c.Request.URL.Path = "/v1/workouts/scheduled-workouts/1"
	handler.GetById(c)
	assert.Equal(t, http.StatusOK, w.Code)
	var response helper.BaseHttpResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, true, response.Success)
}

func TestGetScheduledWorkoutById_Handler_InvalidId(t *testing.T) {
	scheduledRepo := &MockScheduledWorkoutsRepository{}
	workoutRepo := &MockWorkoutRepository{}
	handler, tokenProvider, cfg := setupScheduledWorkoutHandler(scheduledRepo, workoutRepo)

	// Use the helper function with parameters
	params := gin.Params{{Key: "id", Value: "invalid"}}
	c, w := createAuthenticatedGinContextWithParams("GET", "/v1/workouts/scheduled-workouts/invalid", nil, params, tokenProvider, cfg)

	// Set up the route and call the handler
	c.Request.URL.Path = "/v1/workouts/scheduled-workouts/invalid"
	handler.GetById(c)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	var response helper.BaseHttpResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, false, response.Success)
	assert.Equal(t, helper.ValidationError, response.ResultCode)
	assert.Equal(t, "strconv.Atoi: parsing \"invalid\": invalid syntax", response.Error)

}

func TestGetScheduledWorkoutById_Handler_NotFound(t *testing.T) {
	scheduledRepo := &MockScheduledWorkoutsRepository{
		GetByIdFn: func(ctx context.Context, id int) (models.ScheduledWorkouts, error) {
			return models.ScheduledWorkouts{}, &service_errors.ServiceError{EndUserMessage: service_errors.RecordNotFound}
		},
	}
	workoutRepo := &MockWorkoutRepository{}
	handler, tokenProvider, cfg := setupScheduledWorkoutHandler(scheduledRepo, workoutRepo)

	// Use the helper function with parameters
	params := gin.Params{{Key: "id", Value: "2"}}
	c, w := createAuthenticatedGinContextWithParams("GET", "/v1/workouts/scheduled-workouts/2", nil, params, tokenProvider, cfg)

	// Set up the route and call the handler
	c.Request.URL.Path = "/v1/workouts/scheduled-workouts/2"
	handler.GetById(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
	var response helper.BaseHttpResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, false, response.Success)
	assert.Equal(t, "record not found", response.Error)
}

func TestUpdateScheduledWorkout_Handler_Success(t *testing.T) {
	scheduledRepo := &MockScheduledWorkoutsRepository{}
	workoutRepo := &MockWorkoutRepository{}
	handler, tokenProvider, cfg := setupScheduledWorkoutHandler(scheduledRepo, workoutRepo)

	// Use the helper function with parameters
	params := gin.Params{{Key: "id", Value: "1"}}
	requestBody := dto.UpdateScheduledWorkoutsRequest{
		Status:        "active",
		ScheduledTime: time.Now(),
	}
	jsonBody, _ := json.Marshal(requestBody)
	c, w := createAuthenticatedGinContextWithParams("PUT", "/v1/workouts/scheduled-workouts/1", jsonBody, params, tokenProvider, cfg)

	// Set up the route and call the handler
	c.Request.URL.Path = "/v1/workouts/scheduled-workouts/1"
	handler.Update(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var response helper.BaseHttpResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, true, response.Success)
}

func TestUpdateScheduledWorkout_Handler_InvalidId(t *testing.T) {
	scheduledRepo := &MockScheduledWorkoutsRepository{
		UpdateFn: func(ctx context.Context, id int, scheduledWorkouts models.ScheduledWorkouts) (models.ScheduledWorkouts, error) {
			return models.ScheduledWorkouts{}, &service_errors.ServiceError{EndUserMessage: service_errors.RecordNotFound}
		},
	}
	workoutRepo := &MockWorkoutRepository{}
	handler, tokenProvider, cfg := setupScheduledWorkoutHandler(scheduledRepo, workoutRepo)
	// Use the helper function with parameters
	params := gin.Params{{Key: "id", Value: "invalid"}}
	requestBody := dto.UpdateScheduledWorkoutsRequest{
		Status:        "active",
		ScheduledTime: time.Now(),
	}
	jsonBody, _ := json.Marshal(requestBody)
	c, w := createAuthenticatedGinContextWithParams("PUT", "/v1/workouts/scheduled-workouts/invalid", jsonBody, params, tokenProvider, cfg)

	// Set up the route and call the handler
	c.Request.URL.Path = "/v1/workouts/scheduled-workouts/invalid"
	handler.Update(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response helper.BaseHttpResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, false, response.Success)
	assert.Equal(t, "strconv.Atoi: parsing \"invalid\": invalid syntax", response.Error)

}

func TestDeleteScheduledWorkout_Handler_InvalidId(t *testing.T) {
	scheduledRepo := &MockScheduledWorkoutsRepository{}
	workoutRepo := &MockWorkoutRepository{}
	handler, tokenProvider, cfg := setupScheduledWorkoutHandler(scheduledRepo, workoutRepo)

	// Use the helper function with parameters
	params := gin.Params{{Key: "id", Value: "invalid"}}
	c, w := createAuthenticatedGinContextWithParams("DELETE", "/v1/workouts/scheduled-workouts/invalid", nil, params, tokenProvider, cfg)

	// Set up the route and call the handler
	c.Request.URL.Path = "/v1/workouts/scheduled-workouts/invalid"
	handler.Delete(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response helper.BaseHttpResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, false, response.Success)
	assert.Equal(t, "strconv.Atoi: parsing \"invalid\": invalid syntax", response.Error)

}

func TestDeleteScheduledWorkout_Handler_Success(t *testing.T) {
	scheduledRepo := &MockScheduledWorkoutsRepository{}
	workoutRepo := &MockWorkoutRepository{}
	handler, tokenProvider, cfg := setupScheduledWorkoutHandler(scheduledRepo, workoutRepo)

	// Use the helper function with parameters
	params := gin.Params{{Key: "id", Value: "1"}}
	c, w := createAuthenticatedGinContextWithParams("DELETE", "/v1/workouts/scheduled-workouts/1", nil, params, tokenProvider, cfg)

	// Set up the route and call the handler
	c.Request.URL.Path = "/v1/workouts/scheduled-workouts/1"
	handler.Delete(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var response helper.BaseHttpResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, true, response.Success)
}

func TestDeleteScheduledWorkout_Handler_NotFound(t *testing.T) {
	scheduledRepo := &MockScheduledWorkoutsRepository{
		DeleteFn: func(ctx context.Context, id int) error {
			return &service_errors.ServiceError{EndUserMessage: service_errors.RecordNotFound}
		},
	}
	workoutRepo := &MockWorkoutRepository{}
	handler, tokenProvider, cfg := setupScheduledWorkoutHandler(scheduledRepo, workoutRepo)

	// Use the helper function with parameters
	params := gin.Params{{Key: "id", Value: "1"}}
	c, w := createAuthenticatedGinContextWithParams("DELETE", "/v1/workouts/scheduled-workouts/1", nil, params, tokenProvider, cfg)

	// Set up the route and call the handler
	c.Request.URL.Path = "/v1/workouts/scheduled-workouts/1"
	handler.Delete(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
	var response helper.BaseHttpResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, false, response.Success)
	assert.Equal(t, service_errors.RecordNotFound, response.Error)
}
