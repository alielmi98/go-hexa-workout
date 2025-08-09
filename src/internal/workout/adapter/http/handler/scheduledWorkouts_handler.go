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

type ScheduledWorkoutsHandler struct {
	Usecase *usecase.ScheduledWorkoutsUseCase
}

func NewScheduledWorkoutsHandler(cfg *config.Config) *ScheduledWorkoutsHandler {
	return &ScheduledWorkoutsHandler{
		Usecase: usecase.NewScheduledWorkoutsUsecase(cfg, dependency.GetScheduledWorkoutsRepository(), dependency.GetWorkoutRepository()),
	}
}

// CreateScheduledWorkouts godoc
// @Summary Create a ScheduledWorkouts
// @Description Create a ScheduledWorkouts
// @Tags ScheduledWorkouts
// @Accept json
// @produces json
// @Param Request body dto.CreateScheduledWorkoutsRequest true "Create a ScheduledWorkouts"
// @Success 201 {object} helper.BaseHttpResponse{result=dto.ScheduledWorkoutsResponse} "ScheduledWorkouts response"
// @Failure 400 {object} helper.BaseHttpResponse "Bad request"
// @Router /v1/workouts/scheduled-workouts/ [post]
// @Security AuthBearer
func (h *ScheduledWorkoutsHandler) Create(c *gin.Context) {
	Create(c, dto.ToCreateScheduledWorkoutsRequest, dto.ToScheduledWorkoutsResponse, h.Usecase.Create)
}

// GetScheduledWorkoutss godoc
// @Summary Get a ScheduledWorkouts by ID
// @Description Get a ScheduledWorkouts by ID
// @Tags ScheduledWorkouts
// @Accept json
// @Produce json
// @Param id path int true "ScheduledWorkouts ID"
// @Success 200 {object} helper.BaseHttpResponse{result=dto.ScheduledWorkoutsResponse} "ScheduledWorkouts response"
// @Failure 400 {object} helper.BaseHttpResponse "Bad request"
// @Failure 404 {object} helper.BaseHttpResponse "Not found"
// @Router /v1/workouts/scheduled-workouts/{id} [get]
// @Security AuthBearer
func (h *ScheduledWorkoutsHandler) GetById(c *gin.Context) {
	GetById(c, dto.ToScheduledWorkoutsResponse, h.Usecase.GetById)
}

// UpdateScheduledWorkouts godoc
// @Summary Update a ScheduledWorkouts
// @Description Update a ScheduledWorkouts
// @Tags ScheduledWorkouts
// @Accept json
// @Produce json
// @Param id path int true "ScheduledWorkouts ID"
// @Param Request body dto.UpdateScheduledWorkoutsRequest true "Update a ScheduledWorkouts"
// @Success 200 {object} helper.BaseHttpResponse{result=dto.ScheduledWorkoutsResponse} "ScheduledWorkouts response"
// @Failure 400 {object} helper.BaseHttpResponse "Bad request"
// @Failure 404 {object} helper.BaseHttpResponse "Not found"
// @Router /v1/workouts/scheduled-workouts/{id} [put]
// @Security AuthBearer
func (h *ScheduledWorkoutsHandler) Update(c *gin.Context) {
	Update(c, dto.ToUpdateScheduledWorkoutsRequest, dto.ToScheduledWorkoutsResponse, h.Usecase.Update)
}

// DeleteScheduledWorkouts godoc
// @Summary Delete a ScheduledWorkouts
// @Description Delete a ScheduledWorkouts
// @Tags ScheduledWorkouts
// @Param id path int true "ScheduledWorkouts ID"
// @Success 204 {object} helper.BaseHttpResponse "No Content"
// @Failure 400 {object} helper.BaseHttpResponse "Bad request"
// @Failure 404 {object} helper.BaseHttpResponse "Not found"
// @Router /v1/workouts/scheduled-workouts/{id} [delete]
// @Security AuthBearer
func (h *ScheduledWorkoutsHandler) Delete(c *gin.Context) {
	Delete(c, h.Usecase.Delete)
}
