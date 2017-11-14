package server

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"legolas/server/models"
)

type Server struct{}

func (server *Server) Run() {
	m := martini.Classic()
	m.Use(martini.Static("server/public"))
	m.Use(render.Renderer())

	// geg all cases
	m.Get("/cases", func(r render.Render) {
		cases, err := models.FindAllCases()
		if err != nil {
			r.JSON(500, map[string]interface{}{"error": err})
		} else {
			r.JSON(200, cases)
		}
	})

	// get actions in one case
	m.Get("/case/:path/:name/actions", func(p martini.Params, r render.Render) {
		cpath := p["path"]
		cname := p["name"]
		actions, err := models.FindActions(cname, cpath)
		if err != nil {
			r.JSON(500, map[string]interface{}{"error": err})
		} else {
			r.JSON(200, actions)
		}
	})

	m.Get("/actions", func(r render.Render) {
		actions, err := models.FindAllActions()
		if err != nil {
			r.JSON(500, map[string]interface{}{"error": err})
		} else {
			r.JSON(200, actions)
		}
	})

	m.Run()
}
