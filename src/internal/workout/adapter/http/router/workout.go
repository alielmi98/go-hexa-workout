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
}
