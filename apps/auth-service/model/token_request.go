package model

type TokenRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
