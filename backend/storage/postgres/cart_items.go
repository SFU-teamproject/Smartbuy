package postgres

import (
	"database/sql"

	"github.com/sfu-teamproject/smartbuy/backend/models"
)

type CartItem = models.CartItem

func (db *PostgresDB) AddToCart(cartItem models.CartItem) (CartItem, error) {
	row := db.QueryRow("INSERT INTO cart_items (cart_id, smartphone_id) values ($1, $2) RETURNING *",
		cartItem.CartID, cartItem.SmartphoneID)
	return db.extractCartItem(row)
}

func (db *PostgresDB) GetCartItem(ID int) (CartItem, error) {
	row := db.QueryRow("SELECT * from cart_items where id = $1", ID)
	return db.extractCartItem(row)
}

func (db *PostgresDB) GetCartItems(cartID int) ([]CartItem, error) {
	rows, err := db.Query("SELECT * from cart_items where cart_id = $1", cartID)
	if err != nil {
		return nil, db.wrapError(err)
	}
	return db.extractCartItems(rows)
}

func (db *PostgresDB) SetQuantity(cartItem CartItem) (CartItem, error) {
	row := db.QueryRow("UPDATE cart_items SET quantity = $1 WHERE id = $2 and cart_id = $3 RETURNING *",
		cartItem.Quantity, cartItem.ID, cartItem.CartID)
	return db.extractCartItem(row)
}

func (db *PostgresDB) DeleteFromCart(cartID, itemID int) (CartItem, error) {
	row := db.QueryRow("DELETE FROM cart_items where cart_id = $1 and id = $2 returning *",
		cartID, itemID)
	return db.extractCartItem(row)
}

func (db *PostgresDB) extractCartItem(row *sql.Row) (CartItem, error) {
	ci := CartItem{}
	err := row.Scan(&ci.ID, &ci.CartID, &ci.SmartphoneID, &ci.Quantity)
	return ci, db.wrapError(err)
}

func (db *PostgresDB) extractCartItems(rows *sql.Rows) ([]CartItem, error) {
	defer rows.Close()
	cis := []CartItem{}
	for rows.Next() {
		ci := CartItem{}
		err := rows.Scan(&ci.ID, &ci.CartID, &ci.SmartphoneID, &ci.Quantity)
		if err != nil {
			return nil, db.wrapError(err)
		}
		cis = append(cis, ci)
	}
	return cis, nil
}
