package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/imankit007/url-shortner/inits"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)


func init() {
	inits.ConnectRedis()
}

type ShortenRequest struct {
	Links []LinkRequest `json:"links" binding:"required"`
}

type LinkRequest struct{
	Url string `json:"url" binding:"required"`
}

type URLEntry struct{
	Url string `bson:"url"`
	Code string `bson:"code"`
}

const base62chars = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func main() {

	var mainCounter uint64 =  uint64(123456) 

	var orderCollection *mongo.Collection = inits.OpenCollection(inits.Client, "url-collection")

	orderCollection.DeleteMany(context.TODO() ,bson.M{})

	router := gin.Default()
	router.Use(cors.Default())
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

		code := ctx.Param("code")
		
		var counter_code uint64 = decodeBase62(code)

		println(counter_code)

		var dbCtx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var result URLEntry

		orderCollection.FindOne(dbCtx, bson.M{
			"code": counter_code,
		}).Decode(&result)

		fmt.Println(result)

		ctx.Redirect(http.StatusFound, "https://" + result.Url)
	})

	router.POST("/shorten", func(ctx *gin.Context) {
		var req ShortenRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "url is required"})
			return
		}

		var response []map[string]any

		for _, r := range req.Links {
			currCount := mainCounter
			code := encodeBase62(currCount)
			mainCounter++
			
			doc := map[string]interface{}{
				"url": r.Url,
				"code": currCount,
			}
			var dbCtx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
			defer cancel()
			orderCollection.InsertOne(dbCtx, doc);

			response = append(response, map[string]any{
				"long_url" : r.Url,
				"short_url" : "http://localhost:8080/" + code,
		})
	
	}

		ctx.JSON(http.StatusOK, response)
	})


	router.Run(":8080")

	fmt.Println("Hello, World!")

}


func encodeBase62(num uint64) string {
	if num == 0 {
		return "0"
	}

	buf := make([]byte, 0, 11) // max length for uint64

	for num > 0 {
		rem := num % 62
		buf = append(buf, base62chars[rem])
		num /= 62
	}

	// reverse
	for i, j := 0, len(buf)-1; i < j; i, j = i+1, j-1 {
		buf[i], buf[j] = buf[j], buf[i]
	}

	return string(buf)
}

func decodeBase62(s string) uint64 {
	var num uint64
	for _, c := range s {
		num = num*62 + uint64(indexOf(byte(c)))
	}
	return num
}

func indexOf(c byte) int {
	for i := range base62chars {
		if base62chars[i] == c {
			return i
		}
	}
	return -1
}