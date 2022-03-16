package comment

import "forum/internal/model"

type CommentStorage interface {
	// GetOne(id int) *model.Comment
	GetAll(postId int) ([]*model.Comment, error)
	Create(post *model.Comment) error
	Delete(post *model.Comment) error
	GetIdPosts(userUUID string) ([]int, error)
	// Update(post *model.Comment) error
}
