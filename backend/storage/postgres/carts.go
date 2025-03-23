package postgres

import (
	"database/sql"

	"github.com/sfu-teamproject/smartbuy/backend/models"
)

type Cart = models.Cart

func (db *PostgresDB) GetCarts() ([]Cart, error) {
	rows, err := db.Query("SELECT * FROM carts")
	if err != nil {
		return nil, db.wrapError(err)
	}
	return db.extractCarts(rows)
}

func (db *PostgresDB) GetCart(ID int) (Cart, error) {
	row := db.QueryRow("SELECT * FROM carts WHERE id = $1", ID)
	return db.extractCart(row)
}

func (db *PostgresDB) GetCartByUserID(userID int) (Cart, error) {
	row := db.QueryRow("SELECT * FROM carts WHERE user_id = $1", userID)
	return db.extractCart(row)
}

func (db *PostgresDB) extractCart(row *sql.Row) (Cart, error) {
	cart := Cart{}
	err := row.Scan(&cart.ID, &cart.UserID, &cart.CreatedAt, &cart.UpdatedAt)
	return cart, db.wrapError(err)
}

func (db *PostgresDB) extractCarts(rows *sql.Rows) ([]Cart, error) {
	defer rows.Close()
	carts := []Cart{}
	for rows.Next() {
		cart := Cart{}
		err := rows.Scan(&cart.ID, &cart.UserID, &cart.CreatedAt, &cart.UpdatedAt)
		if err != nil {
			return nil, db.wrapError(err)
		}
		carts = append(carts, cart)
	}
	return carts, nil
}
