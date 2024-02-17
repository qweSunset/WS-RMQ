package httpserverhelper

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
	Middleware  func(next http.HandlerFunc) http.HandlerFunc
}

type Routes []Route

func logger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)

		logStr := fmt.Sprintf("%s\t%s\t%s\t%s", r.Method, r.RequestURI, name, time.Since(start))

		log.Println(logStr)
	})
}

func AccessControlMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func NewRouter(routes Routes) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	router.Use(AccessControlMiddleware)

	for _, route := range routes {
		var handler http.Handler

		if route.Middleware != nil {
			handler = route.Middleware(route.HandlerFunc)
		} else {
			handler = route.HandlerFunc
		}

		handler = logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(cors.Default().Handler(handler))
	}
	return router
}
