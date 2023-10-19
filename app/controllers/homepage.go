package controllers

import (
	"net/http"

	"github.com/unrolled/render"
)

func (server *Server) Home(w http.ResponseWriter, r *http.Request) {

	render := render.New(render.Options{Layout: "layout"})

	_ = render.HTML(w, http.StatusOK, "home", map[string]interface{}{
		"Title": "Home Page",
		"Body":  "Home Description",
	})
}
