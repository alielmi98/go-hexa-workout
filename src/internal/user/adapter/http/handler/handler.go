package handler

import (
	"net/http"

	"github.com/alielmi98/go-hexa-workout/dependency"
	"github.com/alielmi98/go-hexa-workout/internal/user/adapter/http/dto"
	"github.com/alielmi98/go-hexa-workout/internal/user/core/usecase"
	"github.com/alielmi98/go-hexa-workout/pkg/config"
	"github.com/alielmi98/go-hexa-workout/pkg/helper"

	"github.com/gin-gonic/gin"
)

// AccountHandler ...
type AccountHandler struct {
	usecase *usecase.UserUsecase
	cfg     *config.Config
}

// NewAccountHandler ...
func NewAccountHandler(cfg *config.Config) *AccountHandler {
	return &AccountHandler{
		usecase: usecase.NewUserUsecase(cfg, dependency.GetUserRepository(cfg)),
		cfg:     cfg,
	}
}

// RegisterByUsername godoc
// @Summary RegisterByUsername
// @Description RegisterByUsername
// @Tags Account
// @Accept  json
// @Produce  json
// @Param Request body dto.RegisterUserByUsernameRequest true "RegisterUserByUsernameRequest"
// @Success 201 {object} helper.BaseHttpResponse "Success"
// @Failure 400 {object} helper.BaseHttpResponse "Failed"
// @Failure 409 {object} helper.BaseHttpResponse "Failed"
// @Router /v1/account/register [post]
func (h *AccountHandler) Create(c *gin.Context) {
	var req dto.RegisterUserByUsernameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			helper.GenerateBaseResponseWithValidationError(nil, false, helper.ValidationError, err))
		return
	}
	err := h.usecase.RegisterByUsername(c, &req)
	if err != nil {
		c.AbortWithStatusJSON(helper.TranslateErrorToStatusCode(err),
			helper.GenerateBaseResponseWithError(nil, false, helper.InternalError, err))
		return
	}
	c.JSON(http.StatusCreated, helper.GenerateBaseResponse("User created", true, helper.Success))
}
