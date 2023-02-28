package pampur

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	//"github.com/pampur/err"
)

type NextFunction 	func()
type Handler func(c *Ctx, w http.ResponseWriter, r *http.Request, n NextFunction) Error

type Ctx struct {
	Params  map[string]any
}

type Route struct {
	Pattern 	*regexp.Regexp 		// pattern to match against the url requested
	Methods 	map[string][]Handler
}

type Router struct {
	Rts 		[]*Route
	Handlers 	[]Handler 
	Bp			*regexp.Regexp	
}

func (r *Router) Get(path string, handlers ...Handler) {
	// get the parameters if any
	// chunk the string in pieces
	chunks := strings.Split(path, "/")	

	for i, chunk := range chunks {
		if i > 0 {
			if chunk[0] == ':' {
				// add named capture group
				chunks[i] = "(?P<" + chunk[1:] + ">" + "[0-9a-zA-Z]+)"
			}
		}
	}

	// reconstruct the string and compile the regex
	pattern := regexp.MustCompile(strings.Join(chunks, "/"))

	// check a route with Get method does not already exist
	for _, route := range r.Rts {
		if route.Pattern.String() == pattern.String() {

			// check if the method is already registered
			_, ok := route.Methods["GET"]
			if ok {
				return
			} else {
				route.Methods["GET"] = handlers
			}
		}
	}

	// create a new method map
	methMap := make(map[string][]Handler)
	methMap["GET"] = handlers

	// add a new route to the list
	newRoute := &Route{
		Methods: methMap,
		Pattern: pattern,	
	}

	r.Rts = append(r.Rts, newRoute) 
}

func (r *Router) Use(handlers ...Handler) {
	r.Handlers = append(r.Handlers, handlers...)
}

func (r *Router) Print() {
	for _, route := range r.Rts {
		fmt.Println(route.Pattern.String())
	}
}

func (r *Router) FindRoute(path string, method string) *Route {
	// find the first match
	for _,  route:= range r.Rts {
		match := route.Pattern.FindString(path)
		if match == "" {
			return nil;
		}

		// route exists	
		// now check if the method is supported
		_, ok := route.Methods[method]
		if ok == false {
			return nil
		}

		return route;
	}
	return nil
}