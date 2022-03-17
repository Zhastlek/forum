package post

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"forum/internal/adapters/api"
	"forum/internal/domain/category"
	"forum/internal/domain/session"
	"forum/internal/model"
)

type handlerPost struct {
	ctx             context.Context
	postService     api.PostService
	sessionService  session.SessionService
	categoryService category.CategoryService
}

func NewHandler(ctx context.Context, service api.PostService, sessionService session.SessionService, categoryService category.CategoryService) api.Handler {
	return &handlerPost{
		ctx:             ctx,
		postService:     service,
		sessionService:  sessionService,
		categoryService: categoryService,
	}
}

func (h *handlerPost) Register(router *http.ServeMux) {
	router.HandleFunc("/", h.AllPostHandlers)
	router.HandleFunc("/category/", h.SortedPost)
	router.HandleFunc("/posts/", h.PostById)
	router.HandleFunc("/create-post", h.CreatePost)
	router.HandleFunc("/logout", h.Logout)
	router.HandleFunc("/profile", h.ProfileHandler)
	router.HandleFunc("/profile/my-posts", h.MyPostsHandler)
	router.HandleFunc("/profile/my-comments", h.MyCommentsHandler)
	router.HandleFunc("/profile/my-likes", h.MyLikesHandler)

}

func (h *handlerPost) AllPostHandlers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetAll(w, r)
	default:
		log.Println("ERROR Method in home page GetAllPost")
		api.PrintWebError(w, 405)
	}
}

func (h *handlerPost) GetAll(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		api.PrintWebError(w, 404)
		return
	}
	temp, err := template.ParseFiles("./template/main-page.html", "./template/header.html", "./template/navbar.html")

	if err != nil {
		log.Printf("Error main-page html Post Handler GetAll method:--> %v\n", err)
		api.PrintWebError(w, 500)
		return
	}

	post, err := h.postService.GetAll(h.ctx)
	if err != nil {
		log.Printf("ERROR post handler PostCreate method GetById function:--> %v\n", err)
		api.PrintWebError(w, 500)
		return
	}
	sessionStatus, err := h.CheckSession(w, r)
	if err != nil {
		log.Printf("ERROR handler post GetAll method Check Session function %v\n", err)
		sessionStatus = false
	}
	allCategories, err := h.categoryService.GetAll(h.ctx)
	if err != nil {
		log.Printf("ERROR post handler GetAll method: %v\n", err)
		api.PrintWebError(w, 500)
		return
	}

	data := &api.Data{
		Session:    sessionStatus,
		Categories: allCategories,
		Posts:      post,
	}

	err = temp.Execute(w, data)
	if err != nil {
		log.Printf("ERROR post handler GetAll method Execute:---> %v\n", err)
		api.PrintWebError(w, 500)
		return
	}
}

func (h *handlerPost) SortedPost(w http.ResponseWriter, r *http.Request) {
	log.Println("POEHALI SortedPost")
	if r.Method == http.MethodGet || r.Method == http.MethodPost {
		h.SortedByCategory(w, r)
		return
	}
	log.Println("PRIEHALI SortedPost")
	log.Println("ERROR Method in SortedPost")
	api.PrintWebError(w, 405)
}

func (h *handlerPost) SortedByCategory(w http.ResponseWriter, r *http.Request) {
	log.Println("POEHALI SortedByCategory")
	temp, err := template.ParseFiles("./template/main-page.html", "./template/header.html", "./template/navbar.html")
	if err != nil {
		api.PrintWebError(w, 500)
		return
	}
	allCategories, err := h.categoryService.GetAll(h.ctx)
	if err != nil {
		log.Printf("ERROR post handler GetAll method: %v\n", err)
		api.PrintWebError(w, 500)
		return
	}

	valueId := strings.TrimPrefix(r.URL.Path, "/category/")
	id, err := strconv.Atoi(valueId)
	if err != nil {
		api.PrintWebError(w, 400)
		return
	}
	log.Println("ID CATEGORY:------>", id)

	status := api.IsValidCategory(allCategories, id)
	if !status {
		api.PrintWebError(w, 400)
		return
	}

	posts, err := h.postService.SortedByCategory(h.ctx, id)
	if err != nil {
		log.Printf("Error sorted post read sorted service:-->%v\n", err)
		api.PrintWebError(w, 500)
		return
	}

	// for _, val := range posts {
	// 	fmt.Println("result sorted:-->>>>>", val.Title, val.Body, val.CategoryId)
	// }

	sessionStatus, err := h.CheckSession(w, r)
	if err != nil {
		log.Printf("ERROR handler post GetAll method Check Session function %v\n", err)
	}

	data := &api.Data{
		Session:    sessionStatus,
		Categories: allCategories,
		Posts:      posts,
	}
	if err := temp.Execute(w, data); err != nil {
		log.Printf("ERROR post handler SortedByCategory method: %v\n", err)
		api.PrintWebError(w, 500)
		return
	}
	// log.Println("PRIEHALI SortedByCategory")
}

func (h *handlerPost) CheckSession(w http.ResponseWriter, r *http.Request) (bool, error) {
	session, err := r.Cookie("session")
	if err != nil {
		log.Printf("This user is not registered(CheckSession method Handler Post):-->%v\n", err)
		return false, err
	}
	myUUID, err := r.Cookie("My-uuid")
	if err != nil {
		log.Printf("This user is not registered(CheckSession method Handler Post):-->%v\n", err)
		return false, err
	}
	SessionDto := &model.SessionDto{
		MyUUID: myUUID.Value,
		Value:  session.Value,
	}

	if err = h.sessionService.Check(h.ctx, SessionDto); err != nil {
		h.DeleteSession(w, r)
		return false, err
	}

	return true, nil
}

func (h *handlerPost) Logout(w http.ResponseWriter, r *http.Request) {
	session, err := r.Cookie("session")
	if err != nil {
		log.Println("This user is not registered")
		api.PrintWebError(w, 400)
		return
	}
	myUUID, err := r.Cookie("My-uuid")
	if err != nil {
		log.Println("This user is not registered")
		api.PrintWebError(w, 400)
		return
	}
	switch r.Method {
	case http.MethodGet:
		http.SetCookie(w, &http.Cookie{
			Name:    "session",
			MaxAge:  -1,
			Expires: time.Unix(0, 0),
		})
		http.SetCookie(w, &http.Cookie{
			Name:    "My-uuid",
			MaxAge:  -1,
			Expires: time.Unix(0, 0),
		})
		SessionDto := &model.SessionDto{
			MyUUID: myUUID.Value,
			Value:  session.Value,
		}
		if err = h.sessionService.Delete(h.ctx, SessionDto); err != nil {
			log.Printf("ERROR handler post Logout method Delete session entry in DB:--> %v\n", err)
			return
		}
		http.Redirect(w, r, "/", 303)
	default:
		log.Println("ERROR Method in Logout")
		api.PrintWebError(w, 405)
		return
	}
}

func (h *handlerPost) DeleteSession(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie("session")
	if err != nil {
		log.Println("This user is not registered")
		api.PrintWebError(w, 400)
		return
	}
	_, err = r.Cookie("My-uuid")
	if err != nil {
		log.Println("This user is not registered")
		api.PrintWebError(w, 400)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "session",
		MaxAge:  -1,
		Expires: time.Unix(0, 0),
	})
	http.SetCookie(w, &http.Cookie{
		Name:    "My-uuid",
		MaxAge:  -1,
		Expires: time.Unix(0, 0),
	})
	http.Redirect(w, r, "/", 303)
}
