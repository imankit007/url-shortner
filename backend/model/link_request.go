package model

type LinkRequest struct {
	URL string `json:"url" binding:"required"`
}
