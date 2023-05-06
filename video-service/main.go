package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

)

func main() {
	port := os.Getenv("PORT")
	storageServiceURL := os.Getenv("STORAGE_SERVICE_URL")
	if storageServiceURL == "" {
		log.Fatal("Please specify storage service")
	}
	if port == "" {
		port = "8080"
	}
	srv := &http.Server{
		Handler: Routes(),
		Addr: ":" + port,
		ErrorLog: log.New(os.Stderr, "ERROR\t", log.Ldate | log.Ltime | log.Lshortfile),
	}
	fmt.Println("Running on port: ", port)
	log.Fatal(srv.ListenAndServe())

}