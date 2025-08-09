package usecase

import (
	"context"

	"github.com/alielmi98/go-hexa-workout/internal/workout/core/models"
	"github.com/alielmi98/go-hexa-workout/internal/workout/core/usecase/dto"
	"github.com/alielmi98/go-hexa-workout/internal/workout/port"
	"github.com/alielmi98/go-hexa-workout/pkg/config"
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
	err := u.base.CheckOwnership(ctx, u.workoutRepo, req.WorkoutId)
	if err != nil {
		return dto.WorkoutExerciseResponse{}, err
	}

	return u.base.Create(ctx, req)
}
func (u *WorkoutExerciseUsecase) Update(ctx context.Context, id int, req dto.UpdateWorkoutExerciseRequest) (dto.WorkoutExerciseResponse, error) {
	// Check if the user is Owner of the Workout
	workoutExercise, err := u.base.GetById(ctx, id)
	if err != nil {
		return dto.WorkoutExerciseResponse{}, err
	}
	err = u.base.CheckOwnership(ctx, u.workoutRepo, workoutExercise.WorkoutId)
	if err != nil {
		return dto.WorkoutExerciseResponse{}, err
	}

	// check if the user is Owner of the Workout to update the exercise
	err = u.base.CheckOwnership(ctx, u.workoutRepo, req.WorkoutId)
	if err != nil {
		return dto.WorkoutExerciseResponse{}, err
	}

	return u.base.Update(ctx, id, req)
}
func (u *WorkoutExerciseUsecase) Delete(ctx context.Context, id int) error {
	// Check if the user is Owner of the Workout
	workoutExercise, err := u.base.GetById(ctx, id)
	if err != nil {
		return err
	}
	err = u.base.CheckOwnership(ctx, u.workoutRepo, workoutExercise.WorkoutId)
	if err != nil {
		return err
	}

	return u.base.Delete(ctx, id)
}
func (u *WorkoutExerciseUsecase) GetById(ctx context.Context, id int) (dto.WorkoutExerciseResponse, error) {
	// Check if the user is Owner of the Workout
	workoutExercise, err := u.base.GetById(ctx, id)
	if err != nil {
		return dto.WorkoutExerciseResponse{}, err
	}
	err = u.base.CheckOwnership(ctx, u.workoutRepo, workoutExercise.WorkoutId)
	if err != nil {
		return dto.WorkoutExerciseResponse{}, err
	}

	return workoutExercise, nil
}
