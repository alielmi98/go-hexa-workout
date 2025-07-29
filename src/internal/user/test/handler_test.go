package test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alecthomas/assert/v2"
	"github.com/alielmi98/go-hexa-workout/constants"
	"github.com/alielmi98/go-hexa-workout/internal/user/adapter/http/dto"
	"github.com/alielmi98/go-hexa-workout/internal/user/adapter/http/handler"
	model "github.com/alielmi98/go-hexa-workout/internal/user/core/models"
	"github.com/alielmi98/go-hexa-workout/internal/user/core/usecase"
	"github.com/alielmi98/go-hexa-workout/pkg/config"
	"github.com/alielmi98/go-hexa-workout/pkg/helper"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func TestLoginByUsername_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	mockRepo := &MockUserRepository{
		FindByUsernameFn: func(ctx context.Context, username string) (*model.User, error) {
			return &model.User{
				Id:       1,
				Username: "testuser",
				Password: string(hashedPassword),
			}, nil
		},
	}
	mockToken := &MockTokenProvider{}

	cfg := &config.Config{
		JWT: config.JWTConfig{
			RefreshTokenExpireDuration: 60,
		},
		Server: config.ServerConfig{
			Domain: "localhost",
		},
	}
	usecase := usecase.NewUserUsecase(cfg, mockRepo, mockToken)
	accountHandler := &handler.AccountHandler{
		Usecase: usecase,
		Cfg:     cfg,
	}

	router := gin.Default()
	router.POST("/v1/account/login", accountHandler.LoginByUsername)

	loginRequest := dto.LoginByUsernameRequest{
		Username: "testuser",
		Password: "password",
	}

	jsonData, _ := json.Marshal(loginRequest)
	req, _ := http.NewRequest("POST", "/v1/account/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response helper.BaseHttpResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, true, response.Success)
	assert.True(t, response.Result != nil)

	cookie := w.Header().Get("Set-Cookie")
	assert.Contains(t, cookie, constants.RefreshTokenCookieName)
}

func TestLoginByUsername_InvalidCredentials(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := &MockUserRepository{
		FindByUsernameFn: func(ctx context.Context, username string) (*model.User, error) {
			return &model.User{
				Username: "testuser",
				Password: "hashedpassword",
			}, nil
		},
	}
	mockToken := &MockTokenProvider{}

	cfg := &config.Config{}
	useCase := usecase.NewUserUsecase(cfg, repo, mockToken)

	accountHandler := &handler.AccountHandler{
		Usecase: useCase,
		Cfg:     cfg,
	}

	router := gin.Default()
	router.POST("/v1/account/login", accountHandler.LoginByUsername)

	reqBody := dto.LoginByUsernameRequest{
		Username: "testuser",
		Password: "wrongpassword",
	}
	jsonData, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/v1/account/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response helper.BaseHttpResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, false, response.Success)
	assert.Equal(t, "username or password invalid", response.Error)
}

func TestRegisterByUsername_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockRepo := &MockUserRepository{}
	mockToken := &MockTokenProvider{}

	cfg := &config.Config{}
	usecase := usecase.NewUserUsecase(cfg, mockRepo, mockToken)
	accountHandler := &handler.AccountHandler{
		Usecase: usecase,
		Cfg:     cfg,
	}

	router := gin.Default()
	router.POST("/v1/account/register", accountHandler.RegisterByUsername)

	registerRequest := dto.RegisterUserByUsernameRequest{
		Username:  "newuser",
		Password:  "newpassword",
		Email:     "newuser@example.com",
		FirstName: "New",
		LastName:  "User",
	}

	jsonData, _ := json.Marshal(registerRequest)
	req, _ := http.NewRequest("POST", "/v1/account/register", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response helper.BaseHttpResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, true, response.Success)
	assert.Equal(t, "User created", response.Result)
}

func TestRegisterByUsername_ValidationError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockRepo := &MockUserRepository{}
	mockToken := &MockTokenProvider{}

	cfg := &config.Config{}
	usecase := usecase.NewUserUsecase(cfg, mockRepo, mockToken)
	accountHandler := &handler.AccountHandler{
		Usecase: usecase,
		Cfg:     cfg,
	}

	router := gin.Default()
	router.POST("/v1/account/register", accountHandler.RegisterByUsername)

	registerRequest := dto.RegisterUserByUsernameRequest{
		Username:  "", // Invalid username
		Password:  "newpassword",
		Email:     "newuser@example.com",
		FirstName: "New",
		LastName:  "User",
	}

	jsonData, _ := json.Marshal(registerRequest)
	req, _ := http.NewRequest("POST", "/v1/account/register", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var response helper.BaseHttpResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, false, response.Success)
	assert.Equal(t, helper.ValidationError, response.ResultCode)

}

func TestRegisterByUsername_UsernameExists(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockRepo := &MockUserRepository{
		ExistsByUsernameFn: func(username string) (bool, error) {
			return true, nil // Simulate existing username
		},
	}
	mockToken := &MockTokenProvider{}

	cfg := &config.Config{}
	usecase := usecase.NewUserUsecase(cfg, mockRepo, mockToken)
	accountHandler := &handler.AccountHandler{
		Usecase: usecase,
		Cfg:     cfg,
	}

	router := gin.Default()
	router.POST("/v1/account/register", accountHandler.RegisterByUsername)

	registerRequest := dto.RegisterUserByUsernameRequest{
		Username:  "existinguser",
		Password:  "newpassword",
		Email:     "existinguser@example.com",
		FirstName: "Existing",
		LastName:  "User",
	}

	jsonData, _ := json.Marshal(registerRequest)
	req, _ := http.NewRequest("POST", "/v1/account/register", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)

	var response helper.BaseHttpResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, false, response.Success)
	assert.Equal(t, "Username exists", response.Error)

}

func TestRegisterByUsername_EmailExists(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockRepo := &MockUserRepository{
		ExistsByEmailFn: func(email string) (bool, error) {
			return true, nil // Simulate existing email
		},
	}
	mockToken := &MockTokenProvider{}

	cfg := &config.Config{}
	usecase := usecase.NewUserUsecase(cfg, mockRepo, mockToken)
	accountHandler := &handler.AccountHandler{
		Usecase: usecase,
		Cfg:     cfg,
	}

	router := gin.Default()
	router.POST("/v1/account/register", accountHandler.RegisterByUsername)

	registerRequest := dto.RegisterUserByUsernameRequest{
		Username:  "newuser",
		Password:  "newpassword",
		Email:     "newuser@example.com",
		FirstName: "New",
		LastName:  "User",
	}

	jsonData, _ := json.Marshal(registerRequest)
	req, _ := http.NewRequest("POST", "/v1/account/register", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)

	var response helper.BaseHttpResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, false, response.Success)
	assert.Equal(t, "Email exists", response.Error)

}

func TestRefreshTokenHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockRepo := &MockUserRepository{}
	mockToken := &MockTokenProvider{
		RefreshTokenFn: func(token string) (*dto.TokenDetail, error) {
			return &dto.TokenDetail{
				AccessToken:  "newAccessToken",
				RefreshToken: "newRefreshToken",
			}, nil
		},
	}

	cfg := &config.Config{
		JWT: config.JWTConfig{
			RefreshTokenExpireDuration: 60,
		},
		Server: config.ServerConfig{
			Domain: "localhost",
		},
	}
	usecase := usecase.NewUserUsecase(cfg, mockRepo, mockToken)
	accountHandler := &handler.AccountHandler{
		Usecase: usecase,
		Cfg:     cfg,
	}

	router := gin.Default()
	router.POST("/v1/account/refresh-token", accountHandler.RefreshToken)

	req, _ := http.NewRequest("POST", "/v1/account/refresh-token", nil)
	req.Header.Set("Authorization", "Bearer validRefreshToken")
	req.AddCookie(&http.Cookie{
		Name:  constants.RefreshTokenCookieName,
		Value: "refreshTokenValue",
	})

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response helper.BaseHttpResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, true, response.Success)
	// Check if the response contains the new access token
	result, _ := response.Result.(map[string]interface{})
	assert.Equal(t, "newAccessToken", result["accessToken"])
	// Get the new refresh token from the cookie
	cookies := w.Result().Cookies()
	var newRefreshToken string
	for _, cookie := range cookies {
		if cookie.Name == constants.RefreshTokenCookieName {
			newRefreshToken = cookie.Value
			break
		}
	}
	assert.Equal(t, "newRefreshToken", newRefreshToken)

}

func TestRefreshTokenHandler_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockRepo := &MockUserRepository{}
	mockToken := &MockTokenProvider{
		RefreshTokenFn: func(token string) (*dto.TokenDetail, error) {
			return nil, errors.New("refresh token error")
		},
	}

	cfg := &config.Config{}
	usecase := usecase.NewUserUsecase(cfg, mockRepo, mockToken)
	accountHandler := &handler.AccountHandler{
		Usecase: usecase,
		Cfg:     cfg,
	}

	router := gin.Default()
	router.POST("/v1/account/refresh-token", accountHandler.RefreshToken)

	req, _ := http.NewRequest("POST", "/v1/account/refresh-token", nil)
	req.AddCookie(&http.Cookie{
		Name:  constants.RefreshTokenCookieName,
		Value: "invalidRefreshToken",
	})

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response helper.BaseHttpResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, false, response.Success)
	assert.Equal(t, "refresh token error", response.Error)
}

func TestRefreshTokenHandler_BadRequest(t *testing.T) {
	// This test simulates a scenario where the refresh token is not provided or is invalid.
	// It checks if the handler returns a Bad Request status code and the appropriate error message.
	gin.SetMode(gin.TestMode)
	mockRepo := &MockUserRepository{}
	mockToken := &MockTokenProvider{}

	cfg := &config.Config{
		JWT: config.JWTConfig{
			RefreshTokenExpireDuration: 60,
		},
		Server: config.ServerConfig{
			Domain: "localhost",
		},
	}
	usecase := usecase.NewUserUsecase(cfg, mockRepo, mockToken)
	accountHandler := &handler.AccountHandler{
		Usecase: usecase,
		Cfg:     cfg,
	}

	router := gin.Default()
	router.POST("/v1/account/refresh-token", accountHandler.RefreshToken)

	req, _ := http.NewRequest("POST", "/v1/account/refresh-token", nil)
	req.Header.Set("Authorization", "Bearer validRefreshToken")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response helper.BaseHttpResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, false, response.Success)
	assert.Equal(t, "http: named cookie not present", response.Error)
}
