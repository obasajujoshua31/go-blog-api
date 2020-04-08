package main

import (
	"go-blog-api/api"
	"log"
)

func main() {
	if err := api.Start(); err != nil {
		log.Fatalf("Error in starting server: %s", err)
	}
}
