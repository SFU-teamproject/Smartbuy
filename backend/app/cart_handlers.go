package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/sfu-teamproject/smartbuy/backend/models"
)

func (app *App) GetCart(w http.ResponseWriter, r *http.Request) {
	ID, err := app.ExtractPathValue(r, "cart_id")
	if err != nil {
		app.ErrorJSON(w, r, err)
		return
	}
	cart, err := app.DB.GetCart(ID)
	if err != nil {
		app.ErrorJSON(w, r, fmt.Errorf("error getting cart %d: %w", ID, err))
		return
	}
	userID, role, err := app.GetClaims(r)
	if err != nil {
		app.ErrorJSON(w, r, fmt.Errorf("%w: error extracting claims: %w", errUnauthorized, err))
		return
	}
	if cart.UserID != userID && role != models.RoleAdmin {
		app.ErrorJSON(w, r, errForbidden)
		return
	}
	cartItems, err := app.DB.GetCartItems(cart.ID)
	if err != nil {
		app.ErrorJSON(w, r, fmt.Errorf("error getting items for cart %d: %w", cart.ID, err))
		return
	}
	cart.Items = cartItems
	app.Encode(w, r, cart)
}

func (app *App) GetCarts(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr != "" {
		app.GetCartByUserID(w, r)
		return
	}
	userID, role, err := app.GetClaims(r)
	if err != nil {
		app.ErrorJSON(w, r, fmt.Errorf("%w: error extracting claims: %w", errUnauthorized, err))
		return
	}
	if role != models.RoleAdmin {
		app.ErrorJSON(w, r, fmt.Errorf("%w: user %d (role %s) does not have required role",
			errForbidden, userID, role))
		return
	}
	carts, err := app.DB.GetCarts()
	if err != nil {
		app.ErrorJSON(w, r, fmt.Errorf("error getting carts: %w", err))
		return
	}
	app.Encode(w, r, carts)
}

func (app *App) GetCartByUserID(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("user_id")
	pUserID, err := strconv.Atoi(userIDStr)
	if err != nil {
		app.ErrorJSON(w, r, fmt.Errorf("%w: incorrect user id(%s): %w", errBadRequest, userIDStr, err))
		return
	}
	userID, role, err := app.GetClaims(r)
	if err != nil {
		app.ErrorJSON(w, r, fmt.Errorf("%w: error extracting claims: %w", errUnauthorized, err))
		return
	}
	if pUserID != userID && role != models.RoleAdmin {
		app.ErrorJSON(w, r, errForbidden)
		return
	}
	cart, err := app.DB.GetCartByUserID(pUserID)
	if err != nil {
		app.ErrorJSON(w, r, fmt.Errorf("error getting cart by user id(%d): %w", userID, err))
		return
	}
	cartItems, err := app.DB.GetCartItems(cart.ID)
	if err != nil {
		app.ErrorJSON(w, r, fmt.Errorf("error getting items for cart %d: %w", cart.ID, err))
		return
	}
	cart.Items = cartItems
	app.Encode(w, r, cart)
}

func (app *App) GetCartItems(w http.ResponseWriter, r *http.Request) {
	cartID, err := app.ExtractPathValue(r, "cart_id")
	if err != nil {
		app.ErrorJSON(w, r, err)
		return
	}
	cart, err := app.DB.GetCart(cartID)
	if err != nil {
		app.ErrorJSON(w, r, fmt.Errorf("error getting cart %d: %w", cartID, err))
		return
	}
	userID, role, err := app.GetClaims(r)
	if err != nil {
		app.ErrorJSON(w, r, fmt.Errorf("%w: error extracting claims: %w", errUnauthorized, err))
		return
	}
	if userID != cart.UserID && role != models.RoleAdmin {
		app.ErrorJSON(w, r, errForbidden)
		return
	}
	cartItems, err := app.DB.GetCartItems(cartID)
	if err != nil {
		app.ErrorJSON(w, r, fmt.Errorf("error getting cartItems of cart id(%d): %w", cartID, err))
		return
	}
	app.Encode(w, r, cartItems)
}

func (app *App) AddToCart(w http.ResponseWriter, r *http.Request) {
	cartID, err := app.ExtractPathValue(r, "cart_id")
	if err != nil {
		app.ErrorJSON(w, r, err)
		return
	}
	cart, err := app.DB.GetCart(cartID)
	if err != nil {
		app.ErrorJSON(w, r, fmt.Errorf("error getting cart %d: %w", cartID, err))
		return
	}
	userID, role, err := app.GetClaims(r)
	if err != nil {
		app.ErrorJSON(w, r, fmt.Errorf("%w: error extracting claims: %w", errUnauthorized, err))
		return
	}
	if userID != cart.UserID && role != models.RoleAdmin {
		app.ErrorJSON(w, r, errForbidden)
		return
	}
	var cartItem models.CartItem
	err = json.NewDecoder(r.Body).Decode(&cartItem)
	if err != nil {
		app.ErrorJSON(w, r, fmt.Errorf("%w: error decoding cartItem: %w", errBadRequest, err))
		return
	}
	cartItem.CartID = cartID
	addedCartItem, err := app.DB.AddToCart(cartItem)
	if err != nil {
		app.ErrorJSON(w, r, fmt.Errorf("error creating cartItem: %w", err))
		return
	}
	w.WriteHeader(http.StatusCreated)
	app.Encode(w, r, addedCartItem)
}

func (app *App) SetQuantity(w http.ResponseWriter, r *http.Request) {
	itemID, err := app.ExtractPathValue(r, "item_id")
	if err != nil {
		app.ErrorJSON(w, r, err)
		return
	}
	cartID, err := app.ExtractPathValue(r, "cart_id")
	if err != nil {
		app.ErrorJSON(w, r, err)
		return
	}
	cart, err := app.DB.GetCart(cartID)
	if err != nil {
		app.ErrorJSON(w, r, fmt.Errorf("error getting cart %d: %w", cartID, err))
		return
	}
	userID, role, err := app.GetClaims(r)
	if err != nil {
		app.ErrorJSON(w, r, fmt.Errorf("%w: error extracting claims: %w", errUnauthorized, err))
		return
	}
	if userID != cart.UserID && role != models.RoleAdmin {
		app.ErrorJSON(w, r, errForbidden)
		return
	}
	var cartItem models.CartItem
	err = json.NewDecoder(r.Body).Decode(&cartItem)
	if err != nil {
		app.ErrorJSON(w, r, fmt.Errorf("%w: error decoding quantity: %w", errBadRequest, err))
		return
	}
	cartItem.ID = itemID
	cartItem.CartID = cartID
	updatedCartItem, err := app.DB.SetQuantity(cartItem)
	if err != nil {
		app.ErrorJSON(w, r, fmt.Errorf("error updating cartItem: %w", err))
		return
	}
	app.Encode(w, r, updatedCartItem)
}

func (app *App) DeleteFromCart(w http.ResponseWriter, r *http.Request) {
	itemID, err := app.ExtractPathValue(r, "item_id")
	if err != nil {
		app.ErrorJSON(w, r, err)
		return
	}
	cartID, err := app.ExtractPathValue(r, "cart_id")
	if err != nil {
		app.ErrorJSON(w, r, err)
		return
	}
	cart, err := app.DB.GetCart(cartID)
	if err != nil {
		app.ErrorJSON(w, r, fmt.Errorf("error getting cart %d: %w", cartID, err))
		return
	}
	userID, role, err := app.GetClaims(r)
	if err != nil {
		app.ErrorJSON(w, r, fmt.Errorf("%w: error extracting claims: %w", errUnauthorized, err))
		return
	}
	if userID != cart.UserID && role != models.RoleAdmin {
		app.ErrorJSON(w, r, errForbidden)
		return
	}
	cartItem, err := app.DB.DeleteFromCart(cartID, itemID)
	if err != nil {
		app.ErrorJSON(w, r, fmt.Errorf("error deleting cartItem %d: %w", itemID, err))
		return
	}
	app.Encode(w, r, cartItem)
}
