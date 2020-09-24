package api

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
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
	helloWorldRouter.HandleFunc("/", HelloWorldHandler).Methods("GET")

	return router
}
