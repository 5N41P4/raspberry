package main

import (
	"net/http"
	"testing"
)

func TestPing(t *testing.T) {
	// Create an instance of our application struct which holds the application-wide dependencies.
	app := newTestApplication(t)

	// Create a new testServer instance passing in the routes from our application.
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	// Make a GET request to the /ping endpoint on the test server.
	code, _, body := ts.get(t, "/ping")

	// Check that the response code is 200 OK.
	if code != http.StatusOK {
		t.Errorf("want %d; got %d", http.StatusOK, code)
	}

	// Check that the response body is what we expect.
	if string(body) != "OK" {
		t.Errorf("want body to equal %q", "OK")
	}
}

func TestHome(t *testing.T) {
	// Create an instance of our application struct which holds the application-wide dependencies.
	app := newTestApplication(t)

	// Create a new testServer instance passing in the routes from our application.
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	// Make a GET request to the / endpoint on the test server.
	code, _, _ := ts.get(t, "/")

	// Check that the response code is 200 OK.
	if code != http.StatusOK {
		t.Errorf("want %d; got %d", http.StatusOK, code)
	}

	// Make a GET request to the /unknown endpoint on the test server.
	code, _, _ = ts.get(t, "/unknown")

	// Check that the response code is 404 Not Found.
	if code != http.StatusNotFound {
		t.Errorf("want %d; got %d", http.StatusNotFound, code)
	}
}

func TestAPITest(t *testing.T) {
	// Create an instance of our application struct which holds the application-wide dependencies.
	app := newTestApplication(t)

	// Create a new testServer instance passing in the routes from our application.
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	// Make a GET request to the /api/test endpoint on the test server.
	code, _, body := ts.get(t, "/api/test")

	// Check that the response code is 200 OK.
	if code != http.StatusOK {
		t.Errorf("want %d; got %d", http.StatusOK, code)
	}

	// Check that the response body is what we expect.
	expected := "{\"labels\":[\"Data 1\",\"Data 2\",\"Data 3\"],\"datasets\":[{\"label\":\"Placeholder Data\",\"data\":[10,20,30],\"backgroundColor\":[\"rgba(255, 99, 132, 0.3)\",\"rgba(54, 162, 235, 0.3)\",\"rgba(255, 206, 86, 0.3)\"],\"borderColor\":\"white\",\"borderWidth\":1}]}\n"

	if string(body) != expected {
		t.Errorf("want body to equal %q, got: %q", expected, body)
	}
}

func TestDiskUsage(t *testing.T) {
	// Create an instance of our application struct which holds the application-wide dependencies.
	app := newTestApplication(t)

	// Create a new testServer instance passing in the routes from our application.
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	// Make a GET request to the /api/disk endpoint on the test server.
	code, _, _ := ts.get(t, "/api/system")

	// Check that the response code is 200 OK.
	if code != http.StatusOK {
		t.Errorf("want %d; got %d", http.StatusOK, code)
	}
}
