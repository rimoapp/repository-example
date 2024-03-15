package main

import (
	"log"
	"net/http"
	"os"

	"github.com/rimoapp/repository-example/router"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Listening on port %s", port)

	rimoRouter, err := router.NewRouter()
	if err != nil {
		log.Fatal(err)
	}
	handler := rimoRouter.Handler
	srv := http.Server{Addr: ":" + port, Handler: handler}
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}

}
