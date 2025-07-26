package dependency

import (
	"github.com/alielmi98/go-hexa-workout/internal/user/adapter/auth"
	infraRepository "github.com/alielmi98/go-hexa-workout/internal/user/adapter/repo"
	contractRepository "github.com/alielmi98/go-hexa-workout/internal/user/port"
	userPort "github.com/alielmi98/go-hexa-workout/internal/user/port"
	"github.com/alielmi98/go-hexa-workout/pkg/config"
)

func GetUserRepository(cfg *config.Config) (contractRepository.UserRepository, userPort.TokenProvider) {
	return infraRepository.NewUserPgRepo(), auth.NewJwtProvider(cfg)
}
