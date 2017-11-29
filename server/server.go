package server

import (
	// "github.com/codegangsta/martini-contrib/binding"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	// "net/url"

	A "legolas/common/models/action"
	J "legolas/common/models/job"
	JS "legolas/common/models/jobstate"
	R "legolas/common/models/run"
	TC "legolas/common/models/testcase"
	S "legolas/common/storage"
)

type Server struct{}
type Ex map[string]interface{}

func (server *Server) Run() {
	m := martini.Classic()
	m.Use(martini.Static("server/public"))
	m.Use(render.Renderer())

	J.SetRedisPool(S.GetRedisPool())

	// m.Get("/cases/f/:word", func(p martini.Params, r render.Render) {
	// 	word, err := url.QueryUnescape(p["word"])
	// 	if err != nil {
	// 		r.JSON(200, Ex{"error": err.Error()})
	// 		return
	// 	}

	// 	cases, err := models.FilterCases(word)
	// 	if err != nil {
	// 		r.JSON(200, Ex{"error": err.Error()})
	// 	} else {
	// 		r.JSON(200, cases)
	// 	}
	// })

	m.Get("/cases", func(r render.Render) {
		mongo := S.AskForMongo()
		defer mongo.Close()
		TC.SetCol(mongo)

		cases, err := TC.GetAllInTimeOrder(50)
		if err != nil {
			r.JSON(200, Ex{"error": err.Error()})
		} else {
			r.JSON(200, cases)
		}
	})

	// // create case
	// m.Post("/cases", binding.Json(models.Case{}), func(c models.Case, ferr binding.Errors, r render.Render) {
	// 	if ferr.Count() > 0 {
	// 		r.JSON(200, Ex{"error": "binding case post data failed"})
	// 	} else {
	// 		c0 := models.NewCase(c.Path, c.Name, c.Desc)
	// 		if err := c0.Save(); err != nil {
	// 			r.JSON(200, Ex{"error": err.Error()})
	// 		} else {
	// 			r.JSON(200, *c0)
	// 		}
	// 	}
	// })

	// // update case
	// m.Put("/case/:path/:name", binding.Json(models.Case{}), func(p martini.Params, c models.Case, ferr binding.Errors, r render.Render) {
	// 	if ferr.Count() > 0 {
	// 		r.JSON(200, Ex{"error": "binding case post data failed"})
	// 		return
	// 	}

	// 	oldCase, err := models.FindCase(p["path"], p["name"])
	// 	if err != nil {
	// 		r.JSON(200, Ex{"error": "cannot find such case"})
	// 		return
	// 	}

	// 	if err = oldCase.UpdateTo(&c); err != nil {
	// 		r.JSON(200, Ex{"error": err.Error()})
	// 	} else {
	// 		r.JSON(200, c)
	// 	}
	// })

	// // get a case
	// m.Get("/case/:path/:name", func(p martini.Params, r render.Render) {
	// 	c, err := models.FindCase(p["path"], p["name"])
	// 	if err != nil {
	// 		r.JSON(200, Ex{"error": err.Error()})
	// 	} else {
	// 		r.JSON(200, *c)
	// 	}
	// })

	// run a case
	m.Post("/case/:cid/runs", func(p martini.Params, r render.Render) {
		mongo := S.AskForMongo()
		defer mongo.Close()
		R.SetCol(mongo)

		run, err := R.InvokeByCaseIdStr(p["cid"])
		if err != nil {
			r.JSON(200, Ex{"error": err.Error()})
			return
		}
		r.JSON(200, run)
	})

	// // delete a case
	// m.Delete("/case/:path/:name", func(p martini.Params, r render.Render) {
	// 	err := models.DeleteCase(p["path"], p["name"])
	// 	if err != nil {
	// 		r.JSON(200, Ex{"error": err.Error()})
	// 	} else {
	// 		r.JSON(200, Ex{"ok": "deleted"})
	// 	}
	// })

	// get all actions of a case
	m.Get("/case/:cid/actions", func(p martini.Params, r render.Render) {
		mongo := S.AskForMongo()
		defer mongo.Close()
		A.SetCol(mongo)

		actions, err := A.GetAllByCaseIdStr(p["cid"])
		if err != nil {
			r.JSON(200, Ex{"error": err.Error()})
		} else {
			r.JSON(200, actions)
		}
	})

	// get all runs of a case
	m.Get("/case/:cid/runs", func(p martini.Params, r render.Render) {
		mongo := S.AskForMongo()
		defer mongo.Close()
		R.SetCol(mongo)

		runs, err := R.GetAllByCaseIdStr(p["cid"])
		if err != nil {
			r.JSON(200, Ex{"error": err.Error()})
		} else {
			r.JSON(200, runs)
		}
	})

	// get all job states of one run
	m.Get("/run/:rid/jobstates", func(p martini.Params, r render.Render) {
		mongo := S.AskForMongo()
		defer mongo.Close()
		JS.SetCol(mongo)

		jss, err := JS.GetAllByRunIdStr(p["rid"])
		if err != nil {
			r.JSON(200, Ex{"error": err.Error()})
		} else {
			r.JSON(200, jss)
		}
	})

	// // get one action of a case
	// m.Get("/case/:cpath/:cname/:name", func(p martini.Params, r render.Render) {
	// 	action, err := models.FindAction(p["cname"], p["cpath"], p["name"])
	// 	if err != nil {
	// 		r.JSON(200, Ex{"error": err.Error()})
	// 	} else {
	// 		r.JSON(200, action)
	// 	}
	// })

	// // add a new action to a case
	// m.Post("/actions", binding.Json(models.Action{}), func(a models.Action, ferr binding.Errors, r render.Render) {
	// 	if ferr.Count() > 0 {
	// 		r.JSON(200, Ex{"error": "binding action post data failed"})
	// 		return
	// 	}

	// 	if err := a.Save(); err != nil {
	// 		r.JSON(200, Ex{"error": err.Error()})
	// 	} else {
	// 		r.JSON(200, a)
	// 	}
	// })

	// // update an action
	// m.Put("/case/:cpath/:cname/:name", binding.Json(models.Action{}), func(newAction models.Action, p martini.Params, ferr binding.Errors, r render.Render) {
	// 	if ferr.Count() > 0 {
	// 		r.JSON(200, Ex{"error": "binding action post data failed"})
	// 		return
	// 	}

	// 	a0 := models.NewAction(p["cpath"], p["cname"], p["name"])
	// 	if err := a0.UpdateTo(&newAction); err != nil {
	// 		r.JSON(200, Ex{"error": err.Error()})
	// 	} else {
	// 		r.JSON(200, newAction)
	// 	}
	// })

	// m.Delete("/case/:cpath/:cname/:name", func(p martini.Params, r render.Render) {
	// 	err := models.DeleteAction(p["cpath"], p["cname"], p["name"])
	// 	if err != nil {
	// 		r.JSON(200, Ex{"error": err.Error()})
	// 	} else {
	// 		r.JSON(200, Ex{"ok": "deleted"})
	// 	}
	// })

	// m.Get("/actions", func(r render.Render) {
	// 	actions, err := models.FindAllActions()
	// 	if err != nil {
	// 		r.JSON(200, Ex{"error": err.Error()})
	// 	} else {
	// 		r.JSON(200, actions)
	// 	}
	// })

	m.Run()
}
