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

func GetWorkoutExerciseRepository() workoutPort.WorkoutExerciseRepository {
	var preloads []db.PreloadEntity = []db.PreloadEntity{}
	workoutExerciseRepo := workoutInfraRepository.NewBaseRepository[workoutModels.WorkoutExercise](preloads)
	return workoutExerciseRepo
}

func GetScheduledWorkoutsRepository() workoutPort.ScheduledWorkoutsRepository {
	var preloads []db.PreloadEntity = []db.PreloadEntity{}
	scheduledWorkoutsRepo := workoutInfraRepository.NewBaseRepository[workoutModels.ScheduledWorkouts](preloads)
	return scheduledWorkoutsRepo
}

func GetWorkoutReportRepository() workoutPort.WorkoutReportRepository {
	var preloads []db.PreloadEntity = []db.PreloadEntity{}
	workoutReportRepo := workoutInfraRepository.NewBaseRepository[workoutModels.WorkoutReport](preloads)
	return workoutReportRepo
}
