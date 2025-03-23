package postgres

import (
	"database/sql"

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

func (db *PostgresDB) GetUserByName(name string) (User, error) {
	row := db.QueryRow("SELECT * FROM users WHERE name = $1", name)
	return db.extractUser(row)
}

func (db *PostgresDB) DeleteUser(ID int) (User, error) {
	row := db.QueryRow("DELETE FROM users where id = $1 returning *", ID)
	return db.extractUser(row)
}

func (db *PostgresDB) UpdateUser(user User) (User, error) {
	query := `
	UPDATE users
	SET name = $1, password = $2, role = $3
	WHERE id = $4
	RETURNING *
	`
	row := db.QueryRow(query, user.Name, user.Password, user.Role, user.ID)
	return db.extractUser(row)
}

func (db *PostgresDB) CreateUser(user User) (User, error) {
	query := `
	INSERT INTO users (name, password)
	VALUES ($1, $2)
	RETURNING *
	`
	row := db.QueryRow(query, user.Name, user.Password)
	return db.extractUser(row)
}

func (db *PostgresDB) extractUser(row *sql.Row) (User, error) {
	user := User{}
	err := row.Scan(&user.ID, &user.Name, &user.Password, &user.Role, &user.CreatedAt)
	return user, db.wrapError(err)
}

func (db *PostgresDB) extractUsers(rows *sql.Rows) ([]User, error) {
	defer rows.Close()
	users := []User{}
	for rows.Next() {
		user := User{}
		err := rows.Scan(&user.ID, &user.Name, &user.Password, &user.Role, &user.CreatedAt)
		if err != nil {
			return nil, db.wrapError(err)
		}
		user.Password = ""
		users = append(users, user)
	}
	return users, nil
}
