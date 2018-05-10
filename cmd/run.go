package cmd

import (
	"github.com/humaidq/4DG/models"
	"github.com/spf13/cobra"
	macaron "gopkg.in/macaron.v1"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the 4DG web server",
	Run: func(cmd *cobra.Command, args []string) {
		runServer()
	},
}

// runServer runs the macaron HTTP server
func runServer() {
	models.Initialize()
	// Setup Macaron
	macaron.Env = macaron.PROD
	m := macaron.Classic()
	m.Use(macaron.Renderer())

	// Define Routes
	m.Get("/", func(ctx *macaron.Context) {
		ctx.Data["Title"] = "4DG Dashboard"
		ctx.Data["Movies"] = models.LoadedMovies
		ctx.HTML(200, "index")
	})

	m.Get("/play/:movie", func(ctx *macaron.Context) {
		ctx.Data["Title"] = "Run movie script - 4DG"
		ctx.HTML(200, "play")
	})

	m.Get("/edit/:movie", func(ctx *macaron.Context) {
		ctx.Data["Title"] = "Edit " + ctx.Params("movie") + " - 4DG"

		mov, ok := models.GetMovie(ctx.Params("movie"))
		if !ok {
			ctx.Redirect("/404")
			return
		}
		ctx.Data["Movie"] = mov
		ctx.HTML(200, "editmov")
	})

	m.Get("/edit/:movie/:", func(ctx *macaron.Context) {
		ctx.Data["Title"] = "Edit " + ctx.Params("movie") + " - 4DG"
		/*mov, ok := models.GetMovie(ctx.Params("movie"))
		if !ok {
			ctx.Redirect("/404")
			return
		}*/

		ctx.HTML(200, "editmov")
	})

	// Run Macaron HTTP Server
	m.Run("0.0.0.0")
}

func init() {
	rootCmd.AddCommand(runCmd)

}
