package usecase

import (
	"context"

	"github.com/alielmi98/go-hexa-workout/internal/workout/adapter/http/dto"
	"github.com/alielmi98/go-hexa-workout/internal/workout/core/models"
	"github.com/alielmi98/go-hexa-workout/internal/workout/port"
	"github.com/alielmi98/go-hexa-workout/internal/workout/port/filter"
	"github.com/alielmi98/go-hexa-workout/pkg/config"
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
	return u.base.Create(ctx, req)
}

func (u *WorkoutUsecase) Update(ctx context.Context, id int, req dto.UpdateWorkoutRequest) (dto.WorkoutResponse, error) {
	return u.base.Update(ctx, id, req)
}

func (u *WorkoutUsecase) Delete(ctx context.Context, id int) error {
	return u.base.Delete(ctx, id)
}
func (u *WorkoutUsecase) GetById(ctx context.Context, id int) (dto.WorkoutResponse, error) {
	return u.base.GetById(ctx, id)
}
func (u *WorkoutUsecase) GetByFilter(ctx context.Context, req filter.PaginationInputWithFilter) (*filter.PagedList[dto.WorkoutResponse], error) {
	return u.base.GetByFilter(ctx, req)
}
