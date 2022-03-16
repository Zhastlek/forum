package api

import (
	"html/template"
	"net/http"
)

type ErrorPage struct {
	Code    int
	Message string
}

// Errors our port
func PrintWebError(w http.ResponseWriter, code int) {
	tmpl, err := template.ParseFiles("template/errors.html", "./template/header.html")
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(code)
	errPage := ErrorPage{Code: code, Message: http.StatusText(code)}
	tmpl.Execute(w, errPage)
	return
}
