package main

import (
	"log"
	"net/http"

	"github.com/VI-IM/im_backend_go/internal/router"
)

func main() {
	// Initialize router
	router.Init()
	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", router.Router); err != nil {
		log.Fatal(err)
	}
}
