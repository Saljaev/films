package models

type User struct {
	Id       int    `json:"id,omitempty"`
	Login    string `json:"login,omitempty"`
	Password string `json:"password,omitempty"`
	Role     string `json:"role,omitempty"`
}
