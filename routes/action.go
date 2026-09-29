package routes

import (
	"net/http"

	"github.com/RandySteven/onboard-be/enums"
)

func registerEndpointRouter(methodName, path, method string, handler HandlerFunc, middlewares ...enums.Middleware) *Router {
	return &Router{methodName: methodName, path: path, handler: handler, method: method, middlewares: middlewares}
}

func Post(methodName, path string, handler HandlerFunc, middlewares ...enums.Middleware) *Router {
	return registerEndpointRouter(methodName, path, http.MethodPost, handler, middlewares...)
}

func Get(methodName, path string, handler HandlerFunc, middlewares ...enums.Middleware) *Router {
	return registerEndpointRouter(methodName, path, http.MethodGet, handler, middlewares...)
}
