package main

import (
	"net/http"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"AddSourceCode",
		"POST",
		"/add_source_code",
		AddSourceCode,
	},
	Route{
		"CheckSourceCode",
		"POST",
		"/check_source_code",
		CheckSourceCode,
	},
	Route{
		"RunSourceCode",
		"POST",
		"/run_source_code",
		RunSourceCode,
	},
	Route{
		"CompareSourceCode",
		"POST",
		"/compare_source_code",
		CompareSourceCode,
	},
}
