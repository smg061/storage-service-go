package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type videoRequestResponse struct {
	VideoUrl string `json:"videoUrl"`
}
func ShowVideo(w http.ResponseWriter, r *http.Request) {
	path:= r.URL.Query().Get("path")
	if path == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	storageServiceURL := os.Getenv("STORAGE_SERVICE_URL")

	fwdReq, err := http.Get(fmt.Sprintf("%s/videos?path=%s", storageServiceURL, path))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer fwdReq.Body.Close()
	var response videoRequestResponse

	if err := json.NewDecoder(fwdReq.Body).Decode(&response); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response.VideoUrl))
}