package post

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"forum/internal/adapters/api"
	"forum/internal/model"
)

func (h *handlerPost) CreateLikePost(w http.ResponseWriter, r *http.Request) {
	log.Println("POEHALI CREATE Like Post")
	if r.Method != http.MethodPost {
		log.Println("ERROR Method in CreateLikePost")
		api.PrintWebError(w, 405)
		return
	}
	reactionLike := r.FormValue("reaction")

	log.Println("THIS IS REACTION USER -------->", reactionLike)
	postId := r.FormValue("post-id")
	id, err := strconv.Atoi(postId)
	if err != nil {
		log.Printf("ERROR post handler id post invalid value in CreateLikePost method:--->%v\n", err)
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
	checkReaction := strings.TrimSpace(reactionLike)
	if checkReaction == "" || (checkReaction != "like" && checkReaction != "dislike") {
		log.Println("Error Create post function, while reading Title and Body post")
		api.PrintWebError(w, 400)
		return
	}
	authorUUID, err := r.Cookie("My-uuid")
	if err != nil {
		log.Printf("ERROR post handler read cookie in CreateLikePost method:--->%v\n", err)
		h.GetAll(w, r)
		return
	}
	likePostDto := &model.PostLikeAndDislike{
		AuthorUUID: authorUUID.Value,
		PostId:     id,
		Reaction:   reactionLike,
	}
	if err = h.postService.CheckPostLike(h.ctx, likePostDto); err != nil {
		log.Printf("Error CreateLikePost method post handler, when called service create comment:-->%v\n", err)
		h.GetAll(w, r)
		return
	}
	h.GetOnePostById(w, r)
}

func (h *handlerPost) CreateLikeComment(w http.ResponseWriter, r *http.Request) {
	log.Println("POEHALI CREATE Comment Post")
	if r.Method != http.MethodPost {
		log.Println("ERROR Method in CreateLikeComment")
		api.PrintWebError(w, 405)
		return
	}
	reactionLike := r.FormValue("reaction")
	commentId := r.FormValue("comment-id")
	id, err := strconv.Atoi(commentId)
	if err != nil {
		log.Printf("ERROR post handler id post invalid value in CreateLikeComment method:--->%v\n", err)
		api.PrintWebError(w, 400)
		return
	}

	checkReaction := strings.TrimSpace(reactionLike)
	if checkReaction == "" || (checkReaction != "like" && checkReaction != "dislike") {
		log.Println("Error Create post function, while reading Title and Body post")
		api.PrintWebError(w, 400)
		return
	}
	authorUUID, err := r.Cookie("My-uuid")
	if err != nil {
		log.Printf("ERROR post handler read cookie in CreateLikeComment method:--->%v\n", err)
		h.GetAll(w, r)
		return
	}
	likeCommentDto := &model.CommentLikeAndDislike{
		AuthorUUID: authorUUID.Value,
		CommentId:  id,
		Reaction:   reactionLike,
	}
	if err = h.postService.CheckCommentLike(h.ctx, likeCommentDto); err != nil {
		log.Printf("Error CreateLikeComment method post handler, when called service create comment:-->%v\n", err)
		h.GetAll(w, r)
		return
	}
	h.GetOnePostById(w, r)
}
