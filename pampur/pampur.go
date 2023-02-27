package pampur

import (
	"net/http"
	"regexp"
	_"fmt"
	"github.com/pampur/router"
)

type Pampur struct {
	rtrs 	 [] *router.Router; // list of routers
	handlers []	 router.Handler; // list of middleware functions
}

func (p *Pampur) Use(h router.Handler) {
	p.handlers = append(p.handlers, h)
}

func (p *Pampur) CreateRouter(basePath string) *router.Router {
	basePath = "^" + basePath
	router := &router.Router{Bp: regexp.MustCompile(basePath)}
	p.rtrs = append(p.rtrs, router)
	return router
}

func (p *Pampur) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// first get the path from the request
	path := r.URL.Path;

	// now loop through all the routers and match the base path
	for _, rt := range p.rtrs {
		match := rt.Bp.FindString(path);
		if match == "" {
			// no base path matches
			return
		}

		route := rt.FindRoute(path, r.Method)
		if route == nil {
			// no subpath matches
			return
		}
	
		// get the parameters
		params := getParams(r.URL.Path, route.Pattern)

		// set the params inside the context object
		ctx := router.Ctx{ Params: params }

		p.runStack(w, r, route, &ctx)
	}
}

func getParams(path string, pattern *regexp.Regexp) map[string]any {
	matches := pattern.FindStringSubmatch(path)  // values
	if len(matches) > 1 {
		keys := pattern.SubexpNames()
		params := make(map[string]any)
		i := 1
		for i < len(matches) {
			params[keys[i]] = matches[i]
			i++
		}

		return params
	}
	return nil
}

func (p *Pampur) runStack(w http.ResponseWriter, req *http.Request, r *router.Route, ctx *router.Ctx) {
	var next router.NextFunction
	i := 0
	var stack []router.Handler
	if (len(p.handlers) > 0) {
		stack = p.handlers
	} else {
		stack = r.Methods[req.Method]
	}

	done := false
	next = func() {
		if i == len(stack) {
			if (done) {
				return
			}
			stack = r.Methods[req.Method]
			i = 0
			done = true
		}

		// call the current handler
		currHandler := stack[i]
		i++

		currHandler(ctx, w, req, next)
		// increment the counter by one
	}
	next()
}