package routing

import (
	"net/http"
	"net/http/pprof"

	"github.com/gorilla/mux"
	"github.com/paulwrubel/photolum/config"
	"github.com/paulwrubel/photolum/controller/cameracontroller"
	"github.com/paulwrubel/photolum/controller/helloworldcontroller"
	"github.com/paulwrubel/photolum/controller/imagecontroller"
	"github.com/paulwrubel/photolum/controller/materialcontroller"
	"github.com/paulwrubel/photolum/controller/parameterscontroller"
	"github.com/paulwrubel/photolum/controller/primitivecontroller"
	"github.com/paulwrubel/photolum/controller/rendercontroller"
	"github.com/paulwrubel/photolum/controller/scenecontroller"
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

	primitiveRouter := router.PathPrefix("/primitives").Subrouter()
	primitiveRouter.HandleFunc("", func(w http.ResponseWriter, r *http.Request) {
		primitivecontroller.GetHandler(w, r, plData, log)
	}).Methods("GET")
	primitiveRouter.HandleFunc("", func(w http.ResponseWriter, r *http.Request) {
		primitivecontroller.PostHandler(w, r, plData, log)
	}).Methods("POST")

	sceneRouter := router.PathPrefix("/scenes").Subrouter()
	sceneRouter.HandleFunc("", func(w http.ResponseWriter, r *http.Request) {
		scenecontroller.GetHandler(w, r, plData, log)
	}).Methods("GET")
	sceneRouter.HandleFunc("", func(w http.ResponseWriter, r *http.Request) {
		scenecontroller.PostHandler(w, r, plData, log)
	}).Methods("POST")

	renderRouter := router.PathPrefix("/renders").Subrouter()
	renderRouter.HandleFunc("", func(w http.ResponseWriter, r *http.Request) {
		rendercontroller.GetHandler(w, r, plData, log)
	}).Methods("GET")
	renderRouter.HandleFunc("", func(w http.ResponseWriter, r *http.Request) {
		rendercontroller.PostHandler(w, r, plData, log)
	}).Methods("POST")

	imageRouter := router.PathPrefix("/images").Subrouter()
	imageRouter.HandleFunc("", func(w http.ResponseWriter, r *http.Request) {
		imagecontroller.GetHandler(w, r, plData, log)
	}).Methods("GET")

	pprofRouter := router.PathPrefix("/debug/pprof").Subrouter()
	pprofRouter.HandleFunc("/", pprof.Index)
	pprofRouter.HandleFunc("/cmdline", pprof.Cmdline)
	pprofRouter.HandleFunc("/profile", pprof.Profile)
	pprofRouter.HandleFunc("/symbol", pprof.Symbol)
	pprofRouter.HandleFunc("/trace", pprof.Trace)

	// Manually add support for paths linked to by index page at /debug/pprof/
	pprofRouter.Handle("/goroutine", pprof.Handler("goroutine"))
	pprofRouter.Handle("/heap", pprof.Handler("heap"))
	pprofRouter.Handle("/threadcreate", pprof.Handler("threadcreate"))
	pprofRouter.Handle("/block", pprof.Handler("block"))

	return router
}
