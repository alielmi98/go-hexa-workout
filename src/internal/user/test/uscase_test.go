package test

import (
	"context"
	"errors"
	"testing"

	"github.com/alecthomas/assert/v2"
	"github.com/alielmi98/go-hexa-workout/internal/user/adapter/http/dto"
	model "github.com/alielmi98/go-hexa-workout/internal/user/core/models"
	"github.com/alielmi98/go-hexa-workout/internal/user/core/usecase"
	"github.com/alielmi98/go-hexa-workout/pkg/config"
	"golang.org/x/crypto/bcrypt"
)

func TestRegisterUser_Success(t *testing.T) {
	repo := &MockUserRepository{}
	useCase, mockRepo := setup(repo)

	user := &dto.RegisterUserByUsernameRequest{
		Username:  "testuser",
		Password:  "password",
		FirstName: "ali",
		LastName:  "elmi",
		Email:     "ali.elmi@example.com",
	}

	err := useCase.RegisterByUsername(nil, user)
	assert.NoError(t, err)
	assert.True(t, mockRepo.SaveCalled)
	assert.Equal(t, user, mockRepo.SaveUser)
}

func TestRegisterUser_Failure(t *testing.T) {
	repo := &MockUserRepository{
		CreateFn: func(ctx context.Context, user *model.User) error {
			return errors.New("db error")
		},
	}
	useCase, _ := setup(repo)

	user := &dto.RegisterUserByUsernameRequest{
		Username:  "testuser",
		Password:  "password",
		FirstName: "ali",
		LastName:  "elmi",
		Email:     "ali.elmi@example.com",
	}

	err := useCase.RegisterByUsername(nil, user)
	assert.Error(t, err)
}

func TestRegisterUser_UsernameExists(t *testing.T) {
	repo := &MockUserRepository{
		ExistsByUsernameFn: func(username string) (bool, error) {
			if username == "existinguser" {
				return true, nil
			}
			return false, nil
		},
	}
	useCase, _ := setup(repo)

	user := &dto.RegisterUserByUsernameRequest{
		Username:  "existinguser",
		Password:  "password",
		FirstName: "ali",
		LastName:  "elmi",
		Email:     "ali.elmi@example.com",
	}

	err := useCase.RegisterByUsername(nil, user)
	assert.Error(t, err)
	assert.Equal(t, "Username exists", err.Error())
}
func TestRegisterUser_EmailExists(t *testing.T) {
	repo := &MockUserRepository{
		ExistsByEmailFn: func(email string) (bool, error) {
			if email == "existingemail@example.com" {
				return true, nil
			}
			return false, nil
		},
	}
	useCase, _ := setup(repo)

	user := &dto.RegisterUserByUsernameRequest{
		Username:  "newuser",
		Password:  "password",
		FirstName: "ali",
		LastName:  "elmi",
		Email:     "existingemail@example.com",
	}

	err := useCase.RegisterByUsername(nil, user)
	assert.Error(t, err)
	assert.Equal(t, "Email exists", err.Error())
}

func TestLoginUser_Success(t *testing.T) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	repo := &MockUserRepository{
		FindByUsernameFn: func(ctx context.Context, username string) (*model.User, error) {
			return &model.User{
				Username:  "testuser",
				Password:  string(hashedPassword),
				Email:     "testuser@example.com",
				FirstName: "Test",
				LastName:  "User",
			}, nil // hashed password
		},
	}
	useCase, _ := setup(repo)

	req := &dto.LoginByUsernameRequest{
		Username: "testuser",
		Password: "password",
	}

	tokenDetail, err := useCase.LoginByUsername(nil, req)
	assert.NoError(t, err)
	assert.True(t, tokenDetail != nil)
	assert.True(t, tokenDetail.AccessToken != "")
	assert.True(t, tokenDetail.RefreshToken != "")
}

func TestLoginUser_Failure(t *testing.T) {
	repo := &MockUserRepository{
		FindByUsernameFn: func(ctx context.Context, username string) (*model.User, error) {
			return nil, errors.New("user not found")
		},
	}
	useCase, _ := setup(repo)

	req := &dto.LoginByUsernameRequest{
		Username: "nonexistentuser",
		Password: "password",
	}

	tokenDetail, err := useCase.LoginByUsername(nil, req)
	assert.Error(t, err)
	assert.True(t, tokenDetail == nil)
}

func TestLoginUser_InvalidPassword(t *testing.T) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correctpassword"), bcrypt.DefaultCost)
	repo := &MockUserRepository{
		FindByUsernameFn: func(ctx context.Context, username string) (*model.User, error) {
			return &model.User{
				Username:  "testuser",
				Password:  string(hashedPassword),
				Email:     "testuser@example.com",
				FirstName: "Test",
				LastName:  "User",
			}, nil
		},
	}
	useCase, _ := setup(repo)

	req := &dto.LoginByUsernameRequest{
		Username: "testuser",
		Password: "wrongpassword",
	}

	tokenDetail, err := useCase.LoginByUsername(nil, req)
	assert.Error(t, err)
	assert.True(t, tokenDetail == nil)
}
func TestLoginUser_GenerateTokenError(t *testing.T) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	repo := &MockUserRepository{
		FindByUsernameFn: func(ctx context.Context, username string) (*model.User, error) {
			return &model.User{
				Username:  "testuser",
				Password:  string(hashedPassword),
				Email:     "testuser@example.com",
				FirstName: "Test",
				LastName:  "User",
			}, nil
		},
	}
	useCase, _ := setup(repo)

	req := &dto.LoginByUsernameRequest{
		Username: "testuser",
		Password: "wrongpassword",
	}

	tokenDetail, err := useCase.LoginByUsername(nil, req)
	assert.Error(t, err)
	assert.True(t, tokenDetail == nil)
}

func TestRefreshToken_Success(t *testing.T) {
	repo := &MockUserRepository{}
	useCase, _ := setup(repo)

	tokenDetail, err := useCase.RefreshToken("valid-refresh-token")
	assert.NoError(t, err)
	assert.True(t, tokenDetail != nil)
	assert.True(t, tokenDetail.AccessToken != "")
	assert.True(t, tokenDetail.RefreshToken != "")
	assert.Equal(t, int64(0), tokenDetail.AccessTokenExpireTime)
	assert.Equal(t, int64(0), tokenDetail.RefreshTokenExpireTime)
}

func TestRefreshToken_Error(t *testing.T) {
	mockConfig := &config.Config{}
	mockRepo := &MockUserRepository{}
	mockToken := &MockTokenProvider{
		RefreshTokenFn: func(refreshToken string) (*dto.TokenDetail, error) {
			return nil, errors.New("refresh token error")
		},
	}
	useCase := usecase.NewUserUsecase(mockConfig, mockRepo, mockToken)

	tokenDetail, err := useCase.RefreshToken("invalid-refresh-token")
	assert.Error(t, err)
	assert.True(t, tokenDetail == nil)
}
