package category

import (
	"context"
	"forum/internal/model"
	"log"
)

type service struct {
	storage CategoryStorage
}
type CategoryService interface {
	Create(ctx context.Context, category *model.Category) error
	GetAll(ctx context.Context) ([]*model.Category, error)
	Delete(ctx context.Context, category *model.Category) error
	Update(ctx context.Context, category *model.Category) error
}

// NewService ...
func NewService(storage CategoryStorage) CategoryService {
	return &service{storage: storage}
}

func (s *service) Create(ctx context.Context, category *model.Category) error {

	return nil
}

func (s *service) GetAll(ctx context.Context) ([]*model.Category, error) {
	log.Println("start GetAll category service")
	categories, err := s.storage.GetAll(ctx)
	if err != nil {
		log.Printf("ERROR service category GetAll method: %v\n", err)
		return nil, err
	}
	log.Println("end GetAll category service")
	return categories, nil
}

func (s *service) Delete(ctx context.Context, category *model.Category) error {
	return nil
}

func (s *service) Update(ctx context.Context, category *model.Category) error {
	return nil
}
