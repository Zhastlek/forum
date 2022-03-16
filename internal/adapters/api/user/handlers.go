package user

import (
	"context"

	"forum/internal/adapters/api"
	"forum/internal/domain/session"
	"forum/internal/model"
	"forum/pkg"
	"html/template"
	"log"
	"net/http"
	"time"

	uuid "github.com/nu7hatch/gouuid"
)

type handlerUser struct {
	ctx            context.Context
	userService    api.UserService
	sessionService session.SessionService
}

func NewHandler(ctx context.Context, service api.UserService, sessionService session.SessionService) api.Handler {
	return &handlerUser{
		ctx:            ctx,
		userService:    service,
		sessionService: sessionService,
	}
}

func (h *handlerUser) Register(router *http.ServeMux) {
	router.HandleFunc("/sign-in", h.SignInHandler)
	router.HandleFunc("/sign-up", h.SignUpHandler)

}

func (h *handlerUser) SignInHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetSignIn(w, r)
	case http.MethodPost:
		h.PostSignIn(w, r)
	default:
		log.Println("Error method SignIn page")
		api.PrintWebError(w, 405)
	}
}

func (h *handlerUser) SignUpHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetSignUp(w, r)
	case http.MethodPost:
		h.PostSignUp(w, r)
	default:
		log.Println("ERROR Method in signup")
		api.PrintWebError(w, 405)
	}
}

func (h *handlerUser) GetSignIn(w http.ResponseWriter, r *http.Request) {
	// temp, err := template.ParseFiles("./templates/login.html")
	temp, err := template.ParseFiles("./template/sign-in.html", "./template/header.html")

	if err != nil {
		log.Printf("Error in parse html files (GetSignIn method handler user):--> %v\n", err)
		api.PrintWebError(w, 500)
		return
	}
	if err = temp.Execute(w, nil); err != nil {
		api.PrintWebError(w, 500)
		return
	}
}

func (h *handlerUser) PostSignIn(w http.ResponseWriter, r *http.Request) {
	// temp, err := template.ParseFiles("./templates/login.html")
	temp, err := template.ParseFiles("./template/sign-in.html", "./template/header.html")
	if err != nil {
		log.Println("Error in parse html files")
		api.PrintWebError(w, 500)
		return
	}
	data := &api.Data{}
	email := r.FormValue("email")
	password := r.FormValue("password")

	user, err := h.userService.GetByEmail(h.ctx, email)
	if err != nil {
		log.Println(("Error user handler sign-in email is wrong"))
		data.Error = true
		w.WriteHeader(403)
		temp.Execute(w, data)
		return
	}

	if user == nil {
		log.Println("Error not found this user in DB")
		data.Error = true
		w.WriteHeader(400)
		temp.Execute(w, data)
		return
	}

	if !pkg.ComparePassword([]byte(user.Password), password) {
		log.Println(("Error user handler sign-in password is wrong"))
		data.Error = true
		w.WriteHeader(403)
		temp.Execute(w, data)
		return
	}
	h.CreateSession(w, r)
	// 	if err = temp.Execute(w, data); err != nil {
	// 		log.Println("-------------------------------------------------------")
	// 		log.Println(err)
	// 		api.PrintWebError(w, 500)
	// 		return
	// 	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *handlerUser) GetSignUp(w http.ResponseWriter, r *http.Request) {
	// temp, err := template.ParseFiles("./templates/signup.html")
	temp, err := template.ParseFiles("./template/sign-up.html", "./template/header.html")

	if err != nil {
		log.Println("Error in parse html files")
	}
	if err = temp.Execute(w, nil); err != nil {
		api.PrintWebError(w, 500)
		return
	}
}

func (h *handlerUser) PostSignUp(w http.ResponseWriter, r *http.Request) {
	temp, err := template.ParseFiles("./template/sign-up.html", "./template/header.html")
	if err != nil {
		log.Println("Error in parse html files")
		api.PrintWebError(w, 500)
		return
	}
	data := &api.Data{}
	name := r.FormValue("login")
	email := r.FormValue("email")
	checkEmail := pkg.IsEmailValid(email)
	if !checkEmail {
		log.Println("Error user handler incorrect Email")
		data.Error = true
		w.WriteHeader(400)
		temp.Execute(w, data)
		return
	}
	password := r.FormValue("password")
	checkPassword := r.FormValue("checkerPassword")

	if password != checkPassword {
		log.Println("Error, please write the correct password")
		data.Error = true
		w.WriteHeader(400)
		temp.Execute(w, data)
		return
	}

	pass, err := pkg.GeneratePassword(password)
	if err != nil {
		log.Println("Error user handler sign-up generate password")
		data.Error = true
		w.WriteHeader(400)
		temp.Execute(w, data)
		return
	}
	userUUID, err := uuid.NewV4()
	if err != nil {
		log.Println("Error user handler generate user UUID return err")
	}

	dto := &model.UserDto{
		Name:     name,
		Email:    email,
		Password: string(pass),
		UUID:     userUUID.String(),
	}
	err = h.userService.Create(h.ctx, dto)
	if err != nil {
		log.Println("Error sign-in handler in user creation")
		api.PrintWebError(w, 400)
		return
	}
	h.CreateSession(w, r)
	// if err = temp.Execute(w, nil); err != nil {
	// 	api.PrintWebError(w, 500)
	// 	return
	// }
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *handlerUser) CheckSession(r *http.Request) (bool, error) {
	session, err := r.Cookie("session")
	if err != nil {
		log.Println("This user is not registered")
		return false, err
	}
	myUUID, err := r.Cookie("My-uuid")
	if err != nil {
		log.Printf("This user is not registered (CheckSession method handler user):--> %v\n", err)
		return false, err
	}
	SessionDto := &model.SessionDto{
		MyUUID: myUUID.Value,
		Value:  session.Value,
	}

	if err = h.sessionService.Check(h.ctx, SessionDto); err != nil {
		return false, err
	}

	return true, nil
}

func (h *handlerUser) CreateSession(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")

	myUuidValue, err := h.userService.GetByEmail(h.ctx, email)
	if err != nil {
		log.Println("Error user handler -> CreateSession incorrect Email")
	}
	sessionDto := &model.SessionDto{
		MyUUID: myUuidValue.UUID,
	}
	if err = h.sessionService.Delete(h.ctx, sessionDto); err != nil {
		log.Printf("ERROR user handler session delete method:---> %v\n", err)
		return
	}
	sessionValue, err := uuid.NewV4()
	if err != nil {
		log.Println("Error user handler generate cookies session value return err")
		return
	}

	expire := time.Now().AddDate(0, 0, 1)
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    sessionValue.String(),
		HttpOnly: true,
		MaxAge:   24 * 60 * 60,
		Expires:  expire,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "My-uuid",
		Value:    myUuidValue.UUID,
		HttpOnly: true,
		MaxAge:   24 * 60 * 60,
		Expires:  expire,
	})
	sessionCreate := &model.Session{
		SessionUUID: sessionValue.String(),
		// Date: date,
		UserUUID: myUuidValue.UUID,
	}
	if err = h.sessionService.Create(h.ctx, sessionCreate); err != nil {
		log.Println("Error: user handler Create session")
		log.Printf("%v\n", err)
		return
	}
}
