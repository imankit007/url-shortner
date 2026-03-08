package model

type User struct {
	UserID       string
	TenantID     string
	Email        string
	PasswordHash string
}
