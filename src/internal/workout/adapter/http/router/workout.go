package router

import (
	"github.com/alielmi98/go-hexa-workout/internal/middlewares"
	"github.com/alielmi98/go-hexa-workout/internal/user/port"
	"github.com/alielmi98/go-hexa-workout/internal/workout/adapter/http/handler"
	"github.com/alielmi98/go-hexa-workout/pkg/config"
	"github.com/gin-gonic/gin"
)

func WorkoutRouters(r *gin.RouterGroup, cfg *config.Config, tokenProvider port.TokenProvider) {
	workoutHandler := handler.NewWorkoutHandler(cfg)
	// Workout
	r.POST("/workout/", middlewares.Authentication(cfg, tokenProvider), workoutHandler.Create)
	r.PUT("/workout/:id", middlewares.Authentication(cfg, tokenProvider), workoutHandler.Update)
	r.GET("/workout/:id", middlewares.Authentication(cfg, tokenProvider), workoutHandler.GetById)
	r.DELETE("/workout/:id", middlewares.Authentication(cfg, tokenProvider), workoutHandler.Delete)
	r.POST("/workout/get-by-filter", middlewares.Authentication(cfg, tokenProvider), workoutHandler.GetByFilter)

	// WorkoutExercise
	workoutExerciseHandler := handler.NewWorkoutExerciseHandler(cfg)
	r.POST("/workout-exercise/", middlewares.Authentication(cfg, tokenProvider), workoutExerciseHandler.Create)
	r.PUT("/workout-exercise/:id", middlewares.Authentication(cfg, tokenProvider), workoutExerciseHandler.Update)
	r.GET("/workout-exercise/:id", middlewares.Authentication(cfg, tokenProvider), workoutExerciseHandler.GetById)
	r.DELETE("/workout-exercise/:id", middlewares.Authentication(cfg, tokenProvider), workoutExerciseHandler.Delete)

	// ScheduledWorkout
	scheduledWorkoutHandler := handler.NewScheduledWorkoutsHandler(cfg)
	r.POST("/scheduled-workouts/", middlewares.Authentication(cfg, tokenProvider), scheduledWorkoutHandler.Create)
	r.PUT("/scheduled-workouts/:id", middlewares.Authentication(cfg, tokenProvider), scheduledWorkoutHandler.Update)
	r.GET("/scheduled-workouts/:id", middlewares.Authentication(cfg, tokenProvider), scheduledWorkoutHandler.GetById)
	r.DELETE("/scheduled-workouts/:id", middlewares.Authentication(cfg, tokenProvider), scheduledWorkoutHandler.Delete)

	// WorkoutReport
	workoutReportHandler := handler.NewWorkoutReportHandler(cfg)
	r.POST("/workout-report/", middlewares.Authentication(cfg, tokenProvider), workoutReportHandler.Create)
	r.PUT("/workout-report/:id", middlewares.Authentication(cfg, tokenProvider), workoutReportHandler.Update)
	r.GET("/workout-report/:id", middlewares.Authentication(cfg, tokenProvider), workoutReportHandler.GetById)
	r.DELETE("/workout-report/:id", middlewares.Authentication(cfg, tokenProvider), workoutReportHandler.Delete)
}
