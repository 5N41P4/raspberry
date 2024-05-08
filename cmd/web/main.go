package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/5N41P4/raspberry/cmd/modules"
	"github.com/5N41P4/raspberry/internal/data"
)

// Define an application struct to hold the application-wide dependencies for the web application.
type application struct {
	config     *Config
	errorLog   *log.Logger
	infoLog    *log.Logger
	inet       map[string]*modules.Interface
	interfaces map[string]*modules.Interface
	access     map[string]*data.AppAP
	clients    map[string]*data.AppClient
	filters    *FilterList
	updater    chan struct{}
}

// Refresh function for organizing lists and parsing files
func (app *application) refresh() {
	list := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-list.C:
			go app.refreshLists()

		case <-app.updater:
			list.Stop()
			return
		}
	}
}

// Cleanup function for graceful exit and saving of files
func cleanup(app *application) {
	for _, iface := range app.interfaces {
		iface.TryAction(data.ApiAction{
			Identifier: iface.Name,
			Action:     "stop",
		}, &app.access, &app.clients)
	}
	app.filters.cleanup()

	fmt.Println("Cleanup completed.")
}

func main() {

	// Define command-line flags for the HTTP server address.
	var config bool
	flag.BoolVar(&config, "config", false, "Run the configuration program to configure the interfaces, IP and port.")
	flag.Parse()

	if config {
		RunConfig()
	}

	// Read the config file, if it doesn't exist, go through setup.
	cfg, err := GetConfig()
	if err != nil {
		return
	}

	// Initialize a new logger which writes messages to the standard out stream,
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// Initialize a new instance of application containing the dependencies.
	app := &application{
		config:   cfg,
		infoLog:  infoLog,
		errorLog: errorLog,
		access:   make(map[string]*data.AppAP),
		clients:  make(map[string]*data.AppClient),
		filters:  newFilterList(),
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	go func() {
		// Get the network interfaces
		app.inet = modules.GetInterfaces(app.config.Inet)
		app.interfaces = modules.GetInterfaces(app.config.Interfaces)
		if app.interfaces == nil {
			errorLog.Fatal("No interfaces found")
		}

		// Initualize the updater in a go routine for all recurring tasks
		app.updater = make(chan struct{})
		go app.refresh()

		// Initialize a new http.Server struct.
		srv := &http.Server{
			Addr:     fmt.Sprintf("%s:%d", app.config.IP, app.config.Port),
			ErrorLog: errorLog,
			Handler:  app.routes(),
		}

		infoLog.Printf("Starting server on %s", srv.Addr)
		err := srv.ListenAndServe()
		errorLog.Fatal(err)
	}()

	// Wait for the main program or interrupt signal
	select {
	case <-interrupt:
		fmt.Println("Received interrupt signal. Shutting down...")
		close(app.updater)
		cleanup(app)
	}
}
