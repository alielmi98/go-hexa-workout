package usecase

import (
	"context"

	"github.com/alielmi98/go-hexa-workout/internal/workout/core/models"
	"github.com/alielmi98/go-hexa-workout/internal/workout/core/usecase/dto"
	"github.com/alielmi98/go-hexa-workout/internal/workout/port"
	"github.com/alielmi98/go-hexa-workout/pkg/config"
	"github.com/alielmi98/go-hexa-workout/pkg/service_errors"
)

type ScheduledWorkoutsUseCase struct {
	base        *BaseUsecase[models.ScheduledWorkouts, dto.CreateScheduledWorkoutsRequest, dto.UpdateScheduledWorkoutsRequest, dto.ScheduledWorkoutsResponse]
	workoutRepo port.WorkoutRepository
}

func NewScheduledWorkoutsUsecase(cfg *config.Config, ScheduledWorkoutsRepository port.ScheduledWorkoutsRepository, workoutRepository port.WorkoutRepository) *ScheduledWorkoutsUseCase {
	return &ScheduledWorkoutsUseCase{
		base:        NewBaseUsecase[models.ScheduledWorkouts, dto.CreateScheduledWorkoutsRequest, dto.UpdateScheduledWorkoutsRequest, dto.ScheduledWorkoutsResponse](cfg, ScheduledWorkoutsRepository),
		workoutRepo: workoutRepository,
	}
}

func (u *ScheduledWorkoutsUseCase) Create(ctx context.Context, req dto.CreateScheduledWorkoutsRequest) (dto.ScheduledWorkoutsResponse, error) {
	// Check if the user is Owner of the Workout
	err := u.base.CheckOwnership(ctx, u.workoutRepo, req.WorkoutId)
	if err != nil {
		return dto.ScheduledWorkoutsResponse{}, err
	}
	if req.Status != "active" && req.Status != "completed" && req.Status != "cancelled" {
		return dto.ScheduledWorkoutsResponse{}, &service_errors.ServiceError{EndUserMessage: service_errors.InvalidStatus}
	}

	return u.base.Create(ctx, req)
}

func (u *ScheduledWorkoutsUseCase) Update(ctx context.Context, id int, req dto.UpdateScheduledWorkoutsRequest) (dto.ScheduledWorkoutsResponse, error) {
	// Check if the user is Owner of the Workout
	ScheduledWorkouts, err := u.base.GetById(ctx, id)
	if err != nil {
		return dto.ScheduledWorkoutsResponse{}, err
	}
	err = u.base.CheckOwnership(ctx, u.workoutRepo, ScheduledWorkouts.WorkoutId)
	if err != nil {
		return dto.ScheduledWorkoutsResponse{}, err
	}

	if req.Status != "active" && req.Status != "completed" && req.Status != "cancelled" {
		return dto.ScheduledWorkoutsResponse{}, &service_errors.ServiceError{EndUserMessage: service_errors.InvalidStatus}
	}

	return u.base.Update(ctx, id, req)
}

func (u *ScheduledWorkoutsUseCase) Delete(ctx context.Context, id int) error {
	// Check if the user is Owner of the Workout
	ScheduledWorkouts, err := u.base.GetById(ctx, id)
	if err != nil {
		return err
	}
	err = u.base.CheckOwnership(ctx, u.workoutRepo, ScheduledWorkouts.WorkoutId)
	if err != nil {
		return err
	}
	return u.base.Delete(ctx, id)
}

func (u *ScheduledWorkoutsUseCase) GetById(ctx context.Context, id int) (dto.ScheduledWorkoutsResponse, error) {
	// Check if the user is Owner of the Workout
	ScheduledWorkouts, err := u.base.GetById(ctx, id)
	if err != nil {
		return dto.ScheduledWorkoutsResponse{}, err
	}

	err = u.base.CheckOwnership(ctx, u.workoutRepo, ScheduledWorkouts.WorkoutId)
	if err != nil {
		return dto.ScheduledWorkoutsResponse{}, err
	}

	return ScheduledWorkouts, nil
}
