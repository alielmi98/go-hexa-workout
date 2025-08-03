package usecase

import (
	"context"

	"github.com/alielmi98/go-hexa-workout/constants"
	"github.com/alielmi98/go-hexa-workout/internal/workout/core/models"
	"github.com/alielmi98/go-hexa-workout/internal/workout/core/usecase/dto"
	"github.com/alielmi98/go-hexa-workout/internal/workout/port"
	"github.com/alielmi98/go-hexa-workout/pkg/config"
	"github.com/alielmi98/go-hexa-workout/pkg/service_errors"
)

type WorkoutExerciseUsecase struct {
	base        *BaseUsecase[models.WorkoutExercise, dto.CreateWorkoutExerciseRequest, dto.UpdateWorkoutExerciseRequest, dto.WorkoutExerciseResponse]
	workoutRepo port.WorkoutRepository
}

func NewWorkoutExerciseUsecase(cfg *config.Config, workoutExerciseRepository port.WorkoutExerciseRepository, workoutRepository port.WorkoutRepository) *WorkoutExerciseUsecase {
	return &WorkoutExerciseUsecase{
		base:        NewBaseUsecase[models.WorkoutExercise, dto.CreateWorkoutExerciseRequest, dto.UpdateWorkoutExerciseRequest, dto.WorkoutExerciseResponse](cfg, workoutExerciseRepository),
		workoutRepo: workoutRepository,
	}
}

func (u *WorkoutExerciseUsecase) Create(ctx context.Context, req dto.CreateWorkoutExerciseRequest) (dto.WorkoutExerciseResponse, error) {
	// Check if the user is Owner of the Workout
	workout, err := u.workoutRepo.GetById(ctx, req.WorkoutId)
	if err != nil {
		return dto.WorkoutExerciseResponse{}, err
	}

	userId := int(ctx.Value(constants.UserIdKey).(float64))
	if userId != workout.UserId {
		return dto.WorkoutExerciseResponse{}, &service_errors.ServiceError{EndUserMessage: service_errors.PermissionDenied}
	}

	return u.base.Create(ctx, req)
}
func (u *WorkoutExerciseUsecase) Update(ctx context.Context, id int, req dto.UpdateWorkoutExerciseRequest) (dto.WorkoutExerciseResponse, error) {
	// Check if the user is Owner of the Workout
	workoutExercise, err := u.base.GetById(ctx, id)
	if err != nil {
		return dto.WorkoutExerciseResponse{}, err
	}
	userId := int(ctx.Value(constants.UserIdKey).(float64))
	workout, err := u.workoutRepo.GetById(ctx, workoutExercise.WorkoutId)
	if err != nil {
		return dto.WorkoutExerciseResponse{}, err
	}
	if userId != workout.UserId {
		return dto.WorkoutExerciseResponse{}, &service_errors.ServiceError{EndUserMessage: service_errors.PermissionDenied}
	}
	// check if the user is Owner of the Workout to update the exercise
	requestWorkout, err := u.workoutRepo.GetById(ctx, req.WorkoutId)
	if err != nil {
		return dto.WorkoutExerciseResponse{}, err
	}
	if requestWorkout.UserId != userId {
		return dto.WorkoutExerciseResponse{}, &service_errors.ServiceError{EndUserMessage: service_errors.PermissionDenied}
	}

	return u.base.Update(ctx, id, req)
}
func (u *WorkoutExerciseUsecase) Delete(ctx context.Context, id int) error {
	// Check if the user is Owner of the Workout
	userId := int(ctx.Value(constants.UserIdKey).(float64))
	workoutExercise, err := u.base.GetById(ctx, id)
	if err != nil {
		return err
	}
	workout, err := u.workoutRepo.GetById(ctx, workoutExercise.WorkoutId)
	if err != nil {
		return err
	}
	if userId != workout.UserId {
		return &service_errors.ServiceError{EndUserMessage: service_errors.PermissionDenied}
	}

	return u.base.Delete(ctx, id)
}
func (u *WorkoutExerciseUsecase) GetById(ctx context.Context, id int) (dto.WorkoutExerciseResponse, error) {
	// Check if the user is Owner of the Workout
	userId := int(ctx.Value(constants.UserIdKey).(float64))
	workoutExercise, err := u.base.GetById(ctx, id)
	if err != nil {
		return dto.WorkoutExerciseResponse{}, err
	}
	workout, err := u.workoutRepo.GetById(ctx, workoutExercise.WorkoutId)
	if err != nil {
		return dto.WorkoutExerciseResponse{}, err
	}
	if userId != workout.UserId {
		return dto.WorkoutExerciseResponse{}, &service_errors.ServiceError{EndUserMessage: service_errors.PermissionDenied}
	}

	return u.base.GetById(ctx, id)
}
