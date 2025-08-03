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

type WorkoutHandler struct {
	usecase *usecase.WorkoutUsecase
}

func NewWorkoutHandler(cfg *config.Config) *WorkoutHandler {
	return &WorkoutHandler{
		usecase: usecase.NewWorkoutUsecase(cfg, dependency.GetWorkoutRepository()),
	}
}

// CreateWorkout godoc
// @Summary Create a Workout
// @Description Create a Workout
// @Tags Workout
// @Accept json
// @produces json
// @Param Request body dto.CreateWorkoutRequest true "Create a Workout"
// @Success 201 {object} helper.BaseHttpResponse{result=dto.WorkoutResponse} "Workout response"
// @Failure 400 {object} helper.BaseHttpResponse "Bad request"
// @Router /v1/workouts/workout/ [post]
// @Security AuthBearer
func (h *WorkoutHandler) Create(c *gin.Context) {
	Create(c, dto.ToCreateWorkoutRequest, dto.ToWorkoutResponse, h.usecase.Create)
}

// GetWorkouts godoc
// @Summary Get a Workout by ID
// @Description Get a Workout by ID
// @Tags Workout
// @Accept json
// @Produce json
// @Param id path int true "Workout ID"
// @Success 200 {object} helper.BaseHttpResponse{result=dto.WorkoutResponse} "Workout response"
// @Failure 400 {object} helper.BaseHttpResponse "Bad request"
// @Failure 404 {object} helper.BaseHttpResponse "Not found"
// @Router /v1/workouts/workout/{id} [get]
// @Security AuthBearer
func (h *WorkoutHandler) GetById(c *gin.Context) {
	GetById(c, dto.ToWorkoutResponse, h.usecase.GetById)
}

// UpdateWorkout godoc
// @Summary Update a Workout
// @Description Update a Workout
// @Tags Workout
// @Accept json
// @Produce json
// @Param id path int true "Workout ID"
// @Param Request body dto.UpdateWorkoutRequest true "Update a Workout"
// @Success 200 {object} helper.BaseHttpResponse{result=dto.WorkoutResponse} "Workout response"
// @Failure 400 {object} helper.BaseHttpResponse "Bad request"
// @Failure 404 {object} helper.BaseHttpResponse "Not found"
// @Router /v1/workouts/workout/{id} [put]
// @Security AuthBearer
func (h *WorkoutHandler) Update(c *gin.Context) {
	Update(c, dto.ToUpdateWorkoutRequest, dto.ToWorkoutResponse, h.usecase.Update)
}

// DeleteWorkout godoc
// @Summary Delete a Workout
// @Description Delete a Workout
// @Tags Workout
// @Param id path int true "Workout ID"
// @Success 204 {object} helper.BaseHttpResponse "No Content"
// @Failure 400 {object} helper.BaseHttpResponse "Bad request"
// @Failure 404 {object} helper.BaseHttpResponse "Not found"
// @Router /v1/workouts/workout/{id} [delete]
// @Security AuthBearer
func (h *WorkoutHandler) Delete(c *gin.Context) {
	Delete(c, h.usecase.Delete)
}

// GetWorkoutsByFilter godoc
// @Summary Get Workouts by Filter
// @Description Get Workouts by Filter
// @Tags Workout
// @Accept json
// @Param Request body filter.PaginationInputWithFilter true "Request"
// @Success 200 {object} helper.BaseHttpResponse{result=filter.PagedList[dto.WorkoutResponse]} "Workout response"
// @Failure 400 {object} helper.BaseHttpResponse "Bad request"
// @Router /v1/workouts/workout/get-by-filter [post]
// @Security AuthBearer
func (h *WorkoutHandler) GetByFilter(c *gin.Context) {
	GetByFilter(c, dto.ToWorkoutResponse, h.usecase.GetByFilter)
}
