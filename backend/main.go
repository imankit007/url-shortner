package main

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/imankit007/url-shortner/inits"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var urlMap =  make(map[string]string)

func init() {
	inits.ConnectRedis()
}

type ShortenRequest struct {
	URL string `json:"url" binding:"required"`
}


func main() {
	
	var orderCollection *mongo.Collection = inits.OpenCollection(inits.Client, "url-collection")
	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var orders []bson.M
		cursor, err := orderCollection.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		defer cursor.Close(ctx)

		if err = cursor.All(ctx, &orders); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		fmt.Println(orders)
		c.JSON(http.StatusOK, orders)
	})

	router.GET("/:code", func(ctx *gin.Context) {
		var url string = urlMap[ctx.Param("code")]
		ctx.Redirect(http.StatusFound, url)
	}) 

	router.POST("/shorten", func(ctx *gin.Context) {
		var req ShortenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "url is required"})
		return
	}

	code := generateCode(6)
	urlMap[code] = req.URL

	ctx.JSON(http.StatusOK, gin.H{
		"short_url": "http://localhost:8080/" + code,
	})

	})


	router.Run(":8080")

	fmt.Println("Hello, World!")

}

func generateCode(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	rand.Seed(time.Now().UnixNano())

	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
