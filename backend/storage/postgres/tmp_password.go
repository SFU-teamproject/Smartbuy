package postgres

import (
	"database/sql"

	"github.com/sfu-teamproject/smartbuy/backend/models"
)

func (db *PostgresDB) CreateTmpPassword(tmpPassword models.TmpPassword) (models.TmpPassword, error) {
	query := `
	INSERT INTO tmp_passwords (email, password, expires_at)
	VALUES ($1, $2, $3)
	RETURNING *
	`
	row := db.QueryRow(query, tmpPassword.Email, tmpPassword.Password, tmpPassword.ExpiresAt)
	return db.extractTmpPassword(row)
}

func (db *PostgresDB) GetTmpPassword(email string) (models.TmpPassword, error) {
	query := `
	SELECT * FROM tmp_passwords WHERE email = $1
	`
	row := db.QueryRow(query, email)
	return db.extractTmpPassword(row)
}

func (db *PostgresDB) DeleteTmpPassword(email string) (models.TmpPassword, error) {
	row := db.QueryRow("DELETE FROM tmp_passwords where email = $1 returning *", email)
	return db.extractTmpPassword(row)
}

func (db *PostgresDB) extractTmpPassword(row *sql.Row) (models.TmpPassword, error) {
	t := models.TmpPassword{}
	err := row.Scan(&t.Email, &t.Password, &t.ExpiresAt)
	return t, db.wrapError(err)
}
