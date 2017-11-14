package server

import (
	"github.com/codegangsta/martini-contrib/binding"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"legolas/server/models"
)

type Server struct{}

func (server *Server) Run() {
	m := martini.Classic()
	m.Use(martini.Static("server/public"))
	m.Use(render.Renderer())

	m.Get("/cases", func(r render.Render) {
		cases, err := models.FindAllCases()
		if err != nil {
			r.JSON(200, map[string]interface{}{"error": err.Error()})
		} else {
			r.JSON(200, cases)
		}
	})

	m.Post("/cases", binding.Json(models.Case{}), func(c models.Case, ferr binding.Errors, r render.Render) {
		if ferr.Count() > 0 {
			r.JSON(200, map[string]interface{}{"error": "binding case post data failed"})
		} else {
			c0 := models.NewCase(c.Path, c.Name, c.Desc)
			if err := c0.Save(); err != nil {
				r.JSON(200, map[string]interface{}{"error": err.Error()})
			} else {
				r.JSON(200, *c0)
			}
		}
	})

	m.Get("/case/:path/:name", func(p martini.Params, r render.Render) {
		cpath := p["path"]
		cname := p["name"]
		c, err := models.FindCase(cpath, cname)
		if err != nil {
			r.JSON(200, map[string]interface{}{"error": err.Error()})
		} else {
			r.JSON(200, *c)
		}
	})

	m.Get("/case/:path/:name/actions", func(p martini.Params, r render.Render) {
		cpath := p["path"]
		cname := p["name"]
		actions, err := models.FindActions(cname, cpath)
		if err != nil {
			r.JSON(200, map[string]interface{}{"error": err.Error()})
		} else {
			r.JSON(200, actions)
		}
	})

	m.Get("/actions", func(r render.Render) {
		actions, err := models.FindAllActions()
		if err != nil {
			r.JSON(200, map[string]interface{}{"error": err.Error()})
		} else {
			r.JSON(200, actions)
		}
	})

	m.Run()
}
