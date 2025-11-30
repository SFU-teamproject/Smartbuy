package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/sfu-teamproject/smartbuy/backend/apperrors"
	"github.com/sfu-teamproject/smartbuy/backend/models"
)

// Gets a cart by cart_id
// @Summary      Get a Cart
// @Description  Gets a Cart with a certain cart_id
// @Tags         cart
// @Security     BearerAuth
// @Produce      json
// @Param        cart_id path int true "Cart ID"
// @Success      200  {object}  models.Cart
// @Failure      401  {object}  apperrors.ErrorResponse "Unauthorized"
// @Failure      403  {object}  apperrors.ErrorResponse "Forbidden"
// @Router       /carts/{cart_id} [get]
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
		app.ErrorJSON(w, r, fmt.Errorf("%w: error extracting claims: %w", apperrors.ErrUnauthorized, err))
		return
	}
	if cart.UserID != userID && role != models.RoleAdmin {
		app.ErrorJSON(w, r, apperrors.ErrForbidden)
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

// @Summary      Get all Carts
// @Description  Gets all Carts. If user id provided in a query, then return only cart belonging to that user
// @Tags         cart
// @Param        user_id query int false "User ID"
// @Security     BearerAuth
// @Produce      json
// @Success      200  {array}  models.Cart
// @Failure      401  {object}  apperrors.ErrorResponse "Unauthorized"
// @Failure      403  {object}  apperrors.ErrorResponse "Forbidden"
// @Router       /carts [get]
func (app *App) GetCarts(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr != "" {
		app.GetCartByUserID(w, r)
		return
	}
	userID, role, err := app.GetClaims(r)
	if err != nil {
		app.ErrorJSON(w, r, fmt.Errorf("%w: error extracting claims: %w", apperrors.ErrUnauthorized, err))
		return
	}
	if role != models.RoleAdmin {
		app.ErrorJSON(w, r, fmt.Errorf("%w: user %d (role %s) does not have required role",
			apperrors.ErrForbidden, userID, role))
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
		app.ErrorJSON(w, r, fmt.Errorf("%w: incorrect user id(%s): %w", apperrors.ErrBadRequest, userIDStr, err))
		return
	}
	userID, role, err := app.GetClaims(r)
	if err != nil {
		app.ErrorJSON(w, r, fmt.Errorf("%w: error extracting claims: %w", apperrors.ErrUnauthorized, err))
		return
	}
	if pUserID != userID && role != models.RoleAdmin {
		app.ErrorJSON(w, r, apperrors.ErrForbidden)
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

// @Summary      Get all Cart items
// @Description  Gets all Cart items from a certain cart
// @Tags         cart
// @Security     BearerAuth
// @Produce      json
// @Param        cart_id path int true "Cart ID"
// @Success      200  {array}  models.CartItem
// @Failure      401  {object}  apperrors.ErrorResponse "Unauthorized"
// @Failure      403  {object}  apperrors.ErrorResponse "Forbidden"
// @Router       /carts/{cart_id}/items [get]
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
		app.ErrorJSON(w, r, fmt.Errorf("%w: error extracting claims: %w", apperrors.ErrUnauthorized, err))
		return
	}
	if userID != cart.UserID && role != models.RoleAdmin {
		app.ErrorJSON(w, r, apperrors.ErrForbidden)
		return
	}
	cartItems, err := app.DB.GetCartItems(cartID)
	if err != nil {
		app.ErrorJSON(w, r, fmt.Errorf("error getting cartItems of cart id(%d): %w", cartID, err))
		return
	}
	app.Encode(w, r, cartItems)
}

// AddToCart adds an item to cart
// @Summary      Add Item to Cart
// @Description  Adds a smartphone item to the user's cart
// @Tags         cart
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        cart_id path int true "Cart ID"
// @Param        cart_item body models.CartItemRequest true "Item to add"
// @Success      201  {object}  models.CartItem
// @Failure      401  {object}  apperrors.ErrorResponse "Unauthorized"
// @Failure      403  {object}  apperrors.ErrorResponse "Forbidden"
// @Router       /carts/{cart_id}/items [post]
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
		app.ErrorJSON(w, r, fmt.Errorf("%w: error extracting claims: %w", apperrors.ErrUnauthorized, err))
		return
	}
	if userID != cart.UserID && role != models.RoleAdmin {
		app.ErrorJSON(w, r, apperrors.ErrForbidden)
		return
	}
	var cartItemreq models.CartItemRequest
	err = json.NewDecoder(r.Body).Decode(&cartItemreq)
	if err != nil {
		app.ErrorJSON(w, r, fmt.Errorf("%w: error decoding cartItem: %w", apperrors.ErrBadRequest, err))
		return
	}
	cartItem := models.CartItem{SmartphoneID: cartItemreq.SmartphoneID, Quantity: cartItemreq.Quantity}
	cartItem.CartID = cartID
	if cartItem.Quantity < 1 {
		cartItem.Quantity = 1
	}
	addedCartItem, err := app.DB.AddToCart(cartItem)
	if err != nil {
		app.ErrorJSON(w, r, fmt.Errorf("error creating cartItem: %w", err))
		return
	}
	w.WriteHeader(http.StatusCreated)
	app.Encode(w, r, addedCartItem)
}

// @Summary      Sets quantity of an item in a cart
// @Description  Sets quantity of an item in a cart
// @Tags         cart
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        item_id path int true "Item ID"
// @Param        cart_id path int true "Cart ID"
// @Param        quantity body models.SetItemQuantity true "New quantity of the item"
// @Success      200  {object}  models.CartItem
// @Failure      401  {object}  apperrors.ErrorResponse "Unauthorized"
// @Failure      403  {object}  apperrors.ErrorResponse "Forbidden"
// @Router       /carts/{cart_id}/items/{item_id} [patch]
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
		app.ErrorJSON(w, r, fmt.Errorf("%w: error extracting claims: %w", apperrors.ErrUnauthorized, err))
		return
	}
	if userID != cart.UserID && role != models.RoleAdmin {
		app.ErrorJSON(w, r, apperrors.ErrForbidden)
		return
	}
	var quant models.SetItemQuantity
	err = json.NewDecoder(r.Body).Decode(&quant)
	if err != nil {
		app.ErrorJSON(w, r, fmt.Errorf("%w: error decoding quantity: %w", apperrors.ErrBadRequest, err))
		return
	}
	if quant.Quantity < 1 {
		app.ErrorJSON(w, r, fmt.Errorf("%w: qunatity must be a positive integer", apperrors.ErrBadRequest))
		return
	}
	cartItem := models.CartItem{Quantity: quant.Quantity}
	cartItem.ID = itemID
	cartItem.CartID = cartID
	updatedCartItem, err := app.DB.SetQuantity(cartItem)
	if err != nil {
		app.ErrorJSON(w, r, fmt.Errorf("error updating cartItem: %w", err))
		return
	}
	app.Encode(w, r, updatedCartItem)
}

// @Summary      Deletes an item from a cart
// @Description  Deletes an item from a cart
// @Tags         cart
// @Security     BearerAuth
// @Produce      json
// @Param        item_id path int true "Item ID"
// @Param        cart_id path int true "Cart ID"
// @Success      200  {object}  models.CartItem
// @Failure      401  {object}  apperrors.ErrorResponse "Unauthorized"
// @Failure      403  {object}  apperrors.ErrorResponse "Forbidden"
// @Router       /carts/{cart_id}/items/{item_id} [delete]
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
		app.ErrorJSON(w, r, fmt.Errorf("%w: error extracting claims: %w", apperrors.ErrUnauthorized, err))
		return
	}
	if userID != cart.UserID && role != models.RoleAdmin {
		app.ErrorJSON(w, r, apperrors.ErrForbidden)
		return
	}
	cartItem, err := app.DB.DeleteFromCart(cartID, itemID)
	if err != nil {
		app.ErrorJSON(w, r, fmt.Errorf("error deleting cartItem %d: %w", itemID, err))
		return
	}
	app.Encode(w, r, cartItem)
}
