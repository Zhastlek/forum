package main

import (
	"context"
	"log"
	"net/http"

	"forum/config"
	"forum/dbase"
)

func main() {
	ctx := context.Background()
	db := dbase.CheckDB()
	router := config.Config(ctx, db)

	router.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

	log.Println("port: 8080 is listening")
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		log.Printf("%v", err)
		// log.Fatal("ListenAndServe ERROR")
	}
}
