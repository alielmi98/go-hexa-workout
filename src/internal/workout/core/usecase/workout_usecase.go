package usecase

import (
	"context"
	"fmt"

	"github.com/alielmi98/go-hexa-workout/constants"
	"github.com/alielmi98/go-hexa-workout/internal/workout/core/models"
	"github.com/alielmi98/go-hexa-workout/internal/workout/core/usecase/dto"
	"github.com/alielmi98/go-hexa-workout/internal/workout/port"
	"github.com/alielmi98/go-hexa-workout/internal/workout/port/filter"
	"github.com/alielmi98/go-hexa-workout/pkg/config"
	"github.com/alielmi98/go-hexa-workout/pkg/service_errors"
)

type WorkoutUsecase struct {
	base *BaseUsecase[models.Workout, dto.CreateWorkoutRequest, dto.UpdateWorkoutRequest, dto.WorkoutResponse]
}

func NewWorkoutUsecase(cfg *config.Config, workoutRepository port.WorkoutRepository) *WorkoutUsecase {
	return &WorkoutUsecase{
		base: NewBaseUsecase[models.Workout, dto.CreateWorkoutRequest, dto.UpdateWorkoutRequest, dto.WorkoutResponse](cfg, workoutRepository),
	}
}

func (u *WorkoutUsecase) Create(ctx context.Context, req dto.CreateWorkoutRequest) (dto.WorkoutResponse, error) {
	userId := int(ctx.Value(constants.UserIdKey).(float64))
	req.UserId = userId
	return u.base.Create(ctx, req)
}

func (u *WorkoutUsecase) Update(ctx context.Context, id int, req dto.UpdateWorkoutRequest) (dto.WorkoutResponse, error) {
	// Check if the user is Owner of the Workout
	userId := int(ctx.Value(constants.UserIdKey).(float64))
	workout, err := u.base.GetById(ctx, id)
	if err != nil {
		return dto.WorkoutResponse{}, err
	}

	if workout.UserId != userId {
		return dto.WorkoutResponse{}, &service_errors.ServiceError{EndUserMessage: service_errors.PermissionDenied}
	}

	return u.base.Update(ctx, id, req)
}

func (u *WorkoutUsecase) Delete(ctx context.Context, id int) error {
	// Check if the user is Owner of the Workout
	userId := int(ctx.Value(constants.UserIdKey).(float64))
	workout, err := u.base.GetById(ctx, id)
	if err != nil {
		return err
	}

	if workout.UserId != userId {
		return &service_errors.ServiceError{EndUserMessage: service_errors.PermissionDenied}
	}
	return u.base.Delete(ctx, id)
}
func (u *WorkoutUsecase) GetById(ctx context.Context, id int) (dto.WorkoutResponse, error) {
	// Check if the user is Owner of the Workout
	userId := int(ctx.Value(constants.UserIdKey).(float64))
	workout, err := u.base.GetById(ctx, id)
	if err != nil {
		return dto.WorkoutResponse{}, err
	}

	if workout.UserId != userId {
		return dto.WorkoutResponse{}, &service_errors.ServiceError{EndUserMessage: service_errors.PermissionDenied}
	}
	return u.base.GetById(ctx, id)
}
func (u *WorkoutUsecase) GetByFilter(ctx context.Context, req filter.PaginationInputWithFilter) (*filter.PagedList[dto.WorkoutResponse], error) {
	// Add user filter to ensure users only see their own workouts
	userId := int(ctx.Value(constants.UserIdKey).(float64))

	// Add user_id filter to the existing filters
	if req.DynamicFilter.Filter == nil {
		req.DynamicFilter.Filter = make(map[string]filter.Filter)
	}

	// Add user_id as an equals filter
	req.DynamicFilter.Filter["UserId"] = filter.Filter{
		Type:       "equals",
		From:       fmt.Sprintf("%d", userId),
		FilterType: "number",
	}

	return u.base.GetByFilter(ctx, req)
}
