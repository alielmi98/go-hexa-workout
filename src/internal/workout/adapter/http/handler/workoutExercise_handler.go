package handler

import (
	"github.com/alielmi98/go-hexa-workout/dependency"
	"github.com/alielmi98/go-hexa-workout/internal/workout/adapter/http/dto"
	_ "github.com/alielmi98/go-hexa-workout/internal/workout/adapter/http/dto"
	"github.com/alielmi98/go-hexa-workout/internal/workout/core/usecase"
	_ "github.com/alielmi98/go-hexa-workout/internal/workout/port/filter"
	"github.com/alielmi98/go-hexa-workout/pkg/config"
	_ "github.com/alielmi98/go-hexa-workout/pkg/helper"
	"github.com/gin-gonic/gin"
)

type WorkoutExerciseHandler struct {
	usecase *usecase.WorkoutExerciseUsecase
}

func NewWorkoutExerciseHandler(cfg *config.Config) *WorkoutExerciseHandler {
	workoutExerciseRepo, workoutRepo := dependency.GetWorkoutExerciseRepository()
	return &WorkoutExerciseHandler{
		usecase: usecase.NewWorkoutExerciseUsecase(cfg, workoutExerciseRepo, workoutRepo),
	}
}

// CreateWorkoutExercise godoc
// @Summary Create a WorkoutExercise
// @Description Create a WorkoutExercise
// @Tags WorkoutExercise
// @Accept json
// @produces json
// @Param Request body dto.CreateWorkoutExerciseRequest true "Create a WorkoutExercise"
// @Success 201 {object} helper.BaseHttpResponse{result=dto.WorkoutExerciseResponse} "WorkoutExercise response"
// @Failure 400 {object} helper.BaseHttpResponse "Bad request"
// @Router /v1/workouts/workout-exercise/ [post]
// @Security AuthBearer
func (h *WorkoutExerciseHandler) Create(c *gin.Context) {
	Create(c, dto.ToCreateWorkoutExerciseRequest, dto.ToWorkoutExerciseResponse, h.usecase.Create)
}

// GetWorkoutExercises godoc
// @Summary Get a WorkoutExercise by ID
// @Description Get a WorkoutExercise by ID
// @Tags WorkoutExercise
// @Accept json
// @Produce json
// @Param id path int true "WorkoutExercise ID"
// @Success 200 {object} helper.BaseHttpResponse{result=dto.WorkoutExerciseResponse} "WorkoutExercise response"
// @Failure 400 {object} helper.BaseHttpResponse "Bad request"
// @Failure 404 {object} helper.BaseHttpResponse "Not found"
// @Router /v1/workouts/workout-exercise/{id} [get]
// @Security AuthBearer
func (h *WorkoutExerciseHandler) GetById(c *gin.Context) {
	GetById(c, dto.ToWorkoutExerciseResponse, h.usecase.GetById)
}

// UpdateWorkoutExercise godoc
// @Summary Update a WorkoutExercise
// @Description Update a WorkoutExercise
// @Tags WorkoutExercise
// @Accept json
// @Produce json
// @Param id path int true "WorkoutExercise ID"
// @Param Request body dto.UpdateWorkoutExerciseRequest true "Update a WorkoutExercise"
// @Success 200 {object} helper.BaseHttpResponse{result=dto.WorkoutExerciseResponse} "WorkoutExercise response"
// @Failure 400 {object} helper.BaseHttpResponse "Bad request"
// @Failure 404 {object} helper.BaseHttpResponse "Not found"
// @Router /v1/workouts/workout-exercise/{id} [put]
// @Security AuthBearer
func (h *WorkoutExerciseHandler) Update(c *gin.Context) {
	Update(c, dto.ToUpdateWorkoutExerciseRequest, dto.ToWorkoutExerciseResponse, h.usecase.Update)
}

// DeleteWorkoutExercise godoc
// @Summary Delete a WorkoutExercise
// @Description Delete a WorkoutExercise
// @Tags WorkoutExercise
// @Param id path int true "WorkoutExercise ID"
// @Success 204 {object} helper.BaseHttpResponse "No Content"
// @Failure 400 {object} helper.BaseHttpResponse "Bad request"
// @Failure 404 {object} helper.BaseHttpResponse "Not found"
// @Router /v1/workouts/workout-exercise/{id} [delete]
// @Security AuthBearer
func (h *WorkoutExerciseHandler) Delete(c *gin.Context) {
	Delete(c, h.usecase.Delete)
}
