package main

import (
	"log"
	"net/http"
	_ "see-log-go/router"
)

func main() {
	log.Fatalln("error", http.ListenAndServe(":8080", nil))
}
