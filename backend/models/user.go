package models

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type Role string

const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
)

type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Password  string    `json:"password,omitzero"`
	Role      Role      `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	Cart      Cart      `json:"cart,omitzero"`
}

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
