package router

import (
	"github.com/alielmi98/go-hexa-workout/internal/workout/adapter/http/handler"
	"github.com/alielmi98/go-hexa-workout/pkg/config"
	"github.com/gin-gonic/gin"
)

func WorkoutRouters(r *gin.RouterGroup, cfg *config.Config) {
	handler := handler.NewWorkoutHandler(cfg)
	// Create a Workout
	r.POST("/workout/", handler.Create)
	r.PUT("/workout/:id", handler.Update)
	r.GET("/workout/:id", handler.GetById)
	r.DELETE("/workout/:id", handler.Delete)
	r.POST("/workout/get-by-filter", handler.GetByFilter)
}
