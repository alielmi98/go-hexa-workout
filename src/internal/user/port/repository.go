package port

import (
	"context"

	model "github.com/alielmi98/go-hexa-workout/internal/user/core/models"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	GetByID(ctx context.Context, id int) (*model.User, error)
	Update(ctx context.Context, id int, user *model.User) error
	Delete(ctx context.Context, id int) error
	FindByUsername(ctx context.Context, username string) (*model.User, error)
}
