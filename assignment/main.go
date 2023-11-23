package main

import (
	"assignment/repository"
	_ "assignment/repository"
	"assignment/server"
	"fmt"
	"net/http"
)

func main() {
	// Create DBs
	repository.Create()
	// APIs
	http.HandleFunc("/login", server.LoginHandler)
	http.HandleFunc("/person/details", server.PersonHandler)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Print(err)
	}
}
