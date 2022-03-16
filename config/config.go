package config

import (
	"context"
	"database/sql"
	"net/http"

	"forum/internal/adapters/db/category"
	"forum/internal/adapters/db/comment"
	"forum/internal/adapters/db/like"
	"forum/internal/adapters/db/post"
	"forum/internal/adapters/db/session"
	"forum/internal/adapters/db/user"

	post3 "forum/internal/adapters/api/post"
	user3 "forum/internal/adapters/api/user"
	category2 "forum/internal/domain/category"
	comment2 "forum/internal/domain/comment"
	like2 "forum/internal/domain/like"
	post2 "forum/internal/domain/post"
	session2 "forum/internal/domain/session"
	user2 "forum/internal/domain/user"
)

func Config(ctx context.Context, db *sql.DB) *http.ServeMux {
	sessionStorage := session.NewStorage(db)
	likeStorage := like.NewLikeStorage(db)
	commentStorage := comment.NewStorage(db)
	categoryStorage := category.NewStorage(db)

	sessionService := session2.NewService(sessionStorage)
	likeService := like2.NewService(likeStorage)
	commentService := comment2.NewService(commentStorage, likeService)
	categoryService := category2.NewService(categoryStorage)

	router := http.NewServeMux()
	postStorage := post.NewStorage(db)
	postService := post2.NewService(postStorage, commentService, likeService)
	postHandler := post3.NewHandler(ctx, postService, sessionService, categoryService)

	postHandler.Register(router)

	userStorage := user.NewStorage(db)
	userService := user2.NewService(userStorage)
	userHandler := user3.NewHandler(ctx, userService, sessionService)
	userHandler.Register(router)

	return router
}
