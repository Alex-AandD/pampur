package pampur

import (
	"net/http"
	"regexp"
	_"fmt"
	//_"github.com/pampur/router"
)

type Pampur struct {
	rtrs 	 [] *Router; // list of routers
	handlers []	 Handler; // list of middleware functions
}

func (p *Pampur) Use(h Handler) {
	p.handlers = append(p.handlers, h)
}

func (p *Pampur) CreateRouter(basePath string) *Router {
	basePath = "^" + basePath
	router := &Router{Bp: regexp.MustCompile(basePath)}
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
			handleError(w, NewHttpError(http.StatusText(http.StatusNotFound), http.StatusNotFound))
		}

		route := rt.FindRoute(path, r.Method)
		if route == nil {
			handleError(w, NewHttpError(http.StatusNotFound, http.StatusText(http.StatusNotFound)))
		}
	
		// get the parameters
		params := getParams(r.URL.Path, route.Pattern)

		// set the params inside the context object
		ctx := Ctx{ Params: params }

		err := p.runStack(w, r, route, &ctx)
		if err != nil {
			handleError(w, err)
		}
	}
}

func handleError(w http.ResponseWriter, err Error) {
	switch e := err.(type) {
		case HttpError:  {
			http.Error(w, e.Error(), e.Status())
		}
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

func (p *Pampur) runStack(w http.ResponseWriter, req *http.Request, r *Route, ctx *Ctx) Error {
	var next NextFunction
	i := 0
	done := false
	var stack []Handler

	if (len(p.handlers) > 0) {
		stack = p.handlers
	} else {
		stack = r.Methods[req.Method]
		done = true
	}

	var finalError Error
	next = func() {
		if i == len(stack) - 1 && done {
			err := stack[i](ctx, w, req, func() { })
			if err != nil {
				finalError = err
			}
		}

		if i >= len(stack) && !done {
			stack = r.Methods[req.Method]
			i = 0
			done = true
		}


		// call the current handler
		currHandler := stack[i]
		i++

		err := currHandler(ctx, w, req, next)
		if err != nil {
			finalError = err
		}
	}
	next()
	return finalError
}