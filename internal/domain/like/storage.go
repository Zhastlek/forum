package like

import "forum/internal/model"

type LikeStorage interface {
	GetOneCommentLike(LikeAndDisLike *model.CommentLikeAndDislike) ([]*model.CommentLikeAndDislike, error)
	CreateCommentLike(LikeAndDisLike *model.CommentLikeAndDislike) error
	GetAllCommentLike(LikeAndDisLike *model.CommentLikeAndDislike) (*model.CommentLikeAndDislike, error)
	DeleteCommentLike(LikeAndDisLike *model.CommentLikeAndDislike) error
	UpdateCommentLike(LikeAndDisLike *model.CommentLikeAndDislike) error

	CreatePostLike(LikeAndDisLike *model.PostLikeAndDislike) error
	GetAllPostLike(LikeAndDisLike *model.PostLikeAndDislike) (*model.PostLikeAndDislike, error)
	DeletePostLike(LikeAndDisLike *model.PostLikeAndDislike) error
	UpdatePostLike(LikeAndDisLike *model.PostLikeAndDislike) error
	GetOnePostLike(LikeAndDisLike *model.PostLikeAndDislike) ([]*model.PostLikeAndDislike, error)

	GetIdPostByLike(userUUID string) ([]int, error)
}
