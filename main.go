package main

import (
	"fmt"
	"net/http"
	"youdotodo/ghost"
	"youdotodo/pages"
	. "youdotodo/repo"
)

type contextKey int

const (
	ctxUserKey contextKey = iota
)

func main() {
	mux := http.NewServeMux()
	sm := NewSessionManager()

	uc, err := NewUserCollection()
	if err != nil {
		panic(err)
	}
	defer uc.Close()
	tc, err := NewTodoCollection()
	if err != nil {
		panic(err)
	}
	defer tc.Close()

	userApi := ghost.NewApi("users", uc)
	ghost.AddApiToMux(mux, ghost.EndpointCGL, userApi)
	todoApi := NewTodoApi(tc, sm, uc)
	ghost.AddApiToMux(mux, ghost.EndpointAll, todoApi)

	mux.HandleFunc("GET /{$}", indexHandler(sm, tc, uc))
	mux.HandleFunc("GET /signup", pages.SignUp)
	mux.HandleFunc("GET /login", pages.Login)
	mux.HandleFunc("POST /login", loginHandler(sm, uc))
	mux.HandleFunc("/logout", logoutHandler(sm, uc))

	fs := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	fmt.Println("Server starting on port 8080...")

	if err := http.ListenAndServe(":8080", mux); err != nil {
		fmt.Printf("Server failed to start: %v\n", err)
	}
}
