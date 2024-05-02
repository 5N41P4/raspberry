package main

import (
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Create newTestApplication helper which returns an instance of our application struct containing mocked dependencies.
func newTestApplication(t *testing.T) *application {
	return &application{
		errorLog: log.New(io.Discard, "", 0),
		infoLog:  log.New(io.Discard, "", 0),
	}
}

// Define a custom testServer type which anonymously embeds a httptest.Server instance
type testServer struct {
	*httptest.Server
}

// Create a newTestServer helper which initalizes and returns an instance of our custom testServer type.
func newTestServer(t *testing.T, h http.Handler) *testServer {
	ts := httptest.NewTLSServer(h)

	return &testServer{ts}
}

// Implement a get helper which makes a GET request to the test server and returns the response.
func (ts *testServer) get(t *testing.T, path string) (int, http.Header, []byte) {
	rs, err := ts.Client().Get(ts.URL + path)
	if err != nil {
		t.Fatal(err)
	}

	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}

	return rs.StatusCode, rs.Header, body
}