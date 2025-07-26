package usecase

import (
	"context"
	"log"

	"github.com/alielmi98/go-hexa-workout/constants"
	"github.com/alielmi98/go-hexa-workout/internal/user/adapter/http/dto"
	model "github.com/alielmi98/go-hexa-workout/internal/user/core/models"
	"github.com/alielmi98/go-hexa-workout/internal/user/port"
	"github.com/alielmi98/go-hexa-workout/pkg/config"
	"golang.org/x/crypto/bcrypt"
)

type UserUsecase struct {
	cfg  *config.Config
	repo port.UserRepository
}

func NewUserUsecase(cfg *config.Config, repository port.UserRepository) *UserUsecase {
	return &UserUsecase{
		cfg:  cfg,
		repo: repository,
	}
}

// Register by username
func (s *UserUsecase) RegisterByUsername(ctx context.Context, req *dto.RegisterUserByUsernameRequest) error {
	u := &model.User{
		Username:  req.Username,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
	}
	// Hash password
	bp := []byte(req.Password)
	hp, err := bcrypt.GenerateFromPassword(bp, bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Caller:%s Level:%s Msg:%s", constants.General, constants.HashPassword, err.Error())
		return err
	}
	req.Password = string(hp)
	u.Password = req.Password

	err = s.repo.Create(ctx, u)
	if err != nil {
		return err
	}
	return nil

}
