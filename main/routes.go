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
		"/",
		Index,
	},
	Route{
		"Index",
		"GET",
		`/{..*}`,
		GenericHandler,
	},
	Route{
		"Index",
		"GET",
		`/{..*}/{..*}`,
		GenericHandler,
	},
	Route{
		"Index",
		"GET",
		`/{..*}/{..*}/{..*}`,
		GenericHandler,
	},
	Route{
		"Index",
		"POST",
		`/{..*}`,
		GenericHandler,
	},
	Route{
		"Index",
		"POST",
		`/{..*}/{..*}`,
		GenericHandler,
	},
	Route{
		"Index",
		"POST",
		`/{..*}/{..*}/{..*}`,
		GenericHandler,
	},
	Route{
		"Index",
		"DELETE",
		`/{..*}`,
		GenericHandler,
	},
	Route{
		"Index",
		"DELETE",
		`/{..*}/{..*}`,
		GenericHandler,
	},
	Route{
		"Index",
		"DELETE",
		`/{..*}/{..*}/{..*}`,
		GenericHandler,
	},
	Route{
		"Index",
		"PUT",
		`/{..*}`,
		GenericHandler,
	},
	Route{
		"Index",
		"PUT",
		`/{..*}/{..*}`,
		GenericHandler,
	},
	Route{
		"Index",
		"PUT",
		`/{..*}/{..*}/{..*}`,
		GenericHandler,
	},
}
