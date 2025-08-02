package repository

import (
	"context"

	"github.com/alielmi98/go-hexa-workout/internal/workout/port/filter"
	"gorm.io/gorm"
)

type BaseRepository[TEntity any] interface {
	Create(ctx context.Context, entity TEntity) (TEntity, error)
	Update(ctx context.Context, id int, entity TEntity) (TEntity, error)
	Delete(ctx context.Context, id int) error
	GetById(ctx context.Context, id int) (TEntity, error)
	GetByFilter(ctx context.Context, req filter.PaginationInputWithFilter) (int64, *[]TEntity, error)
	BeginTransaction(ctx context.Context) (*gorm.DB, error)
	CreateTx(tx *gorm.DB, entity TEntity) (TEntity, error)
}
