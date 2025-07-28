package router

import (
	"github.com/alielmi98/go-hexa-workout/internal/user/adapter/http/handler"
	"github.com/alielmi98/go-hexa-workout/pkg/config"
	"github.com/gin-gonic/gin"
)

func Account(router *gin.RouterGroup, cfg *config.Config) {
	handler := handler.NewAccountHandler(cfg)
	router.POST("/register", handler.RegisterByUsername)
	router.POST("/login", handler.LoginByUsername)
	router.POST("/refresh-token", handler.RefreshToken)

}
