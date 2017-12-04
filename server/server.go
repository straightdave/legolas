package server

import (
	"github.com/codegangsta/martini-contrib/binding"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"gopkg.in/mgo.v2/bson"
	// "net/url"

	A "legolas/common/models/action"
	J "legolas/common/models/job"
	JS "legolas/common/models/jobstate"
	R "legolas/common/models/run"
	T "legolas/common/models/template"
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

	// create new case
	m.Post("/cases", binding.Json(TC.TestCase{}), func(newTC TC.TestCase, ferr binding.Errors, r render.Render) {
		if ferr.Count() > 0 {
			r.JSON(200, Ex{"error": "binding case post data failed"})
			return
		}

		mongo := S.AskForMongo()
		defer mongo.Close()
		TC.SetCol(mongo)

		if err := newTC.Save(); err != nil {
			r.JSON(200, Ex{"error": err.Error()})
		} else {
			r.JSON(200, newTC)
		}
	})

	// update case
	m.Put("/case/:cid", binding.Json(TC.TestCase{}), func(p martini.Params, newTC TC.TestCase, ferr binding.Errors, r render.Render) {
		if ferr.Count() > 0 {
			r.JSON(200, Ex{"error": "binding case post data failed"})
			return
		}

		cid := bson.ObjectIdHex(p["cid"])
		if !cid.Valid() {
			r.JSON(200, Ex{"error": "invalid case id"})
			return
		}

		if cid != newTC.Id {
			r.JSON(200, Ex{"error": "case id mismatches"})
			return
		}

		mongo := S.AskForMongo()
		defer mongo.Close()
		TC.SetCol(mongo)

		if err := newTC.Save(); err != nil {
			r.JSON(200, Ex{"error": err.Error()})
		} else {
			r.JSON(200, newTC)
		}
	})

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

	// add an action to a case
	m.Post("/actions", binding.Json(A.Action{}), func(a A.Action, ferr binding.Errors, r render.Render) {
		if ferr.Count() > 0 {
			r.JSON(200, Ex{"error": "binding action post data failed"})
			return
		}

		mongo := S.AskForMongo()
		defer mongo.Close()
		A.SetCol(mongo)

		if err := a.Save(); err != nil {
			r.JSON(200, Ex{"error": err.Error()})
		} else {
			r.JSON(200, a)
		}
	})

	// update an action
	m.Put("/action/:aid", binding.Json(A.Action{}), func(newAction A.Action, p martini.Params, ferr binding.Errors, r render.Render) {
		if ferr.Count() > 0 {
			r.JSON(200, Ex{"error": "binding action post data failed"})
			return
		}

		aid := bson.ObjectIdHex(p["aid"])
		if !aid.Valid() {
			r.JSON(200, Ex{"error": "invalid action id"})
			return
		}

		if aid != newAction.Id {
			r.JSON(200, Ex{"error": "action id mismatches"})
			return
		}

		mongo := S.AskForMongo()
		defer mongo.Close()
		A.SetCol(mongo)

		if err := newAction.Save(); err != nil {
			r.JSON(200, Ex{"error": err.Error()})
		} else {
			r.JSON(200, newAction)
		}
	})

	m.Delete("/action/:aid", func(p martini.Params, r render.Render) {
		aid := bson.ObjectIdHex(p["aid"])
		if !aid.Valid() {
			r.JSON(200, Ex{"error": "invalid action id"})
			return
		}

		mongo := S.AskForMongo()
		defer mongo.Close()
		A.SetCol(mongo)

		act, err := A.GetOneById(aid)
		if err != nil {
			r.JSON(200, Ex{"error": "failed to get action"})
		}

		if err := act.Delete(); err != nil {
			r.JSON(200, Ex{"error": err.Error()})
		} else {
			r.JSON(200, Ex{"ok": "marked as removed"})
		}
	})

	m.Get("/templates", func(r render.Render) {
		mongo := S.AskForMongo()
		defer mongo.Close()
		T.SetCol(mongo)

		templates, err := T.GetAll2(25)
		if err != nil {
			r.JSON(200, Ex{"error": err.Error()})
		} else {
			r.JSON(200, templates)
		}
	})

	// get a template inro
	m.Get("/template/:tid", func(p martini.Params, r render.Render) {
		mongo := S.AskForMongo()
		defer mongo.Close()
		T.SetCol(mongo)

		tid := bson.ObjectIdHex(p["tid"])
		if !tid.Valid() {
			r.JSON(200, Ex{"error": "invalid template id"})
			return
		}

		template, err := T.GetOneById(tid)
		if err != nil {
			r.JSON(200, Ex{"error": err.Error()})
		} else {
			r.JSON(200, template)
		}
	})

	// create a template
	m.Post("/templates", binding.Json(T.Template{}), func(newT T.Template, ferr binding.Errors, r render.Render) {
		if ferr.Count() > 0 {
			r.JSON(200, Ex{"error": "binding template post data failed"})
			return
		}

		mongo := S.AskForMongo()
		defer mongo.Close()
		T.SetCol(mongo)

		if err := newT.Save(); err != nil {
			r.JSON(200, Ex{"error": err.Error()})
		} else {
			r.JSON(200, newT)
		}
	})

	// update template
	m.Put("/template/:tid", binding.Json(T.Template{}), func(p martini.Params, newT T.Template, ferr binding.Errors, r render.Render) {
		if ferr.Count() > 0 {
			r.JSON(200, Ex{"error": "binding template post data failed"})
			return
		}

		tid := bson.ObjectIdHex(p["tid"])
		if !tid.Valid() {
			r.JSON(200, Ex{"error": "invalid template id"})
			return
		}

		if tid != newT.Id {
			r.JSON(200, Ex{"error": "template id mismatches"})
			return
		}

		mongo := S.AskForMongo()
		defer mongo.Close()
		T.SetCol(mongo)

		if err := newT.Save(); err != nil {
			r.JSON(200, Ex{"error": err.Error()})
		} else {
			r.JSON(200, newT)
		}
	})

	m.Run()
}
