package models

// User model
type User struct {
	ID       int    `json:"-"`
	Name     string `schema:"name"`
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

// IsEmpty check user struct for empty
func (user *User) IsEmpty() bool {
	if user == (&User{}) {
		return true
	}
	return false
}
