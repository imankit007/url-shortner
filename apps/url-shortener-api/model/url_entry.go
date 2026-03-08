package model

type URLEntry struct {
	OriginalURL     string `bson:"url" json:"url"`
	Code            int64  `bson:"code" json:"code"`
	TenantID        string `bson:"tenant_id" json:"tenant_id"`
	CreatedByUserID string `bson:"created_by_user_id" json:"created_by_user_id"`
}
