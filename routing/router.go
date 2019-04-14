package routing

import (
	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {

	router := mux.NewRouter().StrictSlash(true)

	//following this -> router.Methods("GET").Path("/hello").Name("hello").HandlerFunc(Hello)
	//router.Methods("GET").Path("/hello").Name("hello").HandlerFunc(Hello)

	for _, route := range routes {

		router.
			Methods(route.Method).
			Path(route.Path).
			Name(route.Name).
			Handler(Logger(route.HandlerFunc, route.Name))

	}

	return router
}
