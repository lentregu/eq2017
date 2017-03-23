package main

import (
	"net/http"
	"time"

	"github.com/TDAF/gologops"
	"github.com/gorilla/mux"
)

// Logger is a handler to wrap other handlers and log basic request parameters
func loggerHandler(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		inner.ServeHTTP(w, r)
		gologops.Infof("%s    %s    %s    %s", r.Method, r.RequestURI, name, time.Since(start))
	})
}

func newRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = loggerHandler(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)

	}
	return router
}
