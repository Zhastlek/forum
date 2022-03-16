package category

import (
	"context"
	"forum/internal/model"
)

type CategoryStorage interface {
	GetAll(ctx context.Context) ([]*model.Category, error)
	Create(ctx context.Context, category *model.Category) error
	Delete(ctx context.Context, category *model.Category) error
	Update(ctx context.Context, category *model.Category) error
}
