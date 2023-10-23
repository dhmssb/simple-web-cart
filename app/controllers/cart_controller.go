package controllers

import (
	"fmt"
	"net/http"
	"simpleWebCart/app/models"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
	"gorm.io/gorm"
)

func (server *Server) GetCart(w http.ResponseWriter, r *http.Request) {

	render := render.New(render.Options{
		Layout: "layout",
	})
	var cart *models.Cart

	cartID := GetShoppingCartID(w, r)
	cart, _ = GetShoppingCart(server.DB, cartID)
	items, _ := cart.GetItems(server.DB, cart.ID)

	_ = render.HTML(w, http.StatusOK, "cart", map[string]interface{}{
		"cart":    cart,
		"items":   items,
		"success": GetFlash(w, r, "success"),
		"error":   GetFlash(w, r, "error"),
	})

}

func GetShoppingCartID(w http.ResponseWriter, r *http.Request) string {
	session, _ := store.Get(r, sessionShoppingCart)
	if session.Values["cart-id"] == nil {
		session.Values["cart-id"] = uuid.New().String()
		_ = session.Save(r, w)
	}

	return fmt.Sprintf("%v", session.Values["cart-id"])
}

func GetShoppingCart(db *gorm.DB, cartID string) (*models.Cart, error) {
	var cart models.Cart

	existCart, err := cart.GetCart(db, cartID)
	if err != nil {
		existCart, _ = cart.CreateCart(db, cartID)
	}

	existCart.CalculateCart(db, cartID)
	updatedCart, _ := cart.GetCart(db, cartID)

	return updatedCart, nil
}

func (server *Server) AddItemToCart(w http.ResponseWriter, r *http.Request) {
	producID := r.FormValue("product_id")
	qty, _ := strconv.Atoi(r.FormValue("qty"))

	productModel := models.Product{}
	product, err := productModel.FindByID(server.DB, producID)
	if err != nil {
		http.Redirect(w, r, "/products/"+product.Slug, http.StatusSeeOther)
		return
	}

	if qty > product.Stock {
		SetFlash(w, r, "error", "Stock tidak mencukupi")
		http.Redirect(w, r, "/products/"+product.Slug, http.StatusSeeOther)
		return
	}

	var cart *models.Cart

	cartID := GetShoppingCartID(w, r)
	cart, _ = GetShoppingCart(server.DB, cartID)
	_, err = cart.AddItem(server.DB, models.CartItem{
		ProductID: producID,
		Qty:       qty,
	})
	if err != nil {
		SetFlash(w, r, "error", "error")
		http.Redirect(w, r, "/products/"+product.Slug, http.StatusSeeOther)

	}
	SetFlash(w, r, "success", "Item berhasil ditambahkan")
	http.Redirect(w, r, "/carts", http.StatusSeeOther)
}

func (server *Server) UpdateCart(w http.ResponseWriter, r *http.Request) {
	cartID := GetShoppingCartID(w, r)
	cart, _ := GetShoppingCart(server.DB, cartID)

	for _, item := range cart.CartItems {
		qty, _ := strconv.Atoi(r.FormValue(item.ID))

		_, err := cart.UpdateItemQty(server.DB, item.ID, qty)
		if err != nil {
			http.Redirect(w, r, "/carts", http.StatusSeeOther)
		}
	}
	http.Redirect(w, r, "/carts", http.StatusSeeOther)
}

func (server *Server) RemoveItemByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	if vars["id"] == "" {
		http.Redirect(w, r, "/carts", http.StatusSeeOther)
	}

	cartID := GetShoppingCartID(w, r)
	cart, _ := GetShoppingCart(server.DB, cartID)

	err := cart.RemoveItemByID(server.DB, vars["id"])
	if err != nil {
		http.Redirect(w, r, "/carts", http.StatusSeeOther)
	}
	http.Redirect(w, r, "/carts", http.StatusSeeOther)

}
