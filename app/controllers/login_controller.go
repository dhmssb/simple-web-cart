package controllers

import (
	"net/http"
	"simpleWebCart/app/models"

	"github.com/unrolled/render"
)

func (server *Server) Login(w http.ResponseWriter, r *http.Request) {
	render := render.New(render.Options{
		Layout:     "layout",
		Extensions: []string{".html", ".tmpl"},
	})

	_ = render.HTML(w, http.StatusOK, "login", map[string]interface{}{
		"title": "Login",
		"error": GetFlash(w, r, "error"),
	})
}

func (server *Server) DoLogin(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	userModel := models.User{}
	user, err := userModel.FindByEmail(server.DB, username)
	if err != nil {
		SetFlash(w, r, "error", "email or password invalid")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if !ComparePassword(password, user.Password) {
		SetFlash(w, r, "error", "email or password invalid")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	session, _ := store.Get(r, sessionUser)
	session.Values["id"] = user.ID
	session.Save(r, w)

	http.Redirect(w, r, "/products", http.StatusSeeOther)

}

func (server *Server) Logout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, sessionUser)

	session.Values["id"] = nil
	session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
