package routing

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/paulwrubel/photolum/config"
	"github.com/paulwrubel/photolum/controller/cameracontroller"
	"github.com/paulwrubel/photolum/controller/helloworldcontroller"
	"github.com/paulwrubel/photolum/controller/materialcontroller"
	"github.com/paulwrubel/photolum/controller/parameterscontroller"
	"github.com/paulwrubel/photolum/controller/texturecontroller"
	"github.com/sirupsen/logrus"
)

// ListenAndServe starts the API server on port 8080
func ListenAndServe(plData *config.PhotolumData, log *logrus.Logger) {
	router := getRouter(plData, log)
	go func() {
		if err := http.ListenAndServe(":8080", router); err != nil {
			log.WithError(err).Error("error in http.ListenAndServer()")
		}
	}()
}

// GetRouter gets the main router used for the API
func getRouter(plData *config.PhotolumData, log *logrus.Logger) *mux.Router {
	router := mux.NewRouter().PathPrefix("/v1").Subrouter()

	helloWorldRouter := router.PathPrefix("/helloworld").Subrouter()
	helloWorldRouter.HandleFunc("", func(w http.ResponseWriter, r *http.Request) {
		helloworldcontroller.GetHandler(w, r, log)
	}).Methods("GET")

	parametersRouter := router.PathPrefix("/parameters").Subrouter()
	parametersRouter.HandleFunc("", func(w http.ResponseWriter, r *http.Request) {
		parameterscontroller.GetHandler(w, r, plData, log)
	}).Methods("GET")
	parametersRouter.HandleFunc("", func(w http.ResponseWriter, r *http.Request) {
		parameterscontroller.PostHandler(w, r, plData, log)
	}).Methods("POST")

	cameraRouter := router.PathPrefix("/cameras").Subrouter()
	cameraRouter.HandleFunc("", func(w http.ResponseWriter, r *http.Request) {
		cameracontroller.GetHandler(w, r, plData, log)
	}).Methods("GET")
	cameraRouter.HandleFunc("", func(w http.ResponseWriter, r *http.Request) {
		cameracontroller.PostHandler(w, r, plData, log)
	}).Methods("POST")

	textureRouter := router.PathPrefix("/textures").Subrouter()
	textureRouter.HandleFunc("", func(w http.ResponseWriter, r *http.Request) {
		texturecontroller.GetHandler(w, r, plData, log)
	}).Methods("GET")
	textureRouter.HandleFunc("", func(w http.ResponseWriter, r *http.Request) {
		texturecontroller.PostHandler(w, r, plData, log)
	}).Methods("POST")

	materialRouter := router.PathPrefix("/materials").Subrouter()
	materialRouter.HandleFunc("", func(w http.ResponseWriter, r *http.Request) {
		materialcontroller.GetHandler(w, r, plData, log)
	}).Methods("GET")
	materialRouter.HandleFunc("", func(w http.ResponseWriter, r *http.Request) {
		materialcontroller.PostHandler(w, r, plData, log)
	}).Methods("POST")

	// scenesRouter := router.PathPrefix("/scenes").Subrouter()
	// scenesRouter.HandleFunc("", func(w http.ResponseWriter, r *http.Request) {
	// 	scenecontroller.ScenesPostHandler(w, r, plData)
	// }).Methods("POST")
	// scenesRouter.HandleFunc("", func(w http.ResponseWriter, r *http.Request) {
	// 	scenecontroller.ScenesGetHandler(w, r, plData)
	// }).Methods("GET")
	// scenesRouter.HandleFunc("/{scene_id}", func(w http.ResponseWriter, r *http.Request) {
	// 	scenecontroller.SceneIDGetHandler(w, r, plData)
	// }).Methods("GET")
	// scenesRouter.HandleFunc("/{scene_id}/image", func(w http.ResponseWriter, r *http.Request) {
	// 	scenecontroller.SceneIDImageGetHandler(w, r, plData)
	// }).Methods("GET")
	// scenesRouter.HandleFunc("/{scene_id}", func(w http.ResponseWriter, r *http.Request) {
	// 	scenecontroller.SceneIDDeleteHandler(w, r, plData)
	// }).Methods("DELETE")

	// renderRouter := scenesRouter.PathPrefix("/{scene_id}/render").Subrouter()
	// renderRouter.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
	// 	scenecontroller.ScenesIDRenderStatusGetHandler(w, r, plData)
	// }).Methods("GET")
	// renderRouter.HandleFunc("/start", func(w http.ResponseWriter, r *http.Request) {
	// 	scenecontroller.ScenesIDRenderStartPostHandler(w, r, plData)
	// }).Methods("POST")
	// renderRouter.HandleFunc("/stop", func(w http.ResponseWriter, r *http.Request) {
	// 	scenecontroller.ScenesIDRenderStopPostHandler(w, r, plData)
	// }).Methods("POST")

	return router
}
