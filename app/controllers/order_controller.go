package controllers

import "net/http"

func (server *Server) Checkout(w http.ResponseWriter, r *http.Request) {
	if !IsLoggedIn(r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}
