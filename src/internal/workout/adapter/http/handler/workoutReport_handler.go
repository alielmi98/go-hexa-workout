package handler

import (
	"github.com/alielmi98/go-hexa-workout/dependency"
	"github.com/alielmi98/go-hexa-workout/internal/workout/adapter/http/dto"
	"github.com/alielmi98/go-hexa-workout/internal/workout/core/usecase"
	_ "github.com/alielmi98/go-hexa-workout/internal/workout/port/filter"
	"github.com/alielmi98/go-hexa-workout/pkg/config"
	_ "github.com/alielmi98/go-hexa-workout/pkg/helper"
	"github.com/gin-gonic/gin"
)

type WorkoutReportHandler struct {
	Usecase *usecase.WorkoutReportUsecase
}

func NewWorkoutReportHandler(cfg *config.Config) *WorkoutReportHandler {
	return &WorkoutReportHandler{
		Usecase: usecase.NewWorkoutReportUsecase(cfg, dependency.GetWorkoutReportRepository(), dependency.GetWorkoutRepository()),
	}
}

// CreateWorkoutReport godoc
// @Summary Create a WorkoutReport
// @Description Create a WorkoutReport
// @Tags WorkoutReport
// @Accept json
// @produces json
// @Param Request body dto.CreateWorkoutReportRequest true "Create a WorkoutReport"
// @Success 201 {object} helper.BaseHttpResponse{result=dto.WorkoutReportResponse} "WorkoutReport response"
// @Failure 400 {object} helper.BaseHttpResponse "Bad request"
// @Router /v1/workouts/workout-report/ [post]
// @Security AuthBearer
func (h *WorkoutReportHandler) Create(c *gin.Context) {
	Create(c, dto.ToCreateWorkoutReportRequest, dto.ToWorkoutReportResponse, h.Usecase.Create)
}

// GetWorkoutReports godoc
// @Summary Get a WorkoutReport by ID
// @Description Get a WorkoutReport by ID
// @Tags WorkoutReport
// @Accept json
// @Produce json
// @Param id path int true "WorkoutReport ID"
// @Success 200 {object} helper.BaseHttpResponse{result=dto.WorkoutReportResponse} "WorkoutReport response"
// @Failure 400 {object} helper.BaseHttpResponse "Bad request"
// @Failure 404 {object} helper.BaseHttpResponse "Not found"
// @Router /v1/workouts/workout-report/{id} [get]
// @Security AuthBearer
func (h *WorkoutReportHandler) GetById(c *gin.Context) {
	GetById(c, dto.ToWorkoutReportResponse, h.Usecase.GetById)
}

// UpdateWorkoutReport godoc
// @Summary Update a WorkoutReport
// @Description Update a WorkoutReport
// @Tags WorkoutReport
// @Accept json
// @Produce json
// @Param id path int true "WorkoutReport ID"
// @Param Request body dto.UpdateWorkoutReportRequest true "Update a WorkoutReport"
// @Success 200 {object} helper.BaseHttpResponse{result=dto.WorkoutReportResponse} "WorkoutReport response"
// @Failure 400 {object} helper.BaseHttpResponse "Bad request"
// @Failure 404 {object} helper.BaseHttpResponse "Not found"
// @Router /v1/workouts/workout-report/{id} [put]
// @Security AuthBearer
func (h *WorkoutReportHandler) Update(c *gin.Context) {
	Update(c, dto.ToUpdateWorkoutReportRequest, dto.ToWorkoutReportResponse, h.Usecase.Update)
}

// DeleteWorkoutReport godoc
// @Summary Delete a WorkoutReport
// @Description Delete a WorkoutReport
// @Tags WorkoutReport
// @Param id path int true "WorkoutReport ID"
// @Success 204 {object} helper.BaseHttpResponse "No Content"
// @Failure 400 {object} helper.BaseHttpResponse "Bad request"
// @Failure 404 {object} helper.BaseHttpResponse "Not found"
// @Router /v1/workouts/workout-report/{id} [delete]
// @Security AuthBearer
func (h *WorkoutReportHandler) Delete(c *gin.Context) {
	Delete(c, h.Usecase.Delete)
}
