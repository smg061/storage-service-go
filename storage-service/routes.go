package main

import (
	"github.com/gorilla/mux"
)

func (app *app) Routes() *mux.Router {
	mux := mux.NewRouter()
	mux.HandleFunc("/api/videos", app.GetSignedVideoUrl).Methods("POST")
	mux.HandleFunc("/videos", app.GetVideoFromPath).Methods("GET")
	return mux
}
