package post

import (
	"context"
	"fmt"

	"forum/internal/adapters/api"
	"forum/internal/domain/comment"
	"forum/internal/domain/like"
	"forum/internal/model"
	"log"
)

type service struct {
	storage        PostStorage
	commentService comment.CommentService
	likeService    like.LikeService
}

// NewService ...
func NewService(storage PostStorage, commentService comment.CommentService, likeService like.LikeService) api.PostService {
	return &service{
		storage:        storage,
		commentService: commentService,
		likeService:    likeService,
	}
}

func (s *service) Create(ctx context.Context, dto *model.CreatePostDto) (int, error) {
	post := &model.Post{
		Title:      dto.Title,
		Body:       dto.Body,
		AuthorUUID: dto.AuthorUUID,
		CategoryId: dto.CategoryId,
	}
	log.Println(post)
	infoPost, err := s.storage.Create(post)
	if err != nil {
		log.Printf("ERROR post service Create method:--> %v\n", err)
		return infoPost, err
	}

	for _, value := range post.CategoryId {
		err := s.storage.CreatePostCategory(infoPost, value)
		if err != nil {
			log.Printf("ERROR post service CreatePostCategory method:--> %v\n", err)
			return infoPost, err
		}
	}
	log.Println(infoPost)
	return infoPost, nil
}

func (s *service) Update(ctx context.Context, dtoPost *model.UpdatePostDTO) error {
	p := &model.Post{
		Id:         dtoPost.Id,
		Title:      dtoPost.Title,
		Body:       dtoPost.Body,
		CategoryId: dtoPost.CategoryId,
	}
	err := s.storage.Update(p)
	if err != nil {
		log.Printf("Error post service update returned :---> %v\n", err)
		return err
	}
	if err = s.storage.DeletePostCategory(p.Id); err != nil {
		log.Printf("Error post service UPDATE POST (DeletePostCategory) returned :---> %v\n", err)
		return err
	}

	for _, value := range p.CategoryId {
		err := s.storage.CreatePostCategory(p.Id, value)
		if err != nil {
			log.Printf("ERROR post service CreatePostCategory method:--> %v\n", err)
			return err
		}
	}

	return nil
}

func (s *service) GetByID(ctx context.Context, id int) ([]*model.Post, error) {
	// s.storage.GetOne(id)
	post, err := s.storage.GetOne(id)
	if err != nil {
		log.Printf("ERROR post service GetById method :--> %v\n", err)
		return post, err
	}
	for _, val := range post {
		val.Likes = new(model.PostLikeAndDislike)
		val.Likes.PostId = val.Id
		val.Likes.AuthorUUID = val.AuthorUUID
		likesPost, err := s.likeService.GetAll(val.Likes)
		val.Likes = likesPost.PostLike
		val.Comments, err = s.GetComment(ctx, id)
		if err != nil {
			log.Printf("ERROR add comments for post get comment method in GetById post:---> %v\n", err.Error())
		}
	}
	return post, nil
}

func (s *service) GetAll(ctx context.Context) ([]*model.Post, error) {
	p, err := s.storage.GetAll()
	if err != nil {
		log.Printf("ERROR post service GetAll method :--> %v\n", err)
		return p, err
	}
	log.Println("GetAll method post service")
	log.Println(p)
	for _, val := range p {
		fmt.Println(val)
	}
	return p, nil
}

func (s *service) SortedByCategory(ctx context.Context, categoryId int) ([]*model.Post, error) {
	log.Println("POEHALI SortedByCategory service post")
	// listId, err := s.storage.GetIdPostsFromCategory(categoryId)
	// if err != nil {
	// 	log.Printf("Error post service GetIdPostsFromCategory returned :---> %v\n", err)
	// 	return nil, err
	// }

	posts, err := s.storage.SortedByCategory(categoryId)
	if err != nil {
		log.Printf("Error post service sorted by category returned :---> %v\n", err)
		return posts, err
	}
	log.Println("PRIEHALI SortedByCategory  service post ")
	return posts, nil
}

func (s *service) DeletePost(ctx context.Context, id int) error {
	if err := s.storage.Delete(id); err != nil {
		log.Printf("ERROR post service delete post method :--> %v\n", err)
		return err
	}

	return nil
}

func (s *service) CreateComment(ctx context.Context, dto *model.CreateCommentDto) error {
	comment := &model.Comment{
		PostId:     dto.PostId,
		Body:       dto.CommentBody,
		AuthorUUID: dto.CommentAuthorUUID,
	}
	if err := s.commentService.Create(ctx, comment); err != nil {
		log.Printf("ERROR post service create comment method, when called comment service:--->%v\n", err)
		return err
	}
	return nil
}

func (s *service) GetComment(ctx context.Context, postId int) ([]*model.Comment, error) {
	comments, err := s.commentService.GetAll(ctx, postId)
	if err != nil {
		log.Printf("ERROR post service GetComment method:----> %v\n", err)
		return nil, err
	}
	return comments, nil
}

func (s *service) DeleteComment(ctx context.Context, comment *model.Comment) error {
	if err := s.commentService.Delete(ctx, comment); err != nil {
		log.Printf("ERROR post service DELETEComment method:----> %v\n", err)
		return err
	}
	return nil
}

func (s *service) CheckPostLike(ctx context.Context, like *model.PostLikeAndDislike) error {
	if err := s.likeService.Check(like); err != nil {
		log.Printf("ERROR post service Check Post Like method:----> %v\n", err)
		return err
	}
	return nil
}

func (s *service) CheckCommentLike(ctx context.Context, like *model.CommentLikeAndDislike) error {
	if err := s.likeService.Check(like); err != nil {
		log.Printf("ERROR post service Check Comment Like method:----> %v\n", err)
		return err
	}
	return nil
}

func (s *service) SortedByPost(ctx context.Context, userUUID string) ([]*model.Post, error) {
	listId, err := s.storage.GetIdPosts(userUUID)
	posts := []*model.Post{}
	if err != nil {
		log.Printf("ERROR Post service SortedByPost method:----->%v\n", err)
		return nil, err
	}
	for _, value := range listId {
		onePost, err := s.GetByID(ctx, value)
		if err != nil {
			log.Printf("ERROR Post service Get By Id func in (SortedByPost method):----->%v\n", err)
			return nil, err
		}
		posts = append(posts, onePost...)
	}
	return posts, nil
}

func (s *service) SortedByComment(ctx context.Context, userUUID string) ([]*model.Post, error) {
	listId, err := s.commentService.GetIdPostsByComment(ctx, userUUID)
	if err != nil {
		log.Printf("ERROR Post service SortedByComment method:----->%v\n", err)
		return nil, err
	}
	posts := []*model.Post{}
	for _, value := range listId {
		onePost, err := s.GetByID(ctx, value)
		if err != nil {
			log.Printf("ERROR Post service Get By Id func in (SortedByComment method):----->%v\n", err)
			return nil, err
		}
		posts = append(posts, onePost...)
	}
	return posts, nil
}

func (s *service) SortedByLike(ctx context.Context, userUUID string) ([]*model.Post, error) {
	listId, err := s.likeService.GetIdPost(userUUID)
	if err != nil {
		log.Printf("ERROR POST SERVICE IN (Sorted By Like) METHOD:---->%v\n", err)
		return nil, err
	}
	fmt.Println("slice listId in SORTED BY LIKE METHOD------------>", listId)
	posts := []*model.Post{}
	fmt.Println("len posts in SORTED BY LIKE after (((FOR)))------------>", len(posts))
	for i := len(listId) - 1; i >= 0; i-- {
		fmt.Println("ELEMENT SLICE listId SORTED BY LIKE------------>", i)
		onePost, err := s.GetByID(ctx, listId[i])
		if err != nil {
			log.Printf("ERROR Post service Get By Id func in (SortedByLike method):----->%v\n", err)
			return nil, err
		}
		posts = append(posts, onePost...)
	}
	fmt.Println("len posts in SORTED BY LIKE------------>", len(posts))
	return posts, nil
}
