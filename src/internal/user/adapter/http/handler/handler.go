package handler

import (
	"net/http"

	"github.com/alielmi98/go-hexa-workout/constants"
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
	repo, token := dependency.GetUserRepository(cfg)
	return &AccountHandler{
		usecase: usecase.NewUserUsecase(cfg, repo, token),
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

// LoginByUsername godoc
// @Summary LoginByUsername
// @Description LoginByUsername
// @Tags Account
// @Accept  json
// @Produce  json
// @Param Request body dto.LoginByUsernameRequest true "LoginByUsernameRequest"
// @Success 200 {object} helper.BaseHttpResponse "Success"
// @Failure 400 {object} helper.BaseHttpResponse "Failed"
// @Failure 401 {object} helper.BaseHttpResponse "Failed"
// @Router /v1/account/login [post]
func (h *AccountHandler) LoginByUsername(c *gin.Context) {
	var req dto.LoginByUsernameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			helper.GenerateBaseResponseWithValidationError(nil, false, helper.ValidationError, err))
		return
	}
	td, err := h.usecase.LoginByUsername(c, &req)
	if err != nil {
		c.AbortWithStatusJSON(helper.TranslateErrorToStatusCode(err),
			helper.GenerateBaseResponseWithError(nil, false, helper.InternalError, err))
		return
	}

	// Set the refresh token in a cookie
	// Set the new refresh token in a cookie
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     constants.RefreshTokenCookieName,
		Value:    td.RefreshToken,
		MaxAge:   int(h.cfg.JWT.RefreshTokenExpireDuration * 60),
		Path:     "/",
		Domain:   h.cfg.Server.Domain,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
	c.JSON(http.StatusOK, helper.GenerateBaseResponse(td, true, helper.Success))
}

// RefreshToken godoc
// @Summary RefreshToken
// @Description RefreshToken
// @Tags Account
// @Accept  json
// @Produce  json
// @Success 200 {object} helper.BaseHttpResponse "Success"
// @Failure 400 {object} helper.BaseHttpResponse "Failed"
// @Failure 401 {object} helper.BaseHttpResponse "Failed"
// @Router /v1/account/refresh-token [post]
func (h *AccountHandler) RefreshToken(c *gin.Context) {
	// Get the refresh token from the cookie
	refreshToken, err := c.Cookie(constants.RefreshTokenCookieName)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			helper.GenerateBaseResponseWithError(nil, false, helper.AuthError, err))
		return
	}
	// Call the usecase to refresh the token
	td, err := h.usecase.RefreshToken(refreshToken)
	if err != nil {
		c.AbortWithStatusJSON(helper.TranslateErrorToStatusCode(err),
			helper.GenerateBaseResponseWithError(nil, false, helper.InternalError, err))
		return
	}
	// Set the new refresh token in a cookie
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     constants.RefreshTokenCookieName,
		Value:    td.RefreshToken,
		MaxAge:   int(h.cfg.JWT.RefreshTokenExpireDuration * 60),
		Path:     "/",
		Domain:   h.cfg.Server.Domain,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
	// Return the token details
	c.JSON(http.StatusOK, helper.GenerateBaseResponse(td, true, helper.Success))
}
