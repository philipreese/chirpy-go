package main

import (
	"log"
	"net/http"
)

func main() {
	const port = "8080"
	server := &http.Server{
		Handler: http.NewServeMux(),
		Addr:    ":" + port,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal("failed to start server")
	}
	log.Printf("Serving on port: %s\n", port)
}