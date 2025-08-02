package dependency

import (
	"github.com/alielmi98/go-hexa-workout/internal/user/adapter/auth"
	userInfraRepository "github.com/alielmi98/go-hexa-workout/internal/user/adapter/repo"
	userPort "github.com/alielmi98/go-hexa-workout/internal/user/port"
	workoutInfraRepository "github.com/alielmi98/go-hexa-workout/internal/workout/adapter/repo"
	workoutModels "github.com/alielmi98/go-hexa-workout/internal/workout/core/models"
	workoutPort "github.com/alielmi98/go-hexa-workout/internal/workout/port"
	"github.com/alielmi98/go-hexa-workout/pkg/config"
	"github.com/alielmi98/go-hexa-workout/pkg/db"
)

func GetUserRepository(cfg *config.Config) (userPort.UserRepository, userPort.TokenProvider) {
	return userInfraRepository.NewUserPgRepo(), auth.NewJwtProvider(cfg)
}

// Workout
func GetWorkoutRepository() workoutPort.WorkoutRepository {
	var preloads []db.PreloadEntity = []db.PreloadEntity{}
	return workoutInfraRepository.NewBaseRepository[workoutModels.Workout](preloads)
}
