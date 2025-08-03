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

// midedlewares
func GetTokenProvider(cfg *config.Config) userPort.TokenProvider {
	return auth.NewJwtProvider(cfg)
}

// user
func GetUserRepository(cfg *config.Config) (userPort.UserRepository, userPort.TokenProvider) {
	return userInfraRepository.NewUserPgRepo(), auth.NewJwtProvider(cfg)
}

// Workout
func GetWorkoutRepository() workoutPort.WorkoutRepository {
	var preloads []db.PreloadEntity = []db.PreloadEntity{}
	return workoutInfraRepository.NewBaseRepository[workoutModels.Workout](preloads)
}

func GetWorkoutExerciseRepository() (workoutPort.WorkoutExerciseRepository, workoutPort.WorkoutRepository) {
	var workoutExercisePreloads []db.PreloadEntity = []db.PreloadEntity{}
	workoutExerciseRepo := workoutInfraRepository.NewBaseRepository[workoutModels.WorkoutExercise](workoutExercisePreloads)
	var workoutPreloads []db.PreloadEntity = []db.PreloadEntity{}
	workoutRepo := workoutInfraRepository.NewBaseRepository[workoutModels.Workout](workoutPreloads)
	return workoutExerciseRepo, workoutRepo
}
