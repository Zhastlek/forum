package like

import (
	"errors"
	"forum/internal/model"
	"log"
)

type service struct {
	storage LikeStorage
}

type LikeService interface {
	Check(like interface{}) error
	GetAll(like interface{}) (*model.Like, error)
	GetIdPost(userUUID string) ([]int, error)
}

// NewService ...
func NewService(storage LikeStorage) LikeService {
	return &service{storage: storage}
}

func (s *service) GetIdPost(userUUID string) ([]int, error) {
	listId, err := s.storage.GetIdPostByLike(userUUID)
	if err != nil {
		log.Printf("ERROR LIKE SERVICE IN GET ID POST METHOD:---->%v\n", err)
		return listId, err
	}
	return listId, nil
}

func (s *service) GetAll(like interface{}) (*model.Like, error) {
	dataLike := &model.Like{}
	var err error
	switch value := like.(type) {
	case *model.PostLikeAndDislike:
		dataLike.PostLike, err = s.storage.GetAllPostLike(value)
		if err != nil {
			log.Printf("ERROR like service GetAll post like method:---->%v\n", err)
			return dataLike, err
		}
		return dataLike, nil
	case *model.CommentLikeAndDislike:
		dataLike.CommentLike, err = s.storage.GetAllCommentLike(value)
		if err != nil {
			log.Printf("ERROR like service GetAll comment like method:---->%v\n", err)
			return dataLike, err
		}
		return dataLike, nil
	default:
		err := errors.New("ERROR LikeService GetAll method read undeclared method")
		return nil, err
	}
}

func (s *service) Check(like interface{}) error {
	switch value := like.(type) {
	case *model.PostLikeAndDislike:
		if err := s.checkPostlike(value); err != nil {
			log.Printf("ERROR like service Check postlike method:---> %v\n", err)
			return err
		}
	case *model.CommentLikeAndDislike:
		if err := s.checkCommentlike(value); err != nil {
			log.Printf("ERROR like service Check comment like method:---> %v\n", err)
			return err
		}
	}
	return nil
}

func (s *service) checkPostlike(like *model.PostLikeAndDislike) error {
	status, err := s.storage.GetOnePostLike(like)
	if err != nil {
		log.Printf("ERROR like service checkPostLike method:---> %v\n", err)
		return err
	}
	log.Println("THIS IS LEN POST LIKE SLICE----->", len(status))
	if len(status) == 1 {
		if status[0].Reaction == like.Reaction {
			if err := s.storage.DeletePostLike(like); err != nil {
				log.Printf("ERROR PostLike service delete like method:--->%v\n", err)
				return err
			}
		}
		if status[0].Reaction != like.Reaction {
			if err := s.storage.UpdatePostLike(like); err != nil {
				log.Printf("ERROR PostLike service update like method:--->%v\n", err)
				return err
			}
		}
	} else if len(status) == 0 {
		if err := s.storage.CreatePostLike(like); err != nil {
			log.Printf("ERROR like service create postlike method:---> %v\n", err)
			return err
		}
	}
	return nil
}

func (s *service) checkCommentlike(like *model.CommentLikeAndDislike) error {
	status, err := s.storage.GetOneCommentLike(like)
	if err != nil {
		log.Printf("ERROR like service checkCommentlike method:---> %v\n", err)
		return err
	}

	if len(status) == 1 {
		if status[0].Reaction == like.Reaction {
			if err := s.storage.DeleteCommentLike(like); err != nil {
				log.Printf("ERROR CommentLike service delete like method:--->%v\n", err)
				return err
			}
		}
		if status[0].Reaction != like.Reaction {
			if err := s.storage.UpdateCommentLike(like); err != nil {
				log.Printf("ERROR CommentLike service update like method:--->%v\n", err)
				return err
			}
		}
	} else if len(status) == 0 {
		if err := s.storage.CreateCommentLike(like); err != nil {
			log.Printf("ERROR like service create CommentLike method:---> %v\n", err)
			return err
		}
	}
	return nil
}
