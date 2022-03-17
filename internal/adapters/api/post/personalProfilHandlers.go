package post

import (
	"html/template"
	"log"
	"net/http"

	"forum/internal/adapters/api"
)

func (h *handlerPost) ProfileHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("POEHALI Profile Handler------>", r.Method)
	if r.Method != http.MethodGet {
		api.PrintWebError(w, 405)
		return
	}
	temp, err := template.ParseFiles("./template/profile.html", "./template/header.html", "./template/navbar.html")
	if err != nil {
		log.Printf("ERROR post handler ProfileHandler parse html files:--->%v\n", err)
		api.PrintWebError(w, 500)
		return
	}
	sessionStatus, err := h.CheckSession(w, r)
	if err != nil {
		log.Printf("ERROR post handler ProfileHandler method Check Session function %v\n", err)
		api.PrintWebError(w, 403)
		return
	}
	allCategories, err := h.categoryService.GetAll(h.ctx)
	if err != nil {
		log.Printf("ERROR post handler MyPostsHandler method categoryService.GetAll function : %v\n", err)
		api.PrintWebError(w, 500)
		return
	}
	data := &api.Data{
		Session:    sessionStatus,
		Categories: allCategories,
	}

	temp.Execute(w, data)
}

func (h *handlerPost) MyPostsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("POEHALI Profile MyPostsHandler Handler------>", r.Method)
	if r.Method != http.MethodGet {
		api.PrintWebError(w, 405)
		return
	}
	temp, err := template.ParseFiles("./template/profile.html", "./template/header.html", "./template/navbar.html")
	if err != nil {
		log.Printf("ERROR post handler ProfileHandler parse html files:--->%v\n", err)
		api.PrintWebError(w, 500)
		return
	}
	sessionStatus, err := h.CheckSession(w, r)
	if err != nil {
		log.Printf("ERROR post handler ProfileHandler method Check Session function %v\n", err)
		api.PrintWebError(w, 403)
		return
	}
	myUUID, err := r.Cookie("My-uuid")
	if err != nil {
		log.Printf("ERROR post handler ProfileHandler method Read my-uuid function %v\n", err)
		api.PrintWebError(w, 403)
		return
	}

	posts, err := h.postService.SortedByPost(h.ctx, myUUID.Value)
	if err != nil {
		log.Printf("ERROR Personal profile post handler (sorted by post function):---->%v\n", err)
		api.PrintWebError(w, 500)
		return
	}
	allCategories, err := h.categoryService.GetAll(h.ctx)
	if err != nil {
		log.Printf("ERROR post handler MyPostsHandler method categoryService.GetAll function : %v\n", err)
		api.PrintWebError(w, 500)
		return
	}
	data := &api.Data{
		Session:    sessionStatus,
		Categories: allCategories,
		Posts:      posts,
	}

	for _, value := range posts {
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

	log.Println("STATUS MYPOST:--->>>", data.MyPost)

	w.WriteHeader(200)
	log.Println("-------->PRIEHALI---> Profile MyPostsHandler Handler")
	if err = temp.Execute(w, data); err != nil {
		log.Printf("ERROR Profile MyPostsHandler method:--> %v\n", err)
		api.PrintWebError(w, 500)
		return
	}
}

func (h *handlerPost) MyCommentsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("POEHALI Profile MyCommentsHandler Handler------>", r.Method)
	if r.Method != http.MethodGet {
		api.PrintWebError(w, 405)
		return
	}
	temp, err := template.ParseFiles("./template/profile.html", "./template/header.html", "./template/navbar.html")
	if err != nil {
		log.Printf("ERROR post handler (MyCommentsHandler) ProfileHandler parse html files:--->%v\n", err)
		api.PrintWebError(w, 500)
		return
	}
	sessionStatus, err := h.CheckSession(w, r)
	if err != nil {
		log.Printf("ERROR (MyCommentsHandler) post handler ProfileHandler method Check Session function %v\n", err)
		api.PrintWebError(w, 403)
		return
	}
	myUUID, err := r.Cookie("My-uuid")
	if err != nil {
		log.Printf("ERROR post handler ProfileHandler method Read my-uuid function %v\n", err)
		api.PrintWebError(w, 403)
		return
	}

	posts, err := h.postService.SortedByComment(h.ctx, myUUID.Value)
	if err != nil {
		log.Printf("ERROR Personal profile post handler (MyCommentsHandler function):---->%v\n", err)
		api.PrintWebError(w, 500)
		return
	}
	allCategories, err := h.categoryService.GetAll(h.ctx)
	if err != nil {
		log.Printf("ERROR post handler MyCommentsHandler method categoryService.GetAll function : %v\n", err)
		api.PrintWebError(w, 500)
		return
	}
	data := &api.Data{
		Session:    sessionStatus,
		Categories: allCategories,
		Posts:      posts,
	}

	for _, value := range posts {
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

	log.Println("STATUS MYPOST:--->>>", data.MyPost)

	w.WriteHeader(200)
	log.Println("-------->PRIEHALI---> Profile MyCommentsHandler Handler")
	if err = temp.Execute(w, data); err != nil {
		log.Printf("ERROR Profile MyCommentsHandler method:--> %v\n", err)
		api.PrintWebError(w, 500)
		return
	}
}

func (h *handlerPost) MyLikesHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("POEHALI Profile MyLikesHandler Handler------>", r.Method)
	if r.Method != http.MethodGet {
		api.PrintWebError(w, 405)
		return
	}
	temp, err := template.ParseFiles("./template/profile.html", "./template/header.html", "./template/navbar.html")
	if err != nil {
		log.Printf("ERROR post handler (MyLikesHandler) ProfileHandler parse html files:--->%v\n", err)
		api.PrintWebError(w, 500)
		return
	}
	sessionStatus, err := h.CheckSession(w, r)
	if err != nil {
		log.Printf("ERROR post handler (MyLikesHandler) ProfileHandler method Check Session function %v\n", err)
		api.PrintWebError(w, 403)
		return
	}
	myUUID, err := r.Cookie("My-uuid")
	if err != nil {
		log.Printf("ERROR post handler (MyLikesHandler) ProfileHandler method Read my-uuid function %v\n", err)
		api.PrintWebError(w, 403)
		return
	}

	posts, err := h.postService.SortedByLike(h.ctx, myUUID.Value)
	if err != nil {
		log.Printf("ERROR Personal profile post handler (MyLikesHandler function):---->%v\n", err)
		api.PrintWebError(w, 500)
		return
	}
	allCategories, err := h.categoryService.GetAll(h.ctx)
	if err != nil {
		log.Printf("ERROR post handler  (MyLikesHandler) method categoryService.GetAll function : %v\n", err)
		api.PrintWebError(w, 500)
		return
	}
	data := &api.Data{
		Session:    sessionStatus,
		Categories: allCategories,
		Posts:      posts,
	}

	for _, value := range posts {
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

	log.Println("STATUS MYPOST:--->>>", data.MyPost)

	w.WriteHeader(200)
	log.Println("-------->PRIEHALI---> Profile  (MyLikesHandler) Handler")
	if err = temp.Execute(w, data); err != nil {
		log.Printf("ERROR Profile (MyLikesHandler)  method:--> %v\n", err)
		api.PrintWebError(w, 500)
		return
	}
}
