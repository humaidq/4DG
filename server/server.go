package server

import macaron "gopkg.in/macaron.v1"

// RunServer runs the macaron HTTP server
func RunServer() {
	// Setup Macaron
	macaron.Env = macaron.PROD
	m := macaron.Classic()
	m.Use(macaron.Renderer())

	// Define Routes
	m.Get("/", func(ctx *macaron.Context) {
		ctx.Data["Title"] = "4DG Dashboard"
		ctx.HTML(200, "index")
	})

	m.Get("/", func(ctx *macaron.Context) {
		ctx.Data["Title"] = "New Movie Script - 4DG"
		ctx.HTML(200, "new")
	})

	// Run Macaron HTTP Server
	m.Run("0.0.0.0")
}
