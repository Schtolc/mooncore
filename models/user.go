package models

// User model
type User struct {
	ID       int    `json:"-"`
	Name     string `schema:"name"`
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

func (user *User) IsEmpty() bool {
	if user == (&User{}) {
		return true
	} else {
		return false
	}
}
