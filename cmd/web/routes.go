package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	// Initialize a new mux using the pat package
	router := httprouter.New()
	router.RedirectFixedPath = true
	router.RedirectTrailingSlash = true

	router.Handler(http.MethodGet, "/", http.HandlerFunc(app.home))
	router.Handler(http.MethodGet, "/recon", http.HandlerFunc(app.home))
	router.Handler(http.MethodGet, "/filter", http.HandlerFunc(app.home))
	router.Handler(http.MethodGet, "/capture", http.HandlerFunc(app.home))
	router.Handler(http.MethodGet, "/scheduler", http.HandlerFunc(app.home))
	router.Handler(http.MethodGet, "/attack", http.HandlerFunc(app.home))

	// Ping handler for testing
	router.Handler(http.MethodGet, "/ping", http.HandlerFunc(app.ping))

	// API endpoint for testing
	router.Handler(http.MethodGet, "/api/test", http.HandlerFunc(app.apiTest))

	// System information handler
	router.Handler(http.MethodGet, "/api/system", http.HandlerFunc(app.systemInfo))

	// Security overview handler
	router.Handler(http.MethodGet, "/api/security", http.HandlerFunc(app.getSec))

	// Capture file handlers
	router.Handler(http.MethodGet, "/api/captures", http.HandlerFunc(app.getCaptures))
	router.Handler(http.MethodGet, "/api/capture/:id", http.HandlerFunc(app.getCaptureWithId))
	router.Handler(http.MethodPost, "/api/captures", app.readSimpleAction(http.HandlerFunc(app.captureAction)))

	// Access Point handler for the recon
	router.Handler(http.MethodGet, "/api/accesspoints", http.HandlerFunc(app.getAP))
	router.Handler(http.MethodPost, "/api/accesspoints", app.readSimpleAction(http.HandlerFunc(app.apAction)))

	// Client Handler for the recon
	router.Handler(http.MethodGet, "/api/clients", http.HandlerFunc(app.getClients))
	router.Handler(http.MethodPost, "/api/clients", app.readSimpleAction(http.HandlerFunc(app.clientAction)))

	// Network interfaces handler
	router.Handler(http.MethodGet, "/api/interfaces", http.HandlerFunc(app.getInterfaces))
	router.Handler(http.MethodPost, "/api/interfaces/:id",
		app.findInterfaceById(
			app.readInterfaceActoin(
				http.HandlerFunc(app.interfaceAction))))

	// Filter / White- / Black-list
	router.Handler(http.MethodGet, "/api/filter/:id",
		app.findFilterById(
			http.HandlerFunc(app.getFilters)))
	router.Handler(http.MethodPost, "/api/filter/:id",
		app.findFilterById(
			app.readSimpleAction(
				http.HandlerFunc(app.filterAction))))

	// Handler for scheduler
	router.Handler(http.MethodGet, "/api/schedule", http.HandlerFunc(app.getSchedules))
	router.Handler(http.MethodPut, "/api/schedule", http.HandlerFunc(app.addSchedule))
	router.Handler(http.MethodDelete, "/api/schedule/:id", http.HandlerFunc(app.deleteSchedule))

	router.Handler(http.MethodGet, "/api/fakeap", http.HandlerFunc(app.getFakeAP))

	// Use the http.FileServer handler to serve the static files from the ./ui/dist/ directory.
	fileServer := http.FileServer(http.Dir("/usr/local/raspberry/ui/dist/"))
	router.Handler(http.MethodGet, "/assets/*filepath", fileServer)

	// Use the http.FileServer handler to serve the capture files from the /usr/local/raspberry/captures/
	captureServer := http.FileServer(http.Dir("/usr/local/raspberry/"))
	router.Handler(http.MethodGet, "/captures/*filepath", captureServer)

	// Return the router as the http.Handler
	return router
}
