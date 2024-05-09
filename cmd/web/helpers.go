package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/5N41P4/raspberry/cmd/modules"
)

// Define a writeJSON helper to write a JSON response to the client. This takes
// the same parameters as http.Error() and writes a JSON representation of the
// provided data to the response body, setting the Content-Type header to application/json.
func (app *application) writeJSON(w http.ResponseWriter, status int, data interface{}, headers http.Header) error {
	// Encode the data to JSON, returning the error if there was one.
	js, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Append a newLine character to the response body, then write the status code and
	// content-type header to the response header map.
	js = append(js, '\n')

	// Append the headers to the response header map.
	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}

// Implementation of a JSON reader that handles all possible JSON errors for the API
func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	// Use http.MaxbytesReader() to limit the size to 1MB
	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	// Initialize the JSON decoder and configure it to return an error on wrong fields
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	//Decode the request body to the destination
	err := dec.Decode(dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshallTypeError *json.UnsupportedTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError

		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)

		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")

		case errors.As(err, &unmarshallTypeError):
			if unmarshallTypeError.Error() != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshallTypeError.Error())
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshallTypeError.Type.Align())

		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("body contains unknown key %s", fieldName)

		case err.Error() == "http: request body too large":
			return fmt.Errorf("body must not be larger than %d bytes", maxBytes)

		case errors.As(err, &invalidUnmarshalError):
			panic(err)

		default:
			return err
		}
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must only contain a single JSON value")
	}

	return nil
}

func (app *application) refreshLists() {
	// Start time for logging
	start := time.Now()

	aps, cls, err := modules.ParseCSV("/usr/local/raspberry/recon/discovery-01.csv")

	if err != nil {
		return
	}

	// Add the access points to the map
	for _, ap := range aps {
		if app.filters.ApFilter.IsAllowed(ap.Bssid) {
			app.access[ap.Bssid] = &ap
		}
	}

	// Add the clients to the map
	for _, cli := range cls {
		if app.filters.ClFilter.IsAllowed(cli.Station) {
			app.clients[cli.Station] = &cli
		}
	}

	app.infoLog.Printf("Parsing took: %v\n", time.Since(start))
}
