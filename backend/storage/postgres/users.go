package postgres

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/sfu-teamproject/smartbuy/backend/models"
)

type User = models.User

func (db *PostgresDB) GetUsers() ([]User, error) {
	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		return nil, db.wrapError(err)
	}
	return db.extractUsers(rows)
}

func (db *PostgresDB) GetUser(ID int) (User, error) {
	row := db.QueryRow("SELECT * FROM users WHERE id = $1", ID)
	return db.extractUser(row)
}

func (db *PostgresDB) GetUserByEmail(email string) (User, error) {
	row := db.QueryRow("SELECT * FROM users WHERE email = $1", email)
	return db.extractUser(row)
}

func (db *PostgresDB) DeleteUser(ID int) (User, error) {
	row := db.QueryRow("DELETE FROM users where id = $1 returning *", ID)
	return db.extractUser(row)
}

func (db *PostgresDB) UpdateUser(userID int, updates map[string]any) (User, error) {
	var query strings.Builder
	query.WriteString("UPDATE users SET ")
	args := make([]any, 0, len(updates))
	argId := 1
	for field, value := range updates {
		query.WriteString(fmt.Sprintf("%s = $%d", field, argId))
		args = append(args, value)
		if argId < len(updates) {
			query.WriteString(", ")
		}
		argId++
	}
	query.WriteString(fmt.Sprintf(" WHERE id = $%d RETURNING *", argId))
	args = append(args, userID)
	row := db.QueryRow(query.String(), args...)
	return db.extractUser(row)
}

func (db *PostgresDB) CreateUser(user User) (User, error) {
	query := `
	INSERT INTO users (name, email, avatar, password)
	VALUES ($1, $2, $3, $4)
	RETURNING *
	`
	row := db.QueryRow(query, user.Name, user.Email, user.Avatar, user.Password)
	return db.extractUser(row)
}

func (db *PostgresDB) extractUser(row *sql.Row) (User, error) {
	user := User{}
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Avatar, &user.Password, &user.Role, &user.CreatedAt)
	return user, db.wrapError(err)
}

func (db *PostgresDB) extractUsers(rows *sql.Rows) ([]User, error) {
	defer rows.Close()
	users := []User{}
	for rows.Next() {
		user := User{}
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Avatar, &user.Password, &user.Role, &user.CreatedAt)
		if err != nil {
			return nil, db.wrapError(err)
		}
		users = append(users, user)
	}
	return users, nil
}
