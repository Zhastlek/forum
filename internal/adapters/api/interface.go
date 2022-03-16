package api

import (
	"context"
	"forum/internal/model"
	"net/http"
)

type Handler interface {
	Register(router *http.ServeMux)
}

type PostService interface {
	GetByID(ctx context.Context, id int) ([]*model.Post, error)
	GetAll(ctx context.Context) ([]*model.Post, error)
	Create(ctx context.Context, dto *model.CreatePostDto) (int, error)
	Update(ctx context.Context, post *model.UpdatePostDTO) error
	SortedByCategory(ctx context.Context, categoryId int) ([]*model.Post, error)
	DeletePost(ctx context.Context, id int) error

	CreateComment(ctx context.Context, dto *model.CreateCommentDto) error
	GetComment(ctx context.Context, postId int) ([]*model.Comment, error)
	DeleteComment(ctx context.Context, dto *model.Comment) error

	CheckPostLike(ctx context.Context, like *model.PostLikeAndDislike) error
	CheckCommentLike(ctx context.Context, like *model.CommentLikeAndDislike) error

	SortedByPost(ctx context.Context, userUUID string) ([]*model.Post, error)
	SortedByComment(ctx context.Context, userUUID string) ([]*model.Post, error)
	SortedByLike(ctx context.Context, userUUID string) ([]*model.Post, error)
}

type UserService interface {
	GetByUUID(ctx context.Context, uuid string) (*model.User, error)
	//GetAll(ctx context.Context, limit, offset int) []*model.User
	Create(ctx context.Context, dto *model.UserDto) error
	GetByEmail(ctx context.Context, email string) (*model.User, error)
}
