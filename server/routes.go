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
		"AddProgram",
		"POST",
		"/add_program",
		AddProgram,
	},
	Route{
		"CheckProgram",
		"POST",
		"/check_program",
		CheckProgram,
	},
	Route{
		"RunProgram",
		"POST",
		"/run_program",
		RunProgram,
	},
	Route{
		"CompareProgram",
		"POST",
		"/compare_program",
		CompareProgram,
	},
}
