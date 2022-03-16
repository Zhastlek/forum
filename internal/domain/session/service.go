package session

import (
	"context"
	"forum/internal/model"
	"log"
)

type service struct {
	storage SessionStorage
}

type SessionService interface {
	Create(ctx context.Context, session *model.Session) error
	Check(ctx context.Context, session *model.SessionDto) error
	Delete(ctx context.Context, session *model.SessionDto) error
}

func NewService(storage SessionStorage) SessionService {
	return &service{storage: storage}
}

func (s *service) Create(ctx context.Context, session *model.Session) error {
	err := s.storage.Create(ctx, session)
	if err != nil {
		log.Println("Error session service Create function")
		return err
	}
	return nil
}

func (s *service) Check(ctx context.Context, session *model.SessionDto) error {

	err := s.storage.Check(ctx, session)
	if err != nil {
		log.Println("Error session service check function")
		return err
	}
	return nil
}

func (s *service) Delete(ctx context.Context, session *model.SessionDto) error {
	err := s.storage.Delete(ctx, session)
	if err != nil {
		log.Printf("ERROR session service DELETE method :--> %v\n", err)
		return err
	}
	return nil
}
