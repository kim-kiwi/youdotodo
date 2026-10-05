package pages

import (
	"fmt"
	"io"
)

type Page struct {
	Link  string
	Title string
	Text  string
}

var Pages = [...]Page{
	{"/", "Home", "home"},
	{"/signup", "Sign Up", "signup"},
	{"/login", "Login", "login"},
}

const api = "http://localhost:8080"

func Echo(w io.Writer) func(...any) {
	return func(a ...any) {
		fmt.Fprint(w, a...)
	}
}
