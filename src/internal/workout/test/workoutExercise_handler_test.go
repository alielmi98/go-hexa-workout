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
	"github.com/alielmi98/go-hexa-workout/pkg/config"
	"github.com/alielmi98/go-hexa-workout/pkg/helper"
	"github.com/alielmi98/go-hexa-workout/pkg/service_errors"
	"github.com/gin-gonic/gin"
)

func setupWorkoutExerciseHandler(exerciseRepo *MockWorkoutExerciseRepository, workoutRepo *MockWorkoutRepository) (*handler.WorkoutExerciseHandler, *MockTokenProvider, *config.Config) {
	cfg := &config.Config{}
	tokenProvider := &MockTokenProvider{}
	useCase := usecase.NewWorkoutExerciseUsecase(cfg, exerciseRepo, workoutRepo)
	return &handler.WorkoutExerciseHandler{
		Usecase: useCase,
	}, tokenProvider, cfg
}

func TestCreateWorkoutExercise_Handler_Success(t *testing.T) {
	exerciseRepo := &MockWorkoutExerciseRepository{}
	workoutRepo := &MockWorkoutRepository{}
	handler, tokenProvider, cfg := setupWorkoutExerciseHandler(exerciseRepo, workoutRepo)

	requestBody := dto.CreateWorkoutExerciseRequest{
		WorkoutId:   1,
		Name:        "Push Up",
		Description: "Basic push up exercise",
		Reps:        10,
		Sets:        3,
		Weight:      0.0,
	}

	jsonBody, _ := json.Marshal(requestBody)
	c, w := createAuthenticatedGinContext("POST", "/v1/workouts/workout-exercise/", jsonBody, tokenProvider, cfg)

	// Set up the route and call the handler
	c.Request.URL.Path = "/v1/workouts/workout-exercise/"
	handler.Create(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var response helper.BaseHttpResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, false, response.Success)
}

func TestCreateWorkoutExercise_Handler_ValidationError(t *testing.T) {
	exerciseRepo := &MockWorkoutExerciseRepository{}
	workoutRepo := &MockWorkoutRepository{}
	handler, tokenProvider, cfg := setupWorkoutExerciseHandler(exerciseRepo, workoutRepo)

	// Invalid request body - missing required fields
	requestBody := dto.CreateWorkoutExerciseRequest{
		Name:        "Pu", // Too short (min=3)
		Description: "Basic push up exercise",
		// Missing WorkoutId, Reps, Sets, Weight
	}

	jsonBody, _ := json.Marshal(requestBody)
	c, w := createAuthenticatedGinContext("POST", "/v1/workouts/workout-exercise/", jsonBody, tokenProvider, cfg)

	// Set up the route and call the handler
	c.Request.URL.Path = "/v1/workouts/workout-exercise/"
	handler.Create(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var response helper.BaseHttpResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, false, response.Success)
	assert.Equal(t, helper.ValidationError, response.ResultCode)
}

func TestCreateWorkoutExercise_Handler_UnauthorizedWorkout(t *testing.T) {
	exerciseRepo := &MockWorkoutExerciseRepository{}
	workoutRepo := &MockWorkoutRepository{
		GetByIdFn: func(ctx context.Context, id int) (models.Workout, error) {
			return models.Workout{
				Id:     1,
				UserId: 2, // Different user
				Name:   "Test Workout",
			}, nil
		},
	}
	handler, tokenProvider, cfg := setupWorkoutExerciseHandler(exerciseRepo, workoutRepo)

	requestBody := dto.CreateWorkoutExerciseRequest{
		WorkoutId:   1,
		Name:        "Push Up",
		Description: "Basic push up exercise",
		Reps:        10,
		Sets:        3,
		Weight:      0.0,
	}

	jsonBody, _ := json.Marshal(requestBody)
	c, w := createAuthenticatedGinContext("POST", "/v1/workouts/workout-exercise/", jsonBody, tokenProvider, cfg)

	// Set up the route and call the handler
	c.Request.URL.Path = "/v1/workouts/workout-exercise/"
	handler.Create(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var response helper.BaseHttpResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, false, response.Success)
}

func TestGetWorkoutExerciseById_Handler_Success(t *testing.T) {
	exerciseRepo := &MockWorkoutExerciseRepository{}
	workoutRepo := &MockWorkoutRepository{}
	handler, tokenProvider, cfg := setupWorkoutExerciseHandler(exerciseRepo, workoutRepo)

	// Use the helper function with parameters
	params := gin.Params{{Key: "id", Value: "1"}}
	c, w := createAuthenticatedGinContextWithParams("GET", "/v1/workouts/workout-exercise/1", nil, params, tokenProvider, cfg)

	// Set up the route and call the handler
	c.Request.URL.Path = "/v1/workouts/workout-exercise/1"
	handler.GetById(c)
	assert.Equal(t, http.StatusOK, w.Code)
	var response helper.BaseHttpResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, true, response.Success)
}

func TestGetWorkoutExerciseById_Handler_InvalidId(t *testing.T) {
	exerciseRepo := &MockWorkoutExerciseRepository{}
	workoutRepo := &MockWorkoutRepository{}
	handler, tokenProvider, cfg := setupWorkoutExerciseHandler(exerciseRepo, workoutRepo)

	// Use the helper function with parameters
	params := gin.Params{{Key: "id", Value: "invalid"}}
	c, w := createAuthenticatedGinContextWithParams("GET", "/v1/workouts/workout-exercise/invalid", nil, params, tokenProvider, cfg)

	// Set up the route and call the handler
	c.Request.URL.Path = "/v1/workouts/workout-exercise/invalid"
	handler.GetById(c)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	var response helper.BaseHttpResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, false, response.Success)
	assert.Equal(t, helper.ValidationError, response.ResultCode)
}

func TestGetWorkoutExerciseById_Handler_NotFound(t *testing.T) {
	exerciseRepo := &MockWorkoutExerciseRepository{
		GetByIdFn: func(ctx context.Context, id int) (models.WorkoutExercise, error) {
			return models.WorkoutExercise{}, &service_errors.ServiceError{EndUserMessage: service_errors.RecordNotFound}
		},
	}
	workoutRepo := &MockWorkoutRepository{}
	handler, tokenProvider, cfg := setupWorkoutExerciseHandler(exerciseRepo, workoutRepo)

	// Use the helper function with parameters
	params := gin.Params{{Key: "id", Value: "1"}}
	c, w := createAuthenticatedGinContextWithParams("GET", "/v1/workouts/workout-exercise/1", nil, params, tokenProvider, cfg)

	// Set up the route and call the handler
	c.Request.URL.Path = "/v1/workouts/workout-exercise/1"
	handler.GetById(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
	var response helper.BaseHttpResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, false, response.Success)
	assert.Equal(t, "record not found", response.Error)
}

func TestUpdateWorkoutExercise_Handler_Success(t *testing.T) {
	exerciseRepo := &MockWorkoutExerciseRepository{}
	workoutRepo := &MockWorkoutRepository{}
	handler, tokenProvider, cfg := setupWorkoutExerciseHandler(exerciseRepo, workoutRepo)

	// Use the helper function with parameters
	params := gin.Params{{Key: "id", Value: "1"}}
	requestBody := dto.UpdateWorkoutExerciseRequest{
		WorkoutId:   1,
		Name:        "Updated Push Up",
		Description: "Updated push up exercise",
		Reps:        15,
		Sets:        4,
		Weight:      5.0,
	}
	jsonBody, _ := json.Marshal(requestBody)
	c, w := createAuthenticatedGinContextWithParams("PUT", "/v1/workouts/workout-exercise/1", jsonBody, params, tokenProvider, cfg)

	// Set up the route and call the handler
	c.Request.URL.Path = "/v1/workouts/workout-exercise/1"
	handler.Update(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var response helper.BaseHttpResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, true, response.Success)
}

func TestUpdateWorkoutExercise_Handler_InvalidId(t *testing.T) {
	exerciseRepo := &MockWorkoutExerciseRepository{}
	workoutRepo := &MockWorkoutRepository{}
	handler, tokenProvider, cfg := setupWorkoutExerciseHandler(exerciseRepo, workoutRepo)

	// Use the helper function with parameters
	params := gin.Params{{Key: "id", Value: "invalid"}}
	requestBody := dto.UpdateWorkoutExerciseRequest{
		WorkoutId:   1,
		Name:        "Updated Push Up",
		Description: "Updated push up exercise",
		Reps:        15,
		Sets:        4,
		Weight:      5.0,
	}
	jsonBody, _ := json.Marshal(requestBody)
	c, w := createAuthenticatedGinContextWithParams("PUT", "/v1/workouts/workout-exercise/invalid", jsonBody, params, tokenProvider, cfg)

	// Set up the route and call the handler
	c.Request.URL.Path = "/v1/workouts/workout-exercise/invalid"
	handler.Update(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response helper.BaseHttpResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, false, response.Success)
}

func TestDeleteWorkoutExercise_Handler_InvalidId(t *testing.T) {
	exerciseRepo := &MockWorkoutExerciseRepository{}
	workoutRepo := &MockWorkoutRepository{}
	handler, tokenProvider, cfg := setupWorkoutExerciseHandler(exerciseRepo, workoutRepo)

	// Use the helper function with parameters
	params := gin.Params{{Key: "id", Value: "invalid"}}
	c, w := createAuthenticatedGinContextWithParams("DELETE", "/v1/workouts/workout-exercise/invalid", nil, params, tokenProvider, cfg)

	// Set up the route and call the handler
	c.Request.URL.Path = "/v1/workouts/workout-exercise/invalid"
	handler.Delete(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response helper.BaseHttpResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, false, response.Success)
}

func TestDeleteWorkoutExercise_Handler_Success(t *testing.T) {
	exerciseRepo := &MockWorkoutExerciseRepository{}
	workoutRepo := &MockWorkoutRepository{}
	handler, tokenProvider, cfg := setupWorkoutExerciseHandler(exerciseRepo, workoutRepo)

	// Use the helper function with parameters
	params := gin.Params{{Key: "id", Value: "1"}}
	c, w := createAuthenticatedGinContextWithParams("DELETE", "/v1/workouts/workout-exercise/1", nil, params, tokenProvider, cfg)

	// Set up the route and call the handler
	c.Request.URL.Path = "/v1/workouts/workout-exercise/1"
	handler.Delete(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var response helper.BaseHttpResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, true, response.Success)
}

func TestDeleteWorkoutExercise_Handler_NotFound(t *testing.T) {
	exerciseRepo := &MockWorkoutExerciseRepository{
		DeleteFn: func(ctx context.Context, id int) error {
			return &service_errors.ServiceError{EndUserMessage: service_errors.RecordNotFound}
		},
	}
	workoutRepo := &MockWorkoutRepository{}
	handler, tokenProvider, cfg := setupWorkoutExerciseHandler(exerciseRepo, workoutRepo)

	// Use the helper function with parameters
	params := gin.Params{{Key: "id", Value: "7"}}
	c, w := createAuthenticatedGinContextWithParams("DELETE", "/v1/workouts/workout-exercise/7", nil, params, tokenProvider, cfg)

	// Set up the route and call the handler
	c.Request.URL.Path = "/v1/workouts/workout-exercise/7"
	handler.Delete(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
	var response helper.BaseHttpResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, false, response.Success)
	assert.Equal(t, service_errors.RecordNotFound, response.Error)
}
