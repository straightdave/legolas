package server

import (
	"net/http"
)

type Server struct {
}

func (server Server) Start() error {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	err := http.ListenAndServe("localhost:8888", nil)
	return err
}
