package post

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"

	"forum/internal/adapters/api"
	"forum/internal/model"
)

func (h *handlerPost) PostById(w http.ResponseWriter, r *http.Request) {
	fmt.Println("-------->POEHALI POSTBYUUID")
	fmt.Println("This Is Method One Post----------------------->", r.Method)
	switch r.Method {
	case http.MethodGet:
		h.GetOnePostById(w, r)
	case http.MethodPost:
		fmt.Println("-------->POEHALI POST")
		h.EditOnePost(w, r)
	default:
		fmt.Println("-------->POEHALI default")
		log.Println("ERROR Method in GetOnePost")
		api.PrintWebError(w, 405)
	}
}

func (h *handlerPost) GetOnePostById(w http.ResponseWriter, r *http.Request) {
	fmt.Println("-------->PRIEHALI POSTBYUUID GetOnePostById")
	temp, err := template.ParseFiles("./template/one-post.html", "./template/header.html", "./template/navbar.html")
	if err != nil {
		log.Printf("ERROR post handler GetOnePostById parse html files:--->%v\n", err)
		api.PrintWebError(w, 500)
		return
	}
	// value := r.FormValue("id")
	fmt.Println("parseform---->", r.ParseForm(), r.URL.Path)

	value := strings.TrimPrefix(r.URL.Path, "/posts/")

	id, err := strconv.Atoi(value)
	if err != nil {
		log.Printf("ERROR post handler GetOnePostById form value id:--->%v\n", err)
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

	log.Println("NUMBER POST -------------------->", value)
	post, err := h.postService.GetByID(h.ctx, id)
	if err != nil {
		log.Printf("ERROR post handler GetOnePost method:--> %v\n", err)
		api.PrintWebError(w, 400)
		return
	}
	sessionStatus, err := h.CheckSession(w, r)
	if err != nil {
		log.Printf("ERROR handler post GetOnePost method Check Session function %v\n", err)
	}
	allCategories, err := h.categoryService.GetAll(h.ctx)
	if err != nil {
		log.Printf("ERROR post handler GetOnePost method categoryService.GetAll function : %v\n", err)
		api.PrintWebError(w, 500)
		return
	}
	data := &api.Data{
		Session:    sessionStatus,
		Categories: allCategories,
		Posts:      post,
	}

	myUUID, err := r.Cookie("My-uuid")
	if err != nil {
		log.Printf("This user is not sign in r.Cookie(my-uuid):-->%v\n", err)
	} else if len(post) != 0 {
		for _, value := range post {
			if myUUID.Value == value.AuthorUUID {
				data.MyPost = true
			}
			if len(value.Comments) != 0 || value.Comments != nil {
				MarkYourData(myUUID.Value, value.Comments)
				for _, val := range value.Comments {
					if val.Likes != nil {
						MarkYourData(myUUID.Value, val.Likes)
					}
				}
			}
			if value.Likes != nil {
				MarkYourData(myUUID.Value, value.Likes)
			}
		}
	} else if len(post) == 0 {
		api.PrintWebError(w, 400)
		return
	}
	log.Println("STATUS MYPOST:--->>>", data.MyPost)

	fmt.Println("-------->PRIEHALI---> GetOnePostById")
	if err = temp.Execute(w, data); err != nil {
		log.Printf("ERROR post handler GetOnePostById method:--> %v\n", err)
		api.PrintWebError(w, 500)
		return
	}
}

func (h *handlerPost) EditOnePost(w http.ResponseWriter, r *http.Request) {
	log.Println("ADD ONE POST METHOD----->", r.Method)
	log.Println("ADD ONE POST FORMVALUE --------->", r.FormValue("_method"))
	switch r.FormValue("_method") {
	case "update-post":
		log.Println("POEHALI put method post")
		h.UpdatePost(w, r)
	case "delete-post":
		log.Println("POEHALI delete method post")
		h.DeletePost(w, r)
	case "create-comment":
		h.CreateComment(w, r)
	case "delete-comment":
		h.DeleteComment(w, r)
	case "like-post":
		h.CreateLikePost(w, r)
	case "like-comment":
		h.CreateLikeComment(w, r)
	default:
		log.Println("ERROR METHOD add one post function")
		api.PrintWebError(w, 405)
	}
}

func (h *handlerPost) CreatePost(w http.ResponseWriter, r *http.Request) {
	log.Println("POEHALI")
	switch r.Method {
	case http.MethodGet:
		h.GetCreatePost(w, r)
	case http.MethodPost:
		h.PostCreatePost(w, r)
	default:
		log.Println("ERROR Method in CreatePost")
		api.PrintWebError(w, 405)
	}

}

func (h *handlerPost) GetCreatePost(w http.ResponseWriter, r *http.Request) {
	log.Println("PRIEHALI")
	temp, err := template.ParseFiles("./template/create-post.html", "./template/header.html", "./template/navbar.html")
	if err != nil {
		log.Printf("ERROR post handler parse html file GetCreatePost method:--> %v\n", err)
		api.PrintWebError(w, 500)
		return
	}
	sessionStatus, err := h.CheckSession(w, r)
	if err != nil {
		log.Printf("ERROR handler post GetCreatePost method Check Session function %v\n", err)
		api.PrintWebError(w, 403)
		return
	}
	allCategories, err := h.categoryService.GetAll(h.ctx)
	if err != nil {
		log.Printf("ERROR post handler GetCreatePost method allCategories function: %v\n", err)
		api.PrintWebError(w, 500)
		return
	}

	data := &api.Data{
		Session:    sessionStatus,
		Categories: allCategories,
	}
	if err = temp.Execute(w, data); err != nil {
		log.Printf("ERROR post handler GetCreatePost method:--> %v\n", err)
		api.PrintWebError(w, 500)
	}

}

func (h *handlerPost) PostCreatePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		api.PrintWebError(w, 405)
		return
	}
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

	titlePost := r.FormValue("title")
	bodyPost := r.FormValue("bodyPost")

	err = r.ParseForm() // парсим данные в r
	if err != nil {
		log.Printf("Error Create post function, while reading Category post:-->%v\n", err)
		api.PrintWebError(w, 400)
		return
	}
	categoryPost := r.PostForm["category"] // наши данные

	checkTitle := strings.TrimSpace(titlePost)
	checkBody := strings.TrimSpace(bodyPost)
	if checkTitle == "" || checkBody == "" {
		log.Println("Error Create post function, while reading Title and Body post")
		api.PrintWebError(w, 400)
		return
	}

	var categories []int
	for _, checkCategory := range categoryPost {
		categoryId, err := strconv.Atoi(checkCategory)
		log.Println("CATEGORY--------------------->", categoryId)
		if err != nil {
			log.Printf("Error category Id read post create method :-->%v\n", err)
			api.PrintWebError(w, 400)
			return
		}
		status := api.IsValidCategory(allCategories, categoryId)
		if status {
			categories = append(categories, categoryId)
		} else {
			api.PrintWebError(w, 400)
			return
		}
	}

	authorUUID, err := r.Cookie("My-uuid")
	if err != nil {
		log.Println("Error Create post function, while reading cookie")
		api.PrintWebError(w, 403)
		return
	}
	postDto := &model.CreatePostDto{
		Title:      titlePost,
		Body:       bodyPost,
		AuthorUUID: authorUUID.Value,
		CategoryId: categories,
	}

	postId, err := h.postService.Create(h.ctx, postDto)
	if err != nil {
		log.Printf("ERROR post handler PostCreate method:--> %v\n", err)
		api.PrintWebError(w, 400)
		return
	}
	post, err := h.postService.GetByID(h.ctx, postId)
	if err != nil {
		log.Printf("ERROR post handler PostCreate method GetById function:--> %v\n", err)
		api.PrintWebError(w, 500)
		return
	}
	sessionStatus, err := h.CheckSession(w, r)
	if err != nil {
		log.Printf("ERROR handler post GetAll method Check Session function %v\n", err)
		api.PrintWebError(w, 400)
		return
	}

	data := &api.Data{
		Session:    sessionStatus,
		Categories: allCategories,
		Posts:      post,
	}

	if err = temp.Execute(w, data); err != nil {
		log.Printf("ERROR post handler PostCreatePost method:--> %v\n", err)
		api.PrintWebError(w, 500)
		return
	}
}

func (h *handlerPost) UpdatePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Println("ERROR incorrect method in UpdatePost function")
		api.PrintWebError(w, 405)
		return
	}
	allCategories, err := h.categoryService.GetAll(h.ctx)
	if err != nil {
		log.Printf("ERROR post handler GetAll method: %v\n", err)
		api.PrintWebError(w, 500)
		return
	}

	valueTitle := r.FormValue("update-title")
	valueBody := r.FormValue("update-body")

	checkTitle := strings.TrimSpace(valueTitle)
	checkBody := strings.TrimSpace(valueBody)
	if checkTitle == "" || checkBody == "" {
		log.Println("Error Create post function, while reading Title and Body post")
		api.PrintWebError(w, 400)
		return
	}

	log.Println("--->", r.FormValue("id"))
	value := r.FormValue("id")
	id, err := strconv.Atoi(value)
	if err != nil {
		log.Printf("ERROR post handler UpdatePost form value id:--->%v\n", err)
		api.PrintWebError(w, 400)
		return
	}

	err = r.ParseForm() // парсим данные в r
	if err != nil {
		log.Println("Error Create post function, while reading Category post")
		api.PrintWebError(w, 400)
		return
	}
	categoryPost := r.PostForm["category"] // наши данные

	var categories []int
	for _, checkCategory := range categoryPost {
		categoryId, err := strconv.Atoi(checkCategory)
		log.Println("CATEGORY--------------------->", categoryId)
		if err != nil {
			log.Printf("Error category Id read post update method :-->%v\n", err)
			api.PrintWebError(w, 400)
			return
		}
		status := api.IsValidCategory(allCategories, categoryId)
		if status {
			categories = append(categories, categoryId)
		} else {
			api.PrintWebError(w, 400)
			return
		}
	}

	postDto := &model.UpdatePostDTO{
		Id:         id,
		Title:      valueTitle,
		Body:       valueBody,
		CategoryId: categories,
	}
	if err = h.postService.Update(h.ctx, postDto); err != nil {
		log.Printf("ERROR post handler Update method :--->%v\n", err)
		api.PrintWebError(w, 400)
		return
	}
	h.GetOnePostById(w, r)
}

func (h *handlerPost) DeletePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Println("ERROR incorrect method in DeletePost function")
		h.GetAll(w, r)
		return
	}
	log.Println("delete method post R URL---->", r.URL.Path)
	value := r.FormValue("id")
	id, err := strconv.Atoi(value)
	if err != nil {
		log.Printf("ERROR post handler Delete post form value id:--->%v\n", err)
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
	if err = h.postService.DeletePost(h.ctx, id); err != nil {
		log.Printf("ERROR in DeletePost method order by delete postservice:--> %v\n", err)
		api.PrintWebError(w, 400)
		return
	}
	http.Redirect(w, r, "/", 303)

}

func MarkYourData(uuid string, elements interface{}) {
	switch slice := elements.(type) {
	case []*model.Comment:
		for _, value := range slice {
			if value.AuthorUUID == uuid {
				value.MyComment = true
			}
		}
	case *model.PostLikeAndDislike:
		if slice.AuthorUUID == uuid {
			slice.MyReaction = true
		}
	case *model.CommentLikeAndDislike:
		if slice.AuthorUUID == uuid {
			slice.MyReaction = true
		}
	}
}
