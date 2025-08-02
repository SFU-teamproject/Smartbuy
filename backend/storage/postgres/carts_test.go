package postgres

import (
	"testing"

	"github.com/sfu-teamproject/smartbuy/backend/models"
	"github.com/stretchr/testify/assert"
)

func TestCart(t *testing.T) {
	db, err := NewPostgresDB()
	assert.NoError(t, err, "postgres db creating failed", err.Error())
	t.Run("get cart by ID", func(t *testing.T) {
		cart, err := db.GetCart(1)
		assert.NoError(t, err, "getting cart failed", err.Error())
		assert.NotEqual(t, cart, models.Cart{}, "cart is empty struct")
	})
	t.Run("get cart by ID", func(t *testing.T) {
		cart, err := db.GetCartByUserID(1)
		assert.NoError(t, err, "getting cart by ID failed", err.Error())
		assert.NotEqual(t, cart, models.Cart{}, "cart is empty struct")
	})
	t.Run("get carts", func(t *testing.T) {
		carts, err := db.GetCarts()
		assert.NoError(t, err, "getting carts failed", err.Error())
		assert.NotEmpty(t, carts, "cart slice is empty")
	})
}
