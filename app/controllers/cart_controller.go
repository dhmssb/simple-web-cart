package controllers

import (
	"fmt"
	"log"
	"net/http"
	"simpleWebCart/app/models"
	"strconv"

	"github.com/google/uuid"
	"github.com/unrolled/render"
	"gorm.io/gorm"
)

func GetShoppingCartID(w http.ResponseWriter, r *http.Request) string {
	session, _ := store.Get(r, sessionShoppingCart)
	if session.Values["cart-id"] == nil {
		session.Values["cart-id"] = uuid.New().String()
		session.Save(r, w)
	}

	return fmt.Sprintf("%v", session.Values["cart-id"])
}

func GetShoppingCart(db *gorm.DB, cartID string) (*models.Cart, error) {
	var cart models.Cart

	existCart, err := cart.GetCart(db, cartID)
	if err != nil {
		existCart, _ = cart.CreateCart(db, cartID)
	}

	_, _ = existCart.CalculateCart(db, cartID)

	return existCart, nil

}

func (server *Server) GetCart(w http.ResponseWriter, r *http.Request) {
	render := render.New(render.Options{
		Layout:     "layout",
		Extensions: []string{".html", ".tmpl"},
	})

	var (
		cart *models.Cart
		err  error
	)
	cartID := GetShoppingCartID(w, r)
	cart, _ = GetShoppingCart(server.DB, cartID)
	items, _ := cart.GetItems(server.DB, cartID)

	if err != nil {
		log.Fatal(err)
	}

	_ = render.HTML(w, http.StatusOK, "cart", map[string]interface{}{
		"cart":  cart,
		"items": items,
	})
}

func (server *Server) AddItemToCart(w http.ResponseWriter, r *http.Request) {
	productID := r.FormValue("product_id")
	qty, _ := strconv.Atoi(r.FormValue("qty"))

	productModel := models.Product{}
	product, err := productModel.FindByID(server.DB, productID)
	if err != nil {
		http.Redirect(w, r, "/products/"+product.Slug, http.StatusSeeOther)
		return
	}

	if qty > product.Stock {
		// SetFlash(w, r, "error", "Stok tidak mencukupi")
		http.Redirect(w, r, "/products/"+product.Slug, http.StatusSeeOther)
		return
	}

	var cart *models.Cart

	cartID := GetShoppingCartID(w, r)
	cart, _ = GetShoppingCart(server.DB, cartID)
	_, err = cart.AddItem(server.DB, models.CartItem{
		ProductID: productID,
		Qty:       qty,
	})
	if err != nil {
		http.Redirect(w, r, "/products/"+product.Slug, http.StatusSeeOther)
	}

	// SetFlash(w, r, "success", "Item berhasil ditambahkan")
	http.Redirect(w, r, "/carts", http.StatusSeeOther)
}
