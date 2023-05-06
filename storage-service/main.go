package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
)

type app struct {
	storageClient StorageClient
	cache         map[string]*cacheFile
	mutex         *sync.RWMutex
}

func main() {

	PORT := os.Getenv("PORT")
	ctx := context.Background()
	secrets := os.Getenv("GCLOUD_STORAGE_SECRETS")
	if secrets == "" {
		log.Fatal("Storage secrets not provided")
	}
	client, err := NewGCloudStorageClient(&ctx, secrets)

	if err != nil {
		log.Fatal(err)
	}
	app := &app{storageClient: client, cache: map[string]*cacheFile{}, mutex: &sync.RWMutex{}}
	defer app.storageClient.Close()

	bucketname := os.Getenv("BUCKETNAME")
	if bucketname == "" {
		log.Fatal("bucketname not specified")
	}
	app.storageClient.SetBucket(bucketname)
	srv := &http.Server{
		Handler:  app.Routes(),
		Addr:     ":" + PORT,
		ErrorLog: log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
	}
	fmt.Println("Running on port: ", PORT)
	log.Fatal(srv.ListenAndServe()) //TLS("/Users/max/Documents/github/video-service-go/shared/tls/basic-certificate.cert", "/Users/max/Documents/github/video-service-go/shared/tls/basic-private-key.key"))
}
