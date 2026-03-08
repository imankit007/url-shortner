package main

import (
	"log"

	"github.com/imankit007/url-shortner-redirector-service/bootstrap"
)

func main() {
	router, err := bootstrap.InitializeRouter()
	if err != nil {
		log.Fatal(err)
	}

	if err := router.Run(":8082"); err != nil {
		log.Fatal(err)
	}
}
