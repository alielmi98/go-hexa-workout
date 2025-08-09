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

func setupWorkoutReportHandler(reportRepo *MockWorkoutReportRepository, workoutRepo *MockWorkoutRepository) (*handler.WorkoutReportHandler, *MockTokenProvider, *config.Config) {
	cfg := &config.Config{}
	tokenProvider := &MockTokenProvider{}
	useCase := usecase.NewWorkoutReportUsecase(cfg, reportRepo, workoutRepo)
	return &handler.WorkoutReportHandler{
		Usecase: useCase,
	}, tokenProvider, cfg
}

func TestCreateWorkoutReport_Handler_Success(t *testing.T) {
	reportRepo := &MockWorkoutReportRepository{}
	workoutRepo := &MockWorkoutRepository{}
	handler, tokenProvider, cfg := setupWorkoutReportHandler(reportRepo, workoutRepo)

	requestBody := dto.CreateWorkoutReportRequest{
		WorkoutId: 1,
		Details:   "Great workout session! Completed all sets with good form.",
	}

	jsonBody, _ := json.Marshal(requestBody)
	c, w := createAuthenticatedGinContext("POST", "/v1/workouts/workout-report/", jsonBody, tokenProvider, cfg)

	// Set up the route and call the handler
	c.Request.URL.Path = "/v1/workouts/workout-report/"
	handler.Create(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	var response helper.BaseHttpResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, true, response.Success)
}

func TestCreateWorkoutReport_Handler_ValidationError(t *testing.T) {
	reportRepo := &MockWorkoutReportRepository{}
	workoutRepo := &MockWorkoutRepository{}
	handler, tokenProvider, cfg := setupWorkoutReportHandler(reportRepo, workoutRepo)

	// Invalid request body - missing required fields
	requestBody := dto.CreateWorkoutReportRequest{
		// Missing WorkoutId and Details
	}

	jsonBody, _ := json.Marshal(requestBody)
	c, w := createAuthenticatedGinContext("POST", "/v1/workouts/workout-report/", jsonBody, tokenProvider, cfg)

	// Set up the route and call the handler
	c.Request.URL.Path = "/v1/workouts/workout-report/"
	handler.Create(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var response helper.BaseHttpResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, false, response.Success)
	assert.Equal(t, helper.ValidationError, response.ResultCode)
}

func TestCreateWorkoutReport_Handler_UnauthorizedWorkout(t *testing.T) {
	reportRepo := &MockWorkoutReportRepository{}
	workoutRepo := &MockWorkoutRepository{
		GetByIdFn: func(ctx context.Context, id int) (models.Workout, error) {
			return models.Workout{
				Id:     1,
				UserId: 2, // Different user
				Name:   "Test Workout",
			}, nil
		},
	}
	handler, tokenProvider, cfg := setupWorkoutReportHandler(reportRepo, workoutRepo)

	requestBody := dto.CreateWorkoutReportRequest{
		WorkoutId: 1,
		Details:   "Great workout session! Completed all sets with good form.",
	}

	jsonBody, _ := json.Marshal(requestBody)
	c, w := createAuthenticatedGinContext("POST", "/v1/workouts/workout-report/", jsonBody, tokenProvider, cfg)

	// Set up the route and call the handler
	c.Request.URL.Path = "/v1/workouts/workout-report/"
	handler.Create(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	var response helper.BaseHttpResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, false, response.Success)
}

func TestGetWorkoutReportById_Handler_Success(t *testing.T) {
	reportRepo := &MockWorkoutReportRepository{}
	workoutRepo := &MockWorkoutRepository{}
	handler, tokenProvider, cfg := setupWorkoutReportHandler(reportRepo, workoutRepo)

	// Use the helper function with parameters
	params := gin.Params{{Key: "id", Value: "1"}}
	c, w := createAuthenticatedGinContextWithParams("GET", "/v1/workouts/workout-report/1", nil, params, tokenProvider, cfg)

	// Set up the route and call the handler
	c.Request.URL.Path = "/v1/workouts/workout-report/1"
	handler.GetById(c)
	assert.Equal(t, http.StatusOK, w.Code)
	var response helper.BaseHttpResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, true, response.Success)
}

func TestGetWorkoutReportById_Handler_InvalidId(t *testing.T) {
	reportRepo := &MockWorkoutReportRepository{}
	workoutRepo := &MockWorkoutRepository{}
	handler, tokenProvider, cfg := setupWorkoutReportHandler(reportRepo, workoutRepo)

	// Use the helper function with parameters
	params := gin.Params{{Key: "id", Value: "invalid"}}
	c, w := createAuthenticatedGinContextWithParams("GET", "/v1/workouts/workout-report/invalid", nil, params, tokenProvider, cfg)

	// Set up the route and call the handler
	c.Request.URL.Path = "/v1/workouts/workout-report/invalid"
	handler.GetById(c)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	var response helper.BaseHttpResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, false, response.Success)
	assert.Equal(t, helper.ValidationError, response.ResultCode)
}

func TestGetWorkoutReportById_Handler_NotFound(t *testing.T) {
	reportRepo := &MockWorkoutReportRepository{
		GetByIdFn: func(ctx context.Context, id int) (models.WorkoutReport, error) {
			return models.WorkoutReport{}, &service_errors.ServiceError{EndUserMessage: service_errors.RecordNotFound}
		},
	}
	workoutRepo := &MockWorkoutRepository{}
	handler, tokenProvider, cfg := setupWorkoutReportHandler(reportRepo, workoutRepo)

	// Use the helper function with parameters
	params := gin.Params{{Key: "id", Value: "1"}}
	c, w := createAuthenticatedGinContextWithParams("GET", "/v1/workouts/workout-report/1", nil, params, tokenProvider, cfg)

	// Set up the route and call the handler
	c.Request.URL.Path = "/v1/workouts/workout-report/1"
	handler.GetById(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
	var response helper.BaseHttpResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, false, response.Success)
	assert.Equal(t, "record not found", response.Error)
}

func TestUpdateWorkoutReport_Handler_Success(t *testing.T) {
	reportRepo := &MockWorkoutReportRepository{}
	workoutRepo := &MockWorkoutRepository{}
	handler, tokenProvider, cfg := setupWorkoutReportHandler(reportRepo, workoutRepo)

	// Use the helper function with parameters
	params := gin.Params{{Key: "id", Value: "1"}}
	requestBody := dto.UpdateWorkoutReportRequest{
		WorkoutId: 1,
		Details:   "Updated workout report with additional notes about performance.",
	}
	jsonBody, _ := json.Marshal(requestBody)
	c, w := createAuthenticatedGinContextWithParams("PUT", "/v1/workouts/workout-report/1", jsonBody, params, tokenProvider, cfg)

	// Set up the route and call the handler
	c.Request.URL.Path = "/v1/workouts/workout-report/1"
	handler.Update(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var response helper.BaseHttpResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, true, response.Success)
}

func TestUpdateWorkoutReport_Handler_InvalidId(t *testing.T) {
	reportRepo := &MockWorkoutReportRepository{}
	workoutRepo := &MockWorkoutRepository{}
	handler, tokenProvider, cfg := setupWorkoutReportHandler(reportRepo, workoutRepo)

	// Use the helper function with parameters
	params := gin.Params{{Key: "id", Value: "invalid"}}
	requestBody := dto.UpdateWorkoutReportRequest{
		WorkoutId: 1,
		Details:   "Updated workout report with additional notes about performance.",
	}
	jsonBody, _ := json.Marshal(requestBody)
	c, w := createAuthenticatedGinContextWithParams("PUT", "/v1/workouts/workout-report/invalid", jsonBody, params, tokenProvider, cfg)

	// Set up the route and call the handler
	c.Request.URL.Path = "/v1/workouts/workout-report/invalid"
	handler.Update(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response helper.BaseHttpResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, false, response.Success)
}

func TestDeleteWorkoutReport_Handler_InvalidId(t *testing.T) {
	reportRepo := &MockWorkoutReportRepository{}
	workoutRepo := &MockWorkoutRepository{}
	handler, tokenProvider, cfg := setupWorkoutReportHandler(reportRepo, workoutRepo)

	// Use the helper function with parameters
	params := gin.Params{{Key: "id", Value: "invalid"}}
	c, w := createAuthenticatedGinContextWithParams("DELETE", "/v1/workouts/workout-report/invalid", nil, params, tokenProvider, cfg)

	// Set up the route and call the handler
	c.Request.URL.Path = "/v1/workouts/workout-report/invalid"
	handler.Delete(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response helper.BaseHttpResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, false, response.Success)
}

func TestDeleteWorkoutReport_Handler_Success(t *testing.T) {
	reportRepo := &MockWorkoutReportRepository{}
	workoutRepo := &MockWorkoutRepository{}
	handler, tokenProvider, cfg := setupWorkoutReportHandler(reportRepo, workoutRepo)

	// Use the helper function with parameters
	params := gin.Params{{Key: "id", Value: "1"}}
	c, w := createAuthenticatedGinContextWithParams("DELETE", "/v1/workouts/workout-report/1", nil, params, tokenProvider, cfg)

	// Set up the route and call the handler
	c.Request.URL.Path = "/v1/workouts/workout-report/1"
	handler.Delete(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var response helper.BaseHttpResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, true, response.Success)
}

func TestDeleteWorkoutReport_Handler_NotFound(t *testing.T) {
	reportRepo := &MockWorkoutReportRepository{
		DeleteFn: func(ctx context.Context, id int) error {
			return &service_errors.ServiceError{EndUserMessage: service_errors.RecordNotFound}
		},
	}
	workoutRepo := &MockWorkoutRepository{}
	handler, tokenProvider, cfg := setupWorkoutReportHandler(reportRepo, workoutRepo)

	// Use the helper function with parameters
	params := gin.Params{{Key: "id", Value: "7"}}
	c, w := createAuthenticatedGinContextWithParams("DELETE", "/v1/workouts/workout-report/7", nil, params, tokenProvider, cfg)

	// Set up the route and call the handler
	c.Request.URL.Path = "/v1/workouts/workout-report/7"
	handler.Delete(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
	var response helper.BaseHttpResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, false, response.Success)
	assert.Equal(t, service_errors.RecordNotFound, response.Error)
}
