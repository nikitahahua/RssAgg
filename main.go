package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	log.SetPrefix("georgian greetings: ")
	log.SetFlags(0)
	//Load()

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatalln("cant connect to the port")
	}

	fmt.Println("port: " + port)
}
