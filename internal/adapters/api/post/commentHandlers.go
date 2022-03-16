package post

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"forum/internal/adapters/api"
	"forum/internal/model"
)

func (h *handlerPost) CreateComment(w http.ResponseWriter, r *http.Request) {
	log.Println("POEHALI CREATE COMMENT")
	if r.Method != http.MethodPost {
		log.Println("ERROR Method in CreateComment")
		api.PrintWebError(w, 405)
		return
	}
	commentBody := r.FormValue("comment")
	postId := r.FormValue("post-id")
	id, err := strconv.Atoi(postId)
	if err != nil {
		log.Printf("ERROR post handler id post invalid value in comment create method:--->%v\n", err)
		api.PrintWebError(w, 400)
		return
	}
	checkComment := strings.TrimSpace(commentBody)
	if checkComment == "" {
		api.PrintWebError(w, 400)
		return
	}
	allPostsId, err := h.postService.GetAll(h.ctx)
	if err != nil {
		log.Printf("ERROR post handler allPostsId method:--> %v\n", err)
		api.PrintWebError(w, 500)
		return
	}
	status := api.IsValidPosts(allPostsId, id)
	if !status {
		api.PrintWebError(w, 400)
		return
	}
	authorUUID, err := r.Cookie("My-uuid")
	if err != nil {
		log.Printf("ERROR post handler read cookie in comment create method:--->%v\n", err)
		h.GetAll(w, r)
		return
	}
	commentDto := &model.CreateCommentDto{
		PostId:            id,
		CommentAuthorUUID: authorUUID.Value,
		CommentBody:       commentBody,
	}
	if err = h.postService.CreateComment(h.ctx, commentDto); err != nil {
		log.Printf("Error create comment method post handler, when called service create comment:-->%v\n", err)
		h.GetAll(w, r)
		return
	}
	h.GetOnePostById(w, r)
}

func (h *handlerPost) DeleteComment(w http.ResponseWriter, r *http.Request) {
	log.Println("POEHALI DELETE COMMENT")
	if r.Method != http.MethodPost {
		log.Println("ERROR Method in DeleteComment")
		api.PrintWebError(w, 405)
		return
	}
	commentBody := r.FormValue("comment")
	postId := r.FormValue("post-id")
	id, err := strconv.Atoi(postId)
	if err != nil {
		log.Printf("ERROR post handler id post invalid value in delete comment method:--->%v\n", err)
		api.PrintWebError(w, 400)
		return
	}
	checkComment := strings.TrimSpace(commentBody)
	if checkComment == "" {
		api.PrintWebError(w, 400)
		return
	}
	authorUUID, err := r.Cookie("My-uuid")
	if err != nil {
		log.Printf("ERROR post handler read cookie in delete comment  method:--->%v\n", err)
		h.GetAll(w, r)
		return
	}
	commentDto := &model.Comment{
		PostId:     id,
		AuthorUUID: authorUUID.Value,
		Body:       commentBody,
	}
	if err = h.postService.DeleteComment(h.ctx, commentDto); err != nil {
		log.Printf("Error create comment method post handler, when called service delete comment:-->%v\n", err)
		h.GetAll(w, r)
		return
	}
	h.GetOnePostById(w, r)
}
