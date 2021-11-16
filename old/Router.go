package main

import (
	"net/http"
	"strings"
)

type Router struct {
	Routes 			map[string]*Service
	DefaultRoute 	string
}

func (r *Router) AddService(service *Service, routes ...string) {
	for _, route := range routes {
		r.Routes[route] = service
	}
}

func (r *Router) ServeDefaultRoute(Response http.ResponseWriter, Request *http.Request) {
	if r.DefaultRoute != "" {
		http.ServeFile(Response, Request, r.DefaultRoute)
	} else {
		Response.Write([]byte("Isn't this lonely?"))
	}
}

func (r * Router) ServeHTTP(Response http.ResponseWriter, Request *http.Request) {
	if service, ok := r.Routes[Request.URL.Path]; ok {
		for i := 0; i < len(service.MiddleServices); i++ {
			if !service.MiddleServices[i].ServeHTTP(Response, Request) {
				return;
			}
		}
		service.ServeHTTP(Response, Request)
		return
	}
	for route, service := range r.Routes {
		if (strings.HasPrefix(Request.URL.Path, route)){ 
			for i := 0; i < len(service.MiddleServices); i++ {
				if !service.MiddleServices[i].ServeHTTP(Response, Request) {
					return;
				}
			}
			service.ServeHTTP(Response, Request)
			return
		}
	}
	r.ServeDefaultRoute(Response, Request)
}

func NewRouter() Router {
	var router 		= Router{}
	router.Routes 	= make(map[string]*Service)
	return router
}