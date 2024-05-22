package models

const (
	OrdinaryUser = iota
	Admin
)

type User struct {
	ID       int    `json:"id" db:"id"`
	Username string `json:"username" db:"username"`
	Admin    bool   `json:"admin" db:"admin"`
}
