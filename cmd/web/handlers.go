package main

import (
	"context"
	"errors"
	"net/http"
	"os"

	"github.com/5N41P4/raspberry/internal/data"
	"github.com/julienschmidt/httprouter"
)

// Basic Server Functions

// home handles the HTTP request for the home page.
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	wd, err := os.Getwd()

	if err != nil {
		app.serverError(w, err)
		return
	}

	app.infoLog.Println(wd)

	http.ServeFile(w, r, "/usr/local/raspberry/ui/dist/index.html")
}

// Test Handlers

func (app *application) ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func (app *application) apiTest(w http.ResponseWriter, r *http.Request) {
	data := data.Data{
		Labels: []string{"Data 1", "Data 2", "Data 3"},
		Datasets: []data.Dataset{
			{
				Label:           "Placeholder Data",
				Data:            []int{10, 20, 30},
				BackgroundColor: []string{"rgba(255, 99, 132, 0.3)", "rgba(54, 162, 235, 0.3)", "rgba(255, 206, 86, 0.3)"},
				BorderColor:     "white",
				BorderWidth:     1,
			},
		},
	}

	w.Header().Set("Content-Type", "application/json")
	err := app.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		app.serverError(w, err)
	}
}

func (app *application) readJSONMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			if r.Header.Get("Content-Type") != "application/json" {
				app.clientError(w, http.StatusUnsupportedMediaType)
				return
			}
		}

		var intput data.ApiAction
		err := app.readJSON(w, r, &intput)
		if err != nil {
			app.serverError(w, err)
			return
		}
		ctx := context.WithValue(r.Context(), "input", &intput)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (app *application) findInterfaceByIdMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := httprouter.ParamsFromContext(r.Context()).ByName("id")
		inf, ok := app.interfaces[id]
		if !ok {
			app.errorLog.Println("[INF] Interface could not be found")
			app.serverError(w, errors.New("interface not found"))
			return
		}
		ctx := context.WithValue(r.Context(), "inf", inf)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (app *application) findFilterByIdMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := httprouter.ParamsFromContext(r.Context()).ByName("id")
		var list *Filter

		if id == "clients" {
			list = app.filters.ClFilter
		} else if id == "aps" {
			list = app.filters.ApFilter
		} else {
			app.clientError(w, http.StatusBadRequest)
		}
		ctx := context.WithValue(r.Context(), "list", list)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
