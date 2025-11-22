package main

import (
	"context"
	"crypto/rand"
	"math/big"
	"time"

	"github.com/imankit007/url-shortner/inits"
	"go.mongodb.org/mongo-driver/mongo"
)

func init() {

}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func RandAlphaNumeric(n int) string {
	b := make([]rune, n)
	for i := range b {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		b[i] = letters[num.Int64()]
	}
	return string(b)
}

func main() {

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var orderCollection *mongo.Collection = inits.OpenCollection(inits.Client, "url-collection")

	for i := 0; i < 100; i++ {
		orderCollection.InsertOne(ctx, map[string]string{"value": RandAlphaNumeric(20)})
		println("Added %d entry to collection", i+1)
	}

	
}
