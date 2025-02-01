package main

import "strings"

type Router struct {
	routes map[string]Handler
}

func NewRouter() *Router {
	return &Router{
		routes: make(map[string]Handler),
	}
}

func (r *Router) AddRoute(path string, handler Handler) {
	r.routes[path] = handler
}

func (r *Router) FindHandler(path string) Handler {
	for prefix, handler := range r.routes {
		if strings.HasPrefix(path, prefix) {
			return handler
		}
	}

	return nil
}
