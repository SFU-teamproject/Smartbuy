package postgres

import (
	"database/sql"

	"github.com/sfu-teamproject/smartbuy/backend/models"
)

type Review = models.Review

func (db *PostgresDB) GetReview(id int) (Review, error) {
	row := db.QueryRow(`
	SELECT reviews.id, smartphone_id, user_id, users.name, rating, comment, reviews.created_at, reviews.updated_at
	FROM reviews
	JOIN users on user_id = users.id
	WHERE reviews.id = $1
	`, id)
	review := Review{}
	err := row.Scan(&review.ID, &review.SmartphoneID, &review.UserID, &review.UserName,
		&review.Rating, &review.Comment, &review.CreatedAt, &review.UpdatedAt)
	return review, db.wrapError(err)
}

func (db *PostgresDB) GetReviews(smartphoneID int) ([]Review, error) {
	rows, err := db.Query(`
	SELECT reviews.id, smartphone_id, user_id, users.name, rating, comment, reviews.created_at, reviews.updated_at
	FROM reviews
	JOIN users on user_id = users.id
	WHERE smartphone_id = $1;
	`, smartphoneID)
	if err != nil {
		return nil, db.wrapError(err)
	}
	return db.extractReviews(rows)
}

func (db *PostgresDB) DeleteReview(ID int) (Review, error) {
	row := db.QueryRow("DELETE FROM reviews where id = $1 returning *", ID)
	return db.extractReview(row)
}

func (db *PostgresDB) UpdateReview(review Review) (Review, error) {
	query := `
	UPDATE reviews
	SET rating = $1, comment = $2, updated_at = CURRENT_TIMESTAMP
	WHERE id = $3
	RETURNING *
	`
	row := db.QueryRow(query, review.Rating, review.Comment, review.ID)
	return db.extractReview(row)
}

func (db *PostgresDB) CreateReview(review Review) (Review, error) {
	query := `
	INSERT INTO reviews (smartphone_id, user_id, rating, comment)
	VALUES ($1, $2, $3, $4)
	RETURNING *
	`
	row := db.QueryRow(query, review.SmartphoneID, review.UserID, review.Rating, review.Comment)
	return db.extractReview(row)
}

func (db *PostgresDB) extractReview(row *sql.Row) (Review, error) {
	review := Review{}
	err := row.Scan(&review.ID, &review.SmartphoneID, &review.UserID,
		&review.Rating, &review.Comment, &review.CreatedAt, &review.UpdatedAt)
	return review, db.wrapError(err)
}

func (db *PostgresDB) extractReviews(rows *sql.Rows) ([]Review, error) {
	defer rows.Close()
	reviews := []Review{}
	for rows.Next() {
		review := Review{}
		err := rows.Scan(&review.ID, &review.SmartphoneID, &review.UserID, &review.UserName,
			&review.Rating, &review.Comment, &review.CreatedAt, &review.UpdatedAt)
		if err != nil {
			return nil, db.wrapError(err)
		}
		reviews = append(reviews, review)
	}
	return reviews, nil
}
