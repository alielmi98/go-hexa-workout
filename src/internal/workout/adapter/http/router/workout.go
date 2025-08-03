package router

import (
	"github.com/alielmi98/go-hexa-workout/dependency"
	"github.com/alielmi98/go-hexa-workout/internal/middlewares"
	"github.com/alielmi98/go-hexa-workout/internal/workout/adapter/http/handler"
	"github.com/alielmi98/go-hexa-workout/pkg/config"
	"github.com/gin-gonic/gin"
)

func WorkoutRouters(r *gin.RouterGroup, cfg *config.Config) {
	handler := handler.NewWorkoutHandler(cfg)
	tokenProvider := dependency.GetTokenProvider(cfg)
	// Create a Workout
	r.POST("/workout/", middlewares.Authentication(cfg, tokenProvider), handler.Create)
	r.PUT("/workout/:id", middlewares.Authentication(cfg, tokenProvider), handler.Update)
	r.GET("/workout/:id", middlewares.Authentication(cfg, tokenProvider), handler.GetById)
	r.DELETE("/workout/:id", middlewares.Authentication(cfg, tokenProvider), handler.Delete)
	r.POST("/workout/get-by-filter", middlewares.Authentication(cfg, tokenProvider), handler.GetByFilter)
}
