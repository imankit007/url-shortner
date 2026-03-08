package model

import "time"

type ClickEvent struct {
	ShortCode   string    `json:"short_code"`
	OriginalURL string    `json:"original_url"`
	TenantID    string    `json:"tenant_id"`
	Timestamp   time.Time `json:"timestamp"`
	UserAgent   string    `json:"user_agent"`
	Referer     string    `json:"referer"`
	IPAddress   string    `json:"ip_address"`
}
