package api

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/paulwrubel/photolum/api/helloworld"
	"github.com/paulwrubel/photolum/api/scenes"
)

// ListenAndServe starts the API server on port 8080
func ListenAndServe() {
	router := getRouter()
	go func() {
		if err := http.ListenAndServe(":8080", router); err != nil {
			fmt.Printf("Error in api.ListenAndServe(): %s\n", err.Error())
		}
	}()
}

// GetRouter gets the main router used for the API
func getRouter() *mux.Router {
	router := mux.NewRouter()

	helloWorldRouter := router.PathPrefix("/helloworld").Subrouter()
	helloWorldRouter.HandleFunc("", helloworld.GetHandler).Methods("GET")

	scenesRouter := router.PathPrefix("/scenes").Subrouter()
	scenesRouter.HandleFunc("", scenes.ScenesPostHandler).Methods("POST")
	scenesRouter.HandleFunc("", scenes.ScenesGetHandler).Methods("GET")
	scenesRouter.HandleFunc("/{scene_id}", scenes.SceneIDGetHandler).Methods("GET")
	scenesRouter.HandleFunc("/{scene_id}/image", scenes.SceneIDImageGetHandler).Methods("GET")
	scenesRouter.HandleFunc("/{scene_id}", scenes.SceneIDDeleteHandler).Methods("DELETE")

	return router
}
