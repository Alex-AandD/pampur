package main

import (
	_"fmt"
	"log"
	"net/http"

	"github.com/pampur/pampur"
	"github.com/pampur/router"
)

func main() {
	p := pampur.Pampur{}
	rtr := p.CreateRouter("/auth")
	rtr.Get("/hello/:id", 
	func(ctx *router.Ctx, w http.ResponseWriter, r *http.Request, n router.NextFunction) {
		w.Write([]byte("first"))
		n()
	},
	
	func(c *router.Ctx, w http.ResponseWriter, r *http.Request, n router.NextFunction) {
		w.Write([]byte("second"))
	})
	rtr.Print()
	log.Fatal(http.ListenAndServe(":8080", &p))
}