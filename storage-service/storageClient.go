package main

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

type StorageClient interface {
	GetBlob(blobName string) (getBlobResponse, error)
	SetBucket(name string)
	GetBlobSignedUrl(name string) (getSignedUrlResponse, error)
	Close()
}

type GCloudStorageClient struct {
	ctx    *context.Context
	client *storage.Client
	bucket *storage.BucketHandle
}

func (c *GCloudStorageClient) SetBucket(name string) {
	c.bucket = c.client.Bucket(name)
}

type getBlobResponse struct {
	fileSize int
	content  io.ReadSeeker
	fileName string
	modTime time.Time
}
type getSignedUrlResponse struct {
	Url string `json:"url"`
}
func (c *GCloudStorageClient) GetBlobSignedUrl(name string) (getSignedUrlResponse, error) {
	opts := &storage.SignedURLOptions{
		Scheme:  storage.SigningSchemeV4,
		Method:  "GET",
		Expires: time.Now().Add(15 * time.Minute),
	}
	u, err := c.bucket.SignedURL(name, opts)
	if err != nil {
		return getSignedUrlResponse{},errors.New("File not found: " + name)
	}
	return getSignedUrlResponse{Url:u}, nil
}
func (c *GCloudStorageClient) GetBlob(blobName string) (getBlobResponse, error) {
	if c.bucket == nil {
		return getBlobResponse{}, errors.New("storage bucket is not set")
	}
	blob := c.bucket.Object(blobName)
	attrs, err := blob.Attrs(*c.ctx)
	if err != nil {
		log.Fatal(err)
	}
	rc, err := blob.NewReader(*c.ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer rc.Close()
	content, err := io.ReadAll(rc)
	if err != nil {
		log.Fatal(err)
	}
	return getBlobResponse{
		content:  bytes.NewReader(content),
		fileSize: int(attrs.Size),
		fileName: attrs.Name,
		modTime: attrs.Updated,
	}, nil
}

func NewGCloudStorageClient(ctx *context.Context, secretsPath string) (*GCloudStorageClient, error) {
	client, err := storage.NewClient(*ctx, option.WithCredentialsFile(secretsPath))
	if err != nil {
		return nil, err
	}
	return &GCloudStorageClient{
		client: client,
		ctx:    ctx,
		bucket: nil,
	}, nil
}

func (c *GCloudStorageClient) Close() {
	c.client.Close()
}
