package models

// User model
type User struct {
	ID       int    `json:"-"`
	Name     string `schema:"name"`
	Email    string `schema:"email"`
	Password string `schema:"password"`
}


