package postgres

import (
	"testing"

	"github.com/sfu-teamproject/smartbuy/backend/models"
	"github.com/stretchr/testify/assert"
)

func TestCartItems(t *testing.T) {
	db, err := NewPostgresDB(true)
	assert.NoError(t, err, "postgres db creating failed", err.Error())
	cartItem := models.CartItem{
		CartID:       1,
		SmartphoneID: 1,
		Quantity:     1,
	}
	t.Run("add cart item", func(t *testing.T) {
		newCartItem, err := db.AddToCart(cartItem)
		assert.NoError(t, err, "adding cart item failed", err.Error())
		assert.NotEmpty(t, newCartItem.ID, "cart item id is 0")
		cartItem.ID = newCartItem.ID
		assert.Equal(t, cartItem, newCartItem, "cartItem is different")
	})
	t.Run("get cartItems", func(t *testing.T) {
		cartItems, err := db.GetCartItems(1)
		assert.NoError(t, err, "getting cart items failed", err.Error())
		assert.Equal(t, 1, len(cartItems), "length of cart items is not 1")
		assert.Equal(t, cartItem, cartItems[0], "cart item is not the same")
	})
	t.Run("set cartItem quantity", func(t *testing.T) {
		cartItem.Quantity = 3
		newCartItem, err := db.SetQuantity(cartItem)
		assert.NoError(t, err, "setting cart item quantity failed", err.Error())
		assert.Equal(t, 3, newCartItem.Quantity, "quantity of cart item is not 3")
	})
	t.Run("delete cart item", func(t *testing.T) {
		deletedCartItem, err := db.DeleteFromCart(1, 1)
		assert.NoError(t, err, "deleting cart item failed", err.Error())
		assert.Equal(t, cartItem, deletedCartItem, "cart item is different")
		cartItems, err := db.GetCartItems(1)
		assert.NoError(t, err, "getting cart items failed", err.Error())
		assert.Empty(t, cartItems, "cart is not empty")
	})
}
