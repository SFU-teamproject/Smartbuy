package postgres

import (
	"testing"

	"github.com/sfu-teamproject/smartbuy/backend/models"
	"github.com/stretchr/testify/assert"
)

func TestSmartphones(t *testing.T) {
	db, err := NewPostgresDB(true)
	assert.NoError(t, err, "postgres db creating failed", err.Error())
	t.Run("get smartphones", func(t *testing.T) {
		smartphones, err := db.GetSmartphones()
		assert.NoError(t, err, "getting smartphones failed", err.Error())
		assert.NotEmpty(t, smartphones, "smartphone slice is empty")
	})
	t.Run("get smartphone", func(t *testing.T) {
		smartphone, err := db.GetSmartphone(1)
		assert.NoError(t, err, "getting smartphone failed", err.Error())
		assert.NotEqual(t, smartphone, models.Smartphone{}, "smartphone is an empty struct")
	})
}
