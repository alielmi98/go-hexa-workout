package test

import (
	"context"

	"github.com/alielmi98/go-hexa-workout/internal/user/adapter/http/dto"
	model "github.com/alielmi98/go-hexa-workout/internal/user/core/models"
	"github.com/alielmi98/go-hexa-workout/internal/user/core/usecase"
	"github.com/alielmi98/go-hexa-workout/internal/user/entity"
	"github.com/alielmi98/go-hexa-workout/pkg/config"
	"github.com/golang-jwt/jwt"
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
