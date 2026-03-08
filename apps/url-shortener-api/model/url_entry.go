package model

type URLEntry struct {
	OriginalURL     string `bson:"url"`
	Code            int64  `bson:"code"`
	TenantID        string `bson:"tenant_id"`
	CreatedByUserID string `bson:"created_by_user_id"`
}
