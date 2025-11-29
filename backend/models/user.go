package models

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Avatar    *string   `json:"avatar"`
	Password  *string   `json:"-"`
	Role      Role      `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	Cart      Cart      `json:"cart,omitzero"`
}

type SignUpRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	User  User   `json:"user"`
	Token string `json:"token"`
}

type UpdateRequest struct {
	Name     string `json:"name"`
	Avatar   string `json:"avatar"`
	Password string `json:"password"`
}

type Role string

const (
	RoleAdmin = "admin"
	RoleUser  = "user"
)

func (r *Role) isValid() bool {
	return *r == RoleUser || *r == RoleAdmin
}

func (r *Role) UnmarshalJSON(b []byte) error {
	var roleStr string
	if err := json.Unmarshal(b, &roleStr); err != nil {
		return err
	}
	roleStr = strings.ToLower(roleStr)
	role := Role(roleStr)
	if !role.isValid() {
		return fmt.Errorf("invalid role: %s", roleStr)
	}
	*r = role
	return nil
}
