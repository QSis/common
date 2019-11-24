package server

import (
	ginv1 "github.com/gin-gonic/gin"
)

type (
	RouterFunc func(string, ...ginv1.HandlerFunc) ginv1.IRoutes
	HandleFunc func(*Context)

	Router struct {
		Uri    string
		Method RouterFunc
		Action HandleFunc
	}
)

func NewRouter(uri string, method RouterFunc, action HandleFunc) *Router {
	return &Router{
		Uri:    uri,
		Method: method,
		Action: action,
	}
}
