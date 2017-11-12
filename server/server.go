package server

import (
	"github.com/go-martini/martini"
)

type Server struct{}

type Da struct {
	Title string
	Name  string
}

func (server *Server) Run() {
	m := martini.Classic()
	m.Use(martini.Static("server/public"))
	m.Run()
}
