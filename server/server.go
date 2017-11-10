package server

import (
	"github.com/go-martini/martini"
	"html/template"
	"net/http"
)

type Server struct{}

type Da struct {
	Title string
	Name  string
}

func (server *Server) Run() {
	m := martini.Classic()
	m.Use(martini.Static("server/public"))

	m.Get("/", func(w http.ResponseWriter, r *http.Request) {
		tpl := template.Must(template.ParseFiles("server/views/layout.tpl"))
		tpl.Execute(w, Da{Title: "Home", Name: "Dave"})
	})
	m.Run()
}
