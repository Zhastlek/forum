package session

import (
	"context"
	"forum/internal/model"
)

type SessionStorage interface {
	Create(ctx context.Context, session *model.Session) error
	Check(ctx context.Context, session *model.SessionDto) error
	Delete(ctx context.Context, session *model.SessionDto) error
}
