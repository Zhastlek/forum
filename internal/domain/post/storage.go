package post

import (
	"forum/internal/model"
)

// Storage ...
type PostStorage interface {
	GetOne(id int) ([]*model.Post, error)
	GetAll() ([]*model.Post, error)
	Create(post *model.Post) (int, error)
	CreatePostCategory(postId int, categoryId int) error
	DeletePostCategory(postId int) error
	Delete(id int) error
	Update(post *model.Post) error
	SortedByCategory(categoryId int) ([]*model.Post, error)
	GetIdPosts(userUUID string) ([]int, error)
	GetIdPostsFromCategory(categoryId int) ([]int, error)
}
