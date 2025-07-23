package port

import (
	"context"

	"github.com/alielmi98/go-hexa-workout/internal/user/core"
)

type UserRepository interface {
	Create(ctx context.Context, user *core.User) error
	GetByID(ctx context.Context, id int) (*core.User, error)
	Update(ctx context.Context, id int, user *core.User) error
	Delete(ctx context.Context, id int) error
	FindByUsername(ctx context.Context, username string) (*core.User, error)
}
