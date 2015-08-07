package main

import (
	"github.com/mdevilliers/take-home/cmd/server/app"
	"github.com/zenazn/goji"
)

func main() {

	// creates an instance of the api to serve
	api := app.NewApi()

	// sets up the default routes
	api.Route(goji.DefaultMux)

	goji.Serve()
}
