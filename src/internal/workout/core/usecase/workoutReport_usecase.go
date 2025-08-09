package usecase

import (
	"context"

	"github.com/alielmi98/go-hexa-workout/constants"
	"github.com/alielmi98/go-hexa-workout/internal/workout/core/models"
	"github.com/alielmi98/go-hexa-workout/internal/workout/core/usecase/dto"
	"github.com/alielmi98/go-hexa-workout/internal/workout/port"
	"github.com/alielmi98/go-hexa-workout/pkg/config"
)

type WorkoutReportUsecase struct {
	base        *BaseUsecase[models.WorkoutReport, dto.CreateWorkoutReportRequest, dto.UpdateWorkoutReportRequest, dto.WorkoutReportResponse]
	workoutRepo port.WorkoutRepository
}

func NewWorkoutReportUsecase(cfg *config.Config, workoutReportRepository port.WorkoutReportRepository, workoutRepository port.WorkoutRepository) *WorkoutReportUsecase {
	return &WorkoutReportUsecase{
		base:        NewBaseUsecase[models.WorkoutReport, dto.CreateWorkoutReportRequest, dto.UpdateWorkoutReportRequest, dto.WorkoutReportResponse](cfg, workoutReportRepository),
		workoutRepo: workoutRepository,
	}
}

func (u *WorkoutReportUsecase) Create(ctx context.Context, req dto.CreateWorkoutReportRequest) (dto.WorkoutReportResponse, error) {
	// Check if the user is Owner of the Workout
	err := u.base.CheckOwnership(ctx, u.workoutRepo, req.WorkoutId)
	if err != nil {
		return dto.WorkoutReportResponse{}, err
	}

	userId := int(ctx.Value(constants.UserIdKey).(float64))
	req.UserID = userId

	return u.base.Create(ctx, req)
}

func (u *WorkoutReportUsecase) Update(ctx context.Context, id int, req dto.UpdateWorkoutReportRequest) (dto.WorkoutReportResponse, error) {
	// Check if the user is Owner of the Workout
	workoutReport, err := u.base.GetById(ctx, id)
	if err != nil {
		return dto.WorkoutReportResponse{}, err
	}
	err = u.base.CheckOwnership(ctx, u.workoutRepo, workoutReport.WorkoutId)
	if err != nil {
		return dto.WorkoutReportResponse{}, err
	}
	// Check if the user is Owner of the Workout to update the report
	err = u.base.CheckOwnership(ctx, u.workoutRepo, req.WorkoutId)
	if err != nil {
		return dto.WorkoutReportResponse{}, err
	}

	return u.base.Update(ctx, id, req)
}

func (u *WorkoutReportUsecase) Delete(ctx context.Context, id int) error {
	// Check if the user is Owner of the Workout
	workoutReport, err := u.base.GetById(ctx, id)
	if err != nil {
		return err
	}
	err = u.base.CheckOwnership(ctx, u.workoutRepo, workoutReport.WorkoutId)
	if err != nil {
		return err
	}

	return u.base.Delete(ctx, id)
}

func (u *WorkoutReportUsecase) GetById(ctx context.Context, id int) (dto.WorkoutReportResponse, error) {
	// Check if the user is Owner of the Workout
	workoutReport, err := u.base.GetById(ctx, id)
	if err != nil {
		return dto.WorkoutReportResponse{}, err
	}
	err = u.base.CheckOwnership(ctx, u.workoutRepo, workoutReport.WorkoutId)
	if err != nil {
		return dto.WorkoutReportResponse{}, err
	}

	return workoutReport, nil
}
