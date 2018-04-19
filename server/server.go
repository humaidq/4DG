package server

import (
	"io/ioutil"
	"log"

	"github.com/humaidq/4DG/movie"
	macaron "gopkg.in/macaron.v1"
)

var (
	// loadedMovies An array of FDMovies
	loadedMovies [128]movie.FDMovie
)

// initialize Loads the movies and configuration from file.
func initialize() {
	scriptsDir, err := ioutil.ReadDir("./scripts")
	if err != nil {
		log.Fatal(err)
	}

	var i int
	for _, f := range scriptsDir {
		mov, err := movie.Decode("./scripts/" + f.Name())
		if err != nil {
			log.Fatal("Unable to "+f.Name()+": ", err)
		}
		loadedMovies[i] = mov
		i++
	}
}

// RunServer runs the macaron HTTP server
func RunServer() {
	initialize()
	// Setup Macaron
	macaron.Env = macaron.PROD
	m := macaron.Classic()
	m.Use(macaron.Renderer())

	// Define Routes
	m.Get("/", func(ctx *macaron.Context) {
		ctx.Data["Title"] = "4DG Dashboard"
		ctx.Data["Movies"] = loadedMovies
		ctx.HTML(200, "index")
	})

	m.Get("/play/:movie", func(ctx *macaron.Context) {
		ctx.Data["Title"] = "Run movie script - 4DG"
		ctx.HTML(200, "play")
	})

	// Run Macaron HTTP Server
	m.Run("0.0.0.0")
}
