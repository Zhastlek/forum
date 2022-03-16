package api

import "forum/internal/model"

func IsValidCategory(categories []*model.Category, num int) bool {
	var status bool
	for _, value := range categories {
		if num == value.Id {
			status = true
		}
	}
	return status
}

func IsValidPosts(posts []*model.Post, num int) bool {
	var status bool
	for _, value := range posts {
		if num == value.Id {
			status = true
		}
	}
	return status
}
