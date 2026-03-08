package main

import (
	"github.com/imankit007/url-shortner/bootstrap"
	"github.com/imankit007/url-shortner/infrastructure"
)

func init() {
	infrastructure.InitializeRedisClient()
}

func main() {
	router, err := bootstrap.InitializeRouter()
	if err != nil {
		panic(err)
	}

	if err := router.Run(":8080"); err != nil {
		panic(err)
	}
}
