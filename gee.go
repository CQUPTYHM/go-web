package gee

import (
	"log"
	"net/http"
	"strings"
)

type (
	Engine struct {
		*RouterGroup
		router *router
		groups []*RouterGroup
	}

	RouterGroup struct {
		prefix      string
		middlewares []HandlerFunc
		engine      Engine
	}
)

func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: *engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
}
func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine

	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

func (group *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
	pattern := group.prefix + comp
	log.Printf("Route %4s - %s", method, pattern)
	group.engine.router.addRoute(method, pattern, handler)
}

// GET defines the method to add GET request
func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRoute("GET", pattern, handler)
}

// POST defines the method to add POST request
func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
}

// Run defines the method to start a http server
func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

func (engine *Engine) searchMiddlewares(URL string) []HandlerFunc {
	var middlewares []HandlerFunc
	for _, group := range engine.groups {
		if strings.HasPrefix(URL, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	return middlewares
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	middlewares := engine.searchMiddlewares(req.URL.Path)
	c := newContext(w, req)
	c.handlers = middlewares
	engine.router.handle(c)
}
