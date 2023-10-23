package controllers

import (
	"net/http"

	"github.com/unrolled/render"
)

func (server *Server) Home(w http.ResponseWriter, r *http.Request) {

	render := render.New(render.Options{Layout: "layout"})

	user := server.CurrentUser(w, r)
	_ = render.HTML(w, http.StatusOK, "home", map[string]interface{}{
		"user": user,
	})
}
