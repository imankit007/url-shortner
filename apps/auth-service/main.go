package main

import (
	"log"
	"net/http"

	"github.com/imankit007/url-shortner-auth-service/bootstrap"
)

func main() {
	handler, err := bootstrap.InitializeRouter()
	if err != nil {
		log.Fatal(err)
	}

	if err := http.ListenAndServe(":8081", handler); err != nil {
		log.Fatal(err)
	}
}
