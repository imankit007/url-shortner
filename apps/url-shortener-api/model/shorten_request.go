package model

type ShortenRequest struct {
	Links []LinkRequest `json:"links" binding:"required"`
}
