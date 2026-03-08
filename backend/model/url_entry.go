package model

type URLEntry struct {
	OriginalURL string `bson:"url"`
	Code        int64  `bson:"code"`
}
