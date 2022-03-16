package comment

import (
	"context"
	"forum/internal/domain/like"
	"forum/internal/model"
	"log"
)

type service struct {
	storage     CommentStorage
	likeService like.LikeService
}

type CommentService interface {
	GetAll(ctx context.Context, postId int) ([]*model.Comment, error)
	Create(ctx context.Context, comment *model.Comment) error
	Delete(ctx context.Context, comment *model.Comment) error
	GetIdPostsByComment(ctx context.Context, userUUID string) ([]int, error)
	// Update(ctx context.Context, post *model.Comment) error
}

// NewService ...
func NewService(storage CommentStorage, likeService like.LikeService) CommentService {
	return &service{
		storage:     storage,
		likeService: likeService,
	}
}

// func (s *service) Update(ctx context.Context, comment *model.Comment) error {
// 	if err := s.storage.Update(comment); err != nil {
// 		log.Printf("ERROR comments service Update method:---> %v\n", err)
// 		return err
// 	}
// 	return nil
// }

func (s *service) Create(ctx context.Context, comment *model.Comment) error {
	if err := s.storage.Create(comment); err != nil {
		log.Printf("ERROR comment service create comment method, when called comment storage:--->%v\n", err)
		return err
	}
	return nil
}

func (s *service) GetAll(ctx context.Context, postId int) ([]*model.Comment, error) {
	comments, err := s.storage.GetAll(postId)
	if err != nil {
		log.Printf("ERROR comment service GetAll method:--->%v\n", err)
		return nil, err
	}
	for _, val := range comments {
		val.Likes = new(model.CommentLikeAndDislike)
		val.Likes.CommentId = val.Id
		val.Likes.AuthorUUID = val.AuthorUUID
		likesComment, err := s.likeService.GetAll(val.Likes)
		val.Likes = likesComment.CommentLike
		if err != nil {
			log.Printf("ERROR add comments for post get comment method in GetById post:---> %v\n", err.Error())
		}
	}
	return comments, nil
}

func (s *service) Delete(ctx context.Context, comment *model.Comment) error {
	if err := s.storage.Delete(comment); err != nil {
		log.Printf("ERROR comment service Delete method:--> %v\n", err)
		return err
	}
	return nil
}

func (s *service) GetIdPostsByComment(ctx context.Context, userUUID string) ([]int, error) {
	listId, err := s.storage.GetIdPosts(userUUID)
	if err != nil {
		log.Printf("ERROR COMMENT service  (Get Id Posts By Comment method) :----->%v\n", err)
		return nil, err
	}
	return listId, nil
}
