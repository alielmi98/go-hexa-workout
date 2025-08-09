package test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/alielmi98/go-hexa-workout/constants"
	"github.com/alielmi98/go-hexa-workout/internal/middlewares"
	"github.com/alielmi98/go-hexa-workout/internal/user/adapter/http/dto"
	"github.com/alielmi98/go-hexa-workout/internal/user/entity"
	"github.com/alielmi98/go-hexa-workout/internal/workout/core/models"
	"github.com/alielmi98/go-hexa-workout/internal/workout/core/usecase"
	"github.com/alielmi98/go-hexa-workout/internal/workout/port/filter"
	"github.com/alielmi98/go-hexa-workout/pkg/config"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

// MockWorkoutRepository implements WorkoutRepository interface for testing
type MockWorkoutRepository struct {
	CreateFn      func(ctx context.Context, entity models.Workout) (models.Workout, error)
	UpdateFn      func(ctx context.Context, id int, entity models.Workout) (models.Workout, error)
	DeleteFn      func(ctx context.Context, id int) error
	GetByIdFn     func(ctx context.Context, id int) (models.Workout, error)
	GetByFilterFn func(ctx context.Context, req filter.PaginationInputWithFilter) (int64, *[]models.Workout, error)
}

func (m *MockWorkoutRepository) Create(ctx context.Context, entity models.Workout) (models.Workout, error) {
	if m.CreateFn != nil {
		return m.CreateFn(ctx, entity)
	}
	entity.Id = 1
	entity.CreatedAt = time.Now()
	return entity, nil
}

func (m *MockWorkoutRepository) Update(ctx context.Context, id int, entity models.Workout) (models.Workout, error) {
	if m.UpdateFn != nil {
		return m.UpdateFn(ctx, id, entity)
	}
	entity.Id = id
	return entity, nil
}

func (m *MockWorkoutRepository) Delete(ctx context.Context, id int) error {
	if m.DeleteFn != nil {
		return m.DeleteFn(ctx, id)
	}
	return nil
}

func (m *MockWorkoutRepository) GetById(ctx context.Context, id int) (models.Workout, error) {
	if m.GetByIdFn != nil {
		return m.GetByIdFn(ctx, id)
	}
	return models.Workout{
		Id:          id,
		UserId:      1,
		Name:        "Test Workout",
		Description: "Test Description",
		Comments:    "Test Comments",
		CreatedAt:   time.Now(),
	}, nil
}

func (m *MockWorkoutRepository) GetByFilter(ctx context.Context, req filter.PaginationInputWithFilter) (int64, *[]models.Workout, error) {
	if m.GetByFilterFn != nil {
		return m.GetByFilterFn(ctx, req)
	}
	workouts := []models.Workout{
		{
			Id:          1,
			UserId:      1,
			Name:        "Test Workout",
			Description: "Test Description",
			Comments:    "Test Comments",
			CreatedAt:   time.Now(),
		},
	}
	return 1, &workouts, nil
}

// MockScheduledWorkoutsRepository implements ScheduledWorkoutsRepository interface for testing
type MockScheduledWorkoutsRepository struct {
	CreateFn      func(ctx context.Context, entity models.ScheduledWorkouts) (models.ScheduledWorkouts, error)
	UpdateFn      func(ctx context.Context, id int, entity models.ScheduledWorkouts) (models.ScheduledWorkouts, error)
	DeleteFn      func(ctx context.Context, id int) error
	GetByIdFn     func(ctx context.Context, id int) (models.ScheduledWorkouts, error)
	GetByFilterFn func(ctx context.Context, req filter.PaginationInputWithFilter) (int64, *[]models.ScheduledWorkouts, error)
}

func (m *MockScheduledWorkoutsRepository) Create(ctx context.Context, entity models.ScheduledWorkouts) (models.ScheduledWorkouts, error) {
	if m.CreateFn != nil {
		return m.CreateFn(ctx, entity)
	}
	entity.Id = 1
	entity.CreatedAt = time.Now()
	return entity, nil
}

func (m *MockScheduledWorkoutsRepository) Update(ctx context.Context, id int, entity models.ScheduledWorkouts) (models.ScheduledWorkouts, error) {
	if m.UpdateFn != nil {
		return m.UpdateFn(ctx, id, entity)
	}
	entity.ScheduledTime = time.Now()
	entity.Id = id

	return entity, nil
}

func (m *MockScheduledWorkoutsRepository) Delete(ctx context.Context, id int) error {
	if m.DeleteFn != nil {
		return m.DeleteFn(ctx, id)
	}
	return nil
}

func (m *MockScheduledWorkoutsRepository) GetById(ctx context.Context, id int) (models.ScheduledWorkouts, error) {
	if m.GetByIdFn != nil {
		return m.GetByIdFn(ctx, id)
	}
	return models.ScheduledWorkouts{
		Id:            id,
		WorkoutId:     1,
		ScheduledTime: time.Now(),
		Status:        "active",
		CreatedAt:     time.Now(),
	}, nil
}

func (m *MockScheduledWorkoutsRepository) GetByFilter(ctx context.Context, req filter.PaginationInputWithFilter) (int64, *[]models.ScheduledWorkouts, error) {
	if m.GetByFilterFn != nil {
		return m.GetByFilterFn(ctx, req)
	}
	scheduledWorkouts := []models.ScheduledWorkouts{
		{
			Id:            1,
			WorkoutId:     1,
			ScheduledTime: time.Now(),
			Status:        "active",
			CreatedAt:     time.Now(),
		},
	}
	return 1, &scheduledWorkouts, nil
}

// MockWorkoutExerciseRepository implements WorkoutExerciseRepository interface for testing
type MockWorkoutExerciseRepository struct {
	CreateFn      func(ctx context.Context, entity models.WorkoutExercise) (models.WorkoutExercise, error)
	UpdateFn      func(ctx context.Context, id int, entity models.WorkoutExercise) (models.WorkoutExercise, error)
	DeleteFn      func(ctx context.Context, id int) error
	GetByIdFn     func(ctx context.Context, id int) (models.WorkoutExercise, error)
	GetByFilterFn func(ctx context.Context, req filter.PaginationInputWithFilter) (int64, *[]models.WorkoutExercise, error)
}

func (m *MockWorkoutExerciseRepository) Create(ctx context.Context, entity models.WorkoutExercise) (models.WorkoutExercise, error) {
	if m.CreateFn != nil {
		return m.CreateFn(ctx, entity)
	}
	entity.Id = 1
	entity.CreatedAt = time.Now()
	return entity, nil
}

func (m *MockWorkoutExerciseRepository) Update(ctx context.Context, id int, entity models.WorkoutExercise) (models.WorkoutExercise, error) {
	if m.UpdateFn != nil {
		return m.UpdateFn(ctx, id, entity)
	}
	entity.Id = id
	return entity, nil
}

func (m *MockWorkoutExerciseRepository) Delete(ctx context.Context, id int) error {
	if m.DeleteFn != nil {
		return m.DeleteFn(ctx, id)
	}
	return nil
}

func (m *MockWorkoutExerciseRepository) GetById(ctx context.Context, id int) (models.WorkoutExercise, error) {
	if m.GetByIdFn != nil {
		return m.GetByIdFn(ctx, id)
	}
	return models.WorkoutExercise{
		Id:          id,
		WorkoutId:   1,
		Name:        "Test Exercise",
		Description: "Test Description",
		Repetitions: 10,
		Sets:        3,
		Weight:      50.0,
		CreatedAt:   time.Now(),
	}, nil
}

func (m *MockWorkoutExerciseRepository) GetByFilter(ctx context.Context, req filter.PaginationInputWithFilter) (int64, *[]models.WorkoutExercise, error) {
	if m.GetByFilterFn != nil {
		return m.GetByFilterFn(ctx, req)
	}
	exercises := []models.WorkoutExercise{
		{
			Id:          1,
			WorkoutId:   1,
			Name:        "Test Exercise",
			Description: "Test Description",
			Repetitions: 10,
			Sets:        3,
			Weight:      50.0,
			CreatedAt:   time.Now(),
		},
	}
	return 1, &exercises, nil
}

// MockWorkoutReportRepository implements WorkoutReportRepository interface for testing
type MockWorkoutReportRepository struct {
	CreateFn      func(ctx context.Context, entity models.WorkoutReport) (models.WorkoutReport, error)
	UpdateFn      func(ctx context.Context, id int, entity models.WorkoutReport) (models.WorkoutReport, error)
	DeleteFn      func(ctx context.Context, id int) error
	GetByIdFn     func(ctx context.Context, id int) (models.WorkoutReport, error)
	GetByFilterFn func(ctx context.Context, req filter.PaginationInputWithFilter) (int64, *[]models.WorkoutReport, error)
}

func (m *MockWorkoutReportRepository) Create(ctx context.Context, entity models.WorkoutReport) (models.WorkoutReport, error) {
	if m.CreateFn != nil {
		return m.CreateFn(ctx, entity)
	}
	entity.Id = 1
	entity.CreatedAt = time.Now()
	return entity, nil
}

func (m *MockWorkoutReportRepository) Update(ctx context.Context, id int, entity models.WorkoutReport) (models.WorkoutReport, error) {
	if m.UpdateFn != nil {
		return m.UpdateFn(ctx, id, entity)
	}
	entity.Id = id
	return entity, nil
}

func (m *MockWorkoutReportRepository) Delete(ctx context.Context, id int) error {
	if m.DeleteFn != nil {
		return m.DeleteFn(ctx, id)
	}
	return nil
}

func (m *MockWorkoutReportRepository) GetById(ctx context.Context, id int) (models.WorkoutReport, error) {
	if m.GetByIdFn != nil {
		return m.GetByIdFn(ctx, id)
	}
	return models.WorkoutReport{
		Id:        id,
		WorkoutId: 1,
		UserId:    1,
		Details:   "Test Report Details",
		CreatedAt: time.Now(),
	}, nil
}

func (m *MockWorkoutReportRepository) GetByFilter(ctx context.Context, req filter.PaginationInputWithFilter) (int64, *[]models.WorkoutReport, error) {
	if m.GetByFilterFn != nil {
		return m.GetByFilterFn(ctx, req)
	}
	reports := []models.WorkoutReport{
		{
			Id:        1,
			WorkoutId: 1,
			UserId:    1,
			Details:   "Test Report Details",
			CreatedAt: time.Now(),
		},
	}
	return 1, &reports, nil
}

// Helper function to create context with user ID
func createContextWithUserId(userId float64) context.Context {
	ctx := context.Background()
	return context.WithValue(ctx, constants.UserIdKey, userId)
}

// Helper functions to setup use cases for testing
func setupWorkoutUsecase(workoutRepo *MockWorkoutRepository) *usecase.WorkoutUsecase {
	cfg := &config.Config{}
	return usecase.NewWorkoutUsecase(cfg, workoutRepo)
}

func setupScheduledWorkoutUsecase(scheduledRepo *MockScheduledWorkoutsRepository, workoutRepo *MockWorkoutRepository) *usecase.ScheduledWorkoutsUseCase {
	cfg := &config.Config{}
	return usecase.NewScheduledWorkoutsUsecase(cfg, scheduledRepo, workoutRepo)
}

func setupWorkoutExerciseUsecase(exerciseRepo *MockWorkoutExerciseRepository, workoutRepo *MockWorkoutRepository) *usecase.WorkoutExerciseUsecase {
	cfg := &config.Config{}
	return usecase.NewWorkoutExerciseUsecase(cfg, exerciseRepo, workoutRepo)
}

func setupWorkoutReportUsecase(reportRepo *MockWorkoutReportRepository, workoutRepo *MockWorkoutRepository) *usecase.WorkoutReportUsecase {
	cfg := &config.Config{}
	return usecase.NewWorkoutReportUsecase(cfg, reportRepo, workoutRepo)
}

// MockTokenProvider implements TokenProvider interface for testing
type MockTokenProvider struct {
	GenerateTokenFn func(token *entity.TokenPayload) (*dto.TokenDetail, error)
	VerifyTokenFn   func(token string) (*jwt.Token, error)
	GetClaimsFn     func(token string) (map[string]interface{}, error)
	RefreshTokenFn  func(refreshToken string) (*dto.TokenDetail, error)
}

func (m *MockTokenProvider) GenerateToken(token *entity.TokenPayload) (*dto.TokenDetail, error) {
	if m.GenerateTokenFn != nil {
		return m.GenerateTokenFn(token)
	}
	return &dto.TokenDetail{
		AccessToken:            "mock-access-token",
		RefreshToken:           "mock-refresh-token",
		AccessTokenExpireTime:  time.Now().Add(time.Hour).Unix(),
		RefreshTokenExpireTime: time.Now().Add(24 * time.Hour).Unix(),
	}, nil
}

func (m *MockTokenProvider) VerifyToken(token string) (*jwt.Token, error) {
	if m.VerifyTokenFn != nil {
		return m.VerifyTokenFn(token)
	}
	return &jwt.Token{Valid: true}, nil
}

func (m *MockTokenProvider) GetClaims(token string) (map[string]interface{}, error) {
	if m.GetClaimsFn != nil {
		return m.GetClaimsFn(token)
	}
	// Return mock claims with user information
	return map[string]interface{}{
		constants.UserIdKey:       float64(1),
		constants.FirstNameKey:    "Test",
		constants.LastNameKey:     "User",
		constants.UsernameKey:     "testuser",
		constants.EmailKey:        "test@example.com",
		constants.MobileNumberKey: "1234567890",
		constants.ExpireTimeKey:   float64(time.Now().Add(time.Hour).Unix()),
	}, nil
}

func (m *MockTokenProvider) RefreshToken(refreshToken string) (*dto.TokenDetail, error) {
	if m.RefreshTokenFn != nil {
		return m.RefreshTokenFn(refreshToken)
	}
	return &dto.TokenDetail{
		AccessToken:            "mock-new-access-token",
		RefreshToken:           "mock-new-refresh-token",
		AccessTokenExpireTime:  time.Now().Add(time.Hour).Unix(),
		RefreshTokenExpireTime: time.Now().Add(24 * time.Hour).Unix(),
	}, nil
}

// Helper function to generate a valid JWT token for testing
func generateTestToken() string {
	return "Bearer mock-jwt-token"
}

// Helper function to create authenticated test context
func createAuthenticatedContext(userId int) context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, constants.UserIdKey, float64(userId))
	return ctx
}

// Helper function to create authenticated gin context for testing
func createAuthenticatedGinContext(method, url string, body []byte, tokenProvider *MockTokenProvider, cfg *config.Config) (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)

	// Setup the route with authentication middleware
	r.Use(middlewares.Authentication(cfg, tokenProvider))

	// Create request with proper authorization header
	req, _ := http.NewRequest(method, url, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(constants.AuthorizationHeaderKey, generateTestToken())

	c.Set(constants.UserIdKey, float64(1))
	c.Request = req
	return c, w
}

// Helper function to create authenticated gin context with URL parameters
func createAuthenticatedGinContextWithParams(method, url string, body []byte, params gin.Params, tokenProvider *MockTokenProvider, cfg *config.Config) (*gin.Context, *httptest.ResponseRecorder) {
	c, w := createAuthenticatedGinContext(method, url, body, tokenProvider, cfg)
	c.Params = params
	return c, w
}
