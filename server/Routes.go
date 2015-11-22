package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewRouter() *mux.Router {

	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
	// register handlers for these URIs. Kind of observer design pattern
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}

	return router
}

var routes = Routes{
	Route{
		"CreateLocation",
		"POST",
		"/locations",
		CreateLocation,
	},
	Route{
		"GetLocation",
		"GET",
		"/locations/{location_id}",
		GetLocation,
	},
	Route{
		"PutLocation",
		"PUT",
		"/locations/{location_id}",
		PutLocation,
	},
	Route{
		"DeleteLocation",
		"DELETE",
		"/locations/{location_id}",
		DeleteLocation,
	},
	Route{
		"CreateTripPlan",
		"POST",
		"/trips",
		CreateTripPlan,
	},
	Route{
		"GetTripPlan",
		"GET",
		"/trips/{trip_id}",
		GetTripPlan,
	},
	Route{
		"PutTripPlan",
		"PUT",
		"/trips/{trip_id}/request",
		PutTripPlan,
	},
}
