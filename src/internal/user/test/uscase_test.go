package test

import (
	"context"
	"errors"
	"testing"

	"github.com/alecthomas/assert/v2"
	"github.com/alielmi98/go-hexa-workout/internal/user/adapter/http/dto"
	model "github.com/alielmi98/go-hexa-workout/internal/user/core/models"
	"github.com/alielmi98/go-hexa-workout/internal/user/core/usecase"
	"github.com/alielmi98/go-hexa-workout/internal/user/entity"
	"github.com/alielmi98/go-hexa-workout/pkg/config"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type MockUserRepository struct {
	SaveCalled         bool
	SaveUser           *dto.RegisterUserByUsernameRequest
	FindByUsernameFn   func(ctx context.Context, username string) (*model.User, error)
	CreateFn           func(ctx context.Context, user *model.User) error
	ExistsByUsernameFn func(username string) (bool, error)
	ExistsByEmailFn    func(email string) (bool, error)
}

func (m *MockUserRepository) Create(ctx context.Context, user *model.User) error {
	if m.CreateFn != nil {
		return m.CreateFn(ctx, user)
	}
	m.SaveCalled = true
	m.SaveUser = &dto.RegisterUserByUsernameRequest{
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Password:  user.Password,
	}
	return nil
}
func (m *MockUserRepository) GetByID(ctx context.Context, id int) (*model.User, error) {
	return &model.User{Id: id, Username: "testuser", Password: "password"}, nil
}
func (m *MockUserRepository) Update(ctx context.Context, id int, user *model.User) error {
	return nil
}
func (m *MockUserRepository) Delete(ctx context.Context, id int) error {
	return nil
}
func (m *MockUserRepository) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	if m.FindByUsernameFn != nil {
		return m.FindByUsernameFn(ctx, username)
	}
	return nil, nil
}
func (m *MockUserRepository) Save(user *dto.RegisterUserByUsernameRequest) error {
	m.SaveCalled = true
	m.SaveUser = user
	return nil
}
func (m *MockUserRepository) ExistsByEmail(email string) (bool, error) {
	if m.ExistsByEmailFn != nil {
		return m.ExistsByEmailFn(email)
	}
	return false, nil
}
func (m *MockUserRepository) ExistsByUsername(username string) (bool, error) {
	if m.ExistsByUsernameFn != nil {
		return m.ExistsByUsernameFn(username)
	}
	return false, nil
}

type MockTokenProvider struct {
	RefreshTokenFn func(refreshToken string) (*dto.TokenDetail, error)
}

func (m *MockTokenProvider) GenerateToken(token *entity.TokenPayload) (*dto.TokenDetail, error) {
	return &dto.TokenDetail{AccessToken: "token", RefreshToken: "refresh", AccessTokenExpireTime: 0, RefreshTokenExpireTime: 0}, nil
}
func (m *MockTokenProvider) VerifyToken(token string) (*jwt.Token, error) { return &jwt.Token{}, nil }
func (m *MockTokenProvider) GetClaims(token string) (map[string]interface{}, error) {
	return map[string]interface{}{}, nil
}
func (m *MockTokenProvider) RefreshToken(refreshToken string) (*dto.TokenDetail, error) {
	if m.RefreshTokenFn != nil {
		return m.RefreshTokenFn(refreshToken)
	}
	return &dto.TokenDetail{AccessToken: "new-token", RefreshToken: "new-refresh", AccessTokenExpireTime: 0, RefreshTokenExpireTime: 0}, nil
}

func setup(repo *MockUserRepository) (*usecase.UserUsecase, *MockUserRepository) {
	mockToken := &MockTokenProvider{}
	mockConfig := &config.Config{}
	useCase := usecase.NewUserUsecase(mockConfig, repo, mockToken)
	return useCase, repo
}

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
