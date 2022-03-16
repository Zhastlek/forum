package user

import (
	"context"

	"forum/internal/adapters/api"
	"forum/internal/model"
	"log"
)

type service struct {
	storage UserStorage
}

// NewService ...
func NewService(storage UserStorage) api.UserService {
	return &service{storage: storage}
}

func (s *service) Create(ctx context.Context, dto *model.UserDto) error {
	user := &model.User{
		Name:     dto.Name,
		Email:    dto.Email,
		Password: dto.Password,
		UUID:     dto.UUID,
	}
	err := s.storage.Create(user)
	if err != nil {
		log.Println("Error service in user creation")
		return err
	}
	return nil
}

func (s *service) GetByUUID(ctx context.Context, uuid string) (*model.User, error) {
	user, err := s.storage.GetOne(uuid)
	if err != nil {
		log.Println("Error service in GetByUUID user", err)
		return nil, err
	}
	if user == nil {
		log.Println("Error service in GetByUUID user no one returned")
		return nil, nil
	}

	return user, nil
}

func (s *service) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	user, err := s.storage.GetOne(email)
	if err != nil {
		log.Println("Error service in GetByEmail user", err)
		return nil, err
	}
	if user == nil {
		log.Println("Error service in GetByEmail user no one returned")
	}
	return user, nil

}
