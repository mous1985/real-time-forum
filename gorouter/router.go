package gorouter

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

type Handler func(*Context)

type Route struct {
	Pattern *regexp.Regexp
	Handler Handler
	Method  string
	Keys    []string
}

type Router struct {
	Routes       []Route
	DefaultRoute Handler
}

func NewRouter() *Router {
	return &Router{}
}

func (r *Router) GET(pattern string, handler Handler) {
	r.handle(pattern, handler, http.MethodGet)
}

func (r *Router) POST(pattern string, handler Handler) {
	r.handle(pattern, handler, http.MethodPost)
}

func (r *Router) DELETE(pattern string, handler Handler) {
	r.handle(pattern, handler, http.MethodDelete)
}

func (r *Router) PUT(pattern string, handler Handler) {
	r.handle(pattern, handler, http.MethodPut)
}

func (r *Router) PATCH(pattern string, handler Handler) {
	r.handle(pattern, handler, http.MethodPatch)
}

func WrapHandler(h http.Handler) Handler {
	return func(ctx *Context) {
		h.ServeHTTP(ctx.ResponseWriter, ctx.Request)
	}
}

func (r *Router) handle(pattern string, handler Handler, method string) {
	regex, keys := readPatternAndKeys(pattern)
	route := Route{Pattern: regex, Handler: handler, Method: method, Keys: keys}

	r.Routes = append(r.Routes, route)
}

func readPatternAndKeys(pattern string) (*regexp.Regexp, []string) {
	var keys []string
	splited := strings.Split(pattern, "/")

	for i, v := range splited {
		if strings.HasPrefix(v, ":") {
			keys = append(keys, v[1:])
			splited[i] = `([\w\._-]+)`
		}
		if v == "*" {
			keys = append(keys, fmt.Sprintf("param%d", i))
			splited[i] = `([\w\._-]+)`
		}
	}

	regexStr := fmt.Sprintf("^%s$", strings.Join(splited, "/"))
	return regexp.MustCompile(regexStr), keys
}

func (r *Router) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	ctx := &Context{
		Request:        request,
		ResponseWriter: writer,
		Params:         make(map[string]string),
	}

	for _, rt := range r.Routes {
		if request.Method != rt.Method && request.Method != http.MethodOptions {
			continue
		}

		if matches := rt.Pattern.FindStringSubmatch(ctx.URL.Path); len(matches) > 0 {
			if len(matches) > 1 && len(rt.Keys) == len(matches[1:]) {
				ctx.setURLValues(rt.Keys, matches[1:])
			}
			rt.Handler(ctx)
			return
		}
	}

	ctx.WriteError(http.StatusNotFound, "404 Not Found")
}
