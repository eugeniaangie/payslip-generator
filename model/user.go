package model

type User struct {
	ID           string `db:"id" json:"id"`
	Username     string `db:"username" json:"username"`
	PasswordHash string `db:"password_hash" json:"-"`
	FullName     string `db:"full_name" json:"full_name"`
	UserRole     string `db:"user_role" json:"user_role"`
	CreatedAt    string `db:"created_at" json:"created_at"`
	UpdatedAt    string `db:"updated_at" json:"updated_at"`
}

type UserRole string

const (
	UserRoleAdmin    UserRole = "admin"
	UserRoleEmployee UserRole = "employee"
)
