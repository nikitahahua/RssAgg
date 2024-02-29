package main

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func main() {
	log.SetPrefix("georgian greetings: ")
	log.SetFlags(0)
	godotenv.Load(".env")

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatalln("cant connect to the port")
	}

	router := chi.NewRouter()

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + port,
	}

	log.Printf("Server starting on pory %v", port)
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("port: " + port)
}
