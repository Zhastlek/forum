package api

import "forum/internal/model"

type Data struct {
	Session       bool
	Categories    []*model.Category
	Posts         []*model.Post
	MyPost        bool
	MyComment     bool
	MyReaction    bool
	LikeOrDislike int
	Error         bool
}
