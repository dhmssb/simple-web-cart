package controllers

import (
	"net/http"

	"github.com/unrolled/render"
)

func (server *Server) Checkout(w http.ResponseWriter, r *http.Request) {
	render := render.New(render.Options{
		Layout:     "layout",
		Extensions: []string{".html", ".tmpl"},
	})

	if !IsLoggedIn(r) {
		SetFlash(w, r, "error", "anda perlu login!")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}

	_ = render.HTML(w, http.StatusOK, "checkout", map[string]interface{}{
		"user": server.CurrentUser(w, r),
	})
}
