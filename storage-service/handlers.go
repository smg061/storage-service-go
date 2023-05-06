package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type videoRequest struct {
	Name string `json:"name"`
}
type VideoUrlResponse struct {
	VideoUrl string `json:"videoUrl"`
}

type cacheFile struct {
	content io.ReadSeeker
	modTime time.Time
}

func (app *app) GetVideo(w http.ResponseWriter, r *http.Request) {
	var req videoRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	video, err := app.storageClient.GetBlob(req.Name)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	http.ServeContent(w, r, video.fileName, time.Now(), video.content)
}

func (app *app) GetVideoFromPath(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	if path == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	video, err := app.storageClient.GetBlobSignedUrl(path)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	resp := &VideoUrlResponse{
		VideoUrl: video.Url,
	}
	bytes, err := json.Marshal(&resp)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(bytes)
}

// this is a very silly idea
func (app *app) GetVideoUrlPath(w http.ResponseWriter, r *http.Request) {
	video, found := app.cache["feesh.mp4"]
	if !found {
		app.mutex.Lock()
		defer app.mutex.Unlock()
		blob, err := app.storageClient.GetBlob("feesh.mp4")
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		video = &cacheFile{
			content: blob.content,
			modTime: blob.modTime,
		}
		app.cache[blob.fileName] = video
	}
	http.ServeContent(w, r, "feesh.mp4", time.Now(), video.content)
}

func (app *app) GetSignedVideoUrl(w http.ResponseWriter, r *http.Request) {
	var req videoRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fmt.Printf("%s", req)
	u, err := app.storageClient.GetBlobSignedUrl(req.Name)
	fmt.Printf("%s", u)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	resp := &VideoUrlResponse{
		VideoUrl: u.Url,
	}
	bytes, err := json.Marshal(&resp)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(bytes)

}
