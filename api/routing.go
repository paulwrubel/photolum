package api

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/paulwrubel/photolum/api/helloworld"
	"github.com/paulwrubel/photolum/api/scenes"
	"github.com/paulwrubel/photolum/config"
)

// ListenAndServe starts the API server on port 8080
func ListenAndServe(plData *config.PhotolumData) {
	router := getRouter(plData)
	go func() {
		if err := http.ListenAndServe(":8080", router); err != nil {
			fmt.Printf("Error in api.ListenAndServe(): %s\n", err.Error())
		}
	}()
}

// GetRouter gets the main router used for the API
func getRouter(plData *config.PhotolumData) *mux.Router {
	router := mux.NewRouter().PathPrefix("/v1").Subrouter()

	helloWorldRouter := router.PathPrefix("/helloworld").Subrouter()
	helloWorldRouter.HandleFunc("", helloworld.GetHandler).Methods("GET")

	scenesRouter := router.PathPrefix("/scenes").Subrouter()
	scenesRouter.HandleFunc("", func(w http.ResponseWriter, r *http.Request) {
		scenes.ScenesPostHandler(w, r, plData)
	}).Methods("POST")
	scenesRouter.HandleFunc("", func(w http.ResponseWriter, r *http.Request) {
		scenes.ScenesGetHandler(w, r, plData)
	}).Methods("GET")
	scenesRouter.HandleFunc("/{scene_id}", func(w http.ResponseWriter, r *http.Request) {
		scenes.SceneIDGetHandler(w, r, plData)
	}).Methods("GET")
	scenesRouter.HandleFunc("/{scene_id}/image", func(w http.ResponseWriter, r *http.Request) {
		scenes.SceneIDImageGetHandler(w, r, plData)
	}).Methods("GET")
	scenesRouter.HandleFunc("/{scene_id}", func(w http.ResponseWriter, r *http.Request) {
		scenes.SceneIDDeleteHandler(w, r, plData)
	}).Methods("DELETE")

	renderRouter := scenesRouter.PathPrefix("/{scene_id}/render").Subrouter()
	renderRouter.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		scenes.ScenesIDRenderStatusGetHandler(w, r, plData)
	}).Methods("GET")
	renderRouter.HandleFunc("/start", func(w http.ResponseWriter, r *http.Request) {
		scenes.ScenesIDRenderStartPostHandler(w, r, plData)
	}).Methods("POST")
	renderRouter.HandleFunc("/stop", func(w http.ResponseWriter, r *http.Request) {
		scenes.ScenesIDRenderStopPostHandler(w, r, plData)
	}).Methods("POST")

	return router
}
