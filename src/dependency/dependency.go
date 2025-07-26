package dependency

import (
	infraRepository "github.com/alielmi98/go-hexa-workout/internal/user/adapter/repo"
	contractRepository "github.com/alielmi98/go-hexa-workout/internal/user/port"
	"github.com/alielmi98/go-hexa-workout/pkg/config"
)

func GetUserRepository(cfg *config.Config) contractRepository.UserRepository {
	return infraRepository.NewUserPgRepo()
}
