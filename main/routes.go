package main

import "net/http"

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/{.}",
		MakeHandler,
	},
	Route{
		"Index",
		"GET",
		"/{.}/{.}",
		MakeHandler,
	},
	Route{
		"Index",
		"GET",
		"/",
		MakeHandler,
	},
	Route{
		"Index",
		"GET",
		"/{.}/{.}/{.}",
		MakeHandler,
	},
	/*	Route{
			"plugin1",
			"GET",
			"/plugin1",
			MakeHandler,
		},
		Route{
			"plugin1",
			"GET",
			"/plugin1/{.}",
			MakeHandler,
		},
		Route{
			"plugin1",
			"GET",
			"/plugin1/{.}/{.}",
			MakeHandler,
		},*/
}
