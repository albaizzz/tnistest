package http

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/tnistest/config"
	"github.com/urfave/negroni"
	"github.com/wptest/pkg/logrus"
)

type Routes struct {
	Config        *config.Config
	DeviceHandler handler.IDeviceHandler
}

// Main Router
func (r *Routes) NewRoutes() http.Handler {
	// define route
	router := mux.NewRouter().StrictSlash(false)
	route := router.PathPrefix(r.Config.Api.Prefix).Subrouter()

	// health-check
	route.HandleFunc("/health-check", handler.GetHealthCheck).Methods(http.MethodGet)
	// messages
	route.HandleFunc("/data", r.DeviceHandler.InsertDevice).Methods(http.MethodPost)
	route.HandleFunc("/data/{id}", r.DeviceHandler.GetDevice).Methods(http.MethodGet)
	route.HandleFunc("/data", r.DeviceHandler.GetDevices).Methods(http.MethodGet)

	// Use Negroni Log Router
	n := negroni.Classic()
	recovery := negroni.NewRecovery() // Panic handler
	if r.Config.App.Debug == false {
		recovery.PrintStack = false
	}

	n.Use(logrus.NewLoggerMiddleware(r.Config.App.Name))
	n.UseHandler(router)
	return n
}
