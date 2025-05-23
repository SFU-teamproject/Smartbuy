package postgres

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/sfu-teamproject/smartbuy/backend/models"
)

type Smartphone = models.Smartphone

func (db *PostgresDB) GetSmartphones() ([]Smartphone, error) {
	rows, err := db.Query("SELECT * FROM smartphones order by price")
	if err != nil {
		return nil, db.wrapError(err)
	}
	return db.extractSmartphones(rows)
}

func (db *PostgresDB) GetSmartphonesByIDs(IDs []int) ([]Smartphone, error) {
	var IDsStr strings.Builder
	for i, ID := range IDs {
		if i != 0 {
			IDsStr.WriteByte(',')
		}
		IDsStr.WriteString(strconv.Itoa(ID))
	}
	query := fmt.Sprintf("SELECT * FROM smartphones WHERE id IN (%s) order by price", IDsStr.String())
	rows, err := db.Query(query)
	if err != nil {
		return nil, db.wrapError(err)
	}
	return db.extractSmartphones(rows)
}

func (db *PostgresDB) GetSmartphone(id int) (Smartphone, error) {
	row := db.QueryRow("SELECT * FROM smartphones WHERE id = $1", id)
	return db.extractSmartphone(row)
}

func (db *PostgresDB) DeleteSmartphone(id int) (Smartphone, error) {
	row := db.QueryRow("DELETE FROM smartphones where id = $1 returning *", id)
	return db.extractSmartphone(row)
}

func (db *PostgresDB) UpdateSmartphone(sm Smartphone) (Smartphone, error) {
	query := `
	UPDATE smartphones
	SET model = $1, producer = $2, memory = $3, ram = $4, display_size = $5,
	ratings_sum = $6, ratings_count = $7, price = $8, image_path = $9, description = $10
	WHERE id = $11
	RETURNING *
	`
	row := db.QueryRow(query, sm.Model, sm.Producer, sm.Memory, sm.Ram, sm.DisplaySize,
		sm.RatingsSum, sm.RatingsCount, sm.Price, sm.ImagePath, sm.Description, sm.ID)
	return db.extractSmartphone(row)
}

func (db *PostgresDB) CreateSmartphone(sm Smartphone) (Smartphone, error) {
	query := `
	INSERT INTO smartphones (model, producer, memory, ram, display_size,
	ratings_sum, ratings_count, price, image_path, description)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	RETURNING *
	`
	row := db.QueryRow(query, sm.Model, sm.Producer, sm.Memory, sm.Ram, sm.DisplaySize,
		sm.RatingsSum, sm.RatingsCount, sm.Price, sm.ImagePath, sm.Description)
	return db.extractSmartphone(row)
}

func (db *PostgresDB) extractSmartphone(row *sql.Row) (Smartphone, error) {
	sm := Smartphone{}
	err := row.Scan(&sm.ID, &sm.Model, &sm.Producer, &sm.Memory, &sm.Ram, &sm.DisplaySize,
		&sm.Price, &sm.RatingsSum, &sm.RatingsCount, &sm.ImagePath, &sm.Description)
	return sm, db.wrapError(err)
}

func (db *PostgresDB) extractSmartphones(rows *sql.Rows) ([]Smartphone, error) {
	defer rows.Close()
	smartphones := []Smartphone{}
	for rows.Next() {
		sm := Smartphone{}
		err := rows.Scan(&sm.ID, &sm.Model, &sm.Producer, &sm.Memory, &sm.Ram, &sm.DisplaySize,
			&sm.Price, &sm.RatingsSum, &sm.RatingsCount, &sm.ImagePath, &sm.Description)
		if err != nil {
			return nil, db.wrapError(err)
		}
		smartphones = append(smartphones, sm)
	}
	return smartphones, nil
}
