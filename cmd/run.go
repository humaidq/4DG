package cmd

import (
	"sort"
	"strconv"

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
		ctx.Data["Simulation"] = models.GPIOCheck
		ctx.HTML(200, "index")
	})

	m.Get("/play/:movie", func(ctx *macaron.Context) {
		ctx.Data["Title"] = "Run movie script - 4DG"
		mov, ok := models.GetMovie(ctx.Params("movie"))
		if !ok {
			ctx.Redirect("/404")
			return
		}
		go models.RunMovie(mov)

		ctx.Data["Movie"] = mov

		ctx.HTML(200, "play")
	})

	m.Get("/timer", func(ctx *macaron.Context) {
		ctx.Data["Timer"] = models.FormatMovieTime()
		ctx.HTML(200, "timer")
	})

	m.Get("/adj/:offset", func(ctx *macaron.Context) {

		ctx.Redirect("/404")
	})

	m.Get("/edit/:movie", func(ctx *macaron.Context) {
		ctx.Data["Title"] = "Edit " + ctx.Params("movie") + " - 4DG"

		mov, ok := models.GetMovie(ctx.Params("movie"))
		if !ok {
			ctx.Redirect("/404")
			return
		}

		type sortedTimestamps struct {
			Index     string
			Timestamp *models.TimestampEffect
		}

		var sorted []sortedTimestamps
		for k, v := range mov.Effects {
			sorted = append(sorted, sortedTimestamps{k, v})
		}

		sort.Slice(sorted, func(i, j int) bool {
			i1, _ := strconv.ParseFloat(sorted[i].Index, 64)
			i2, _ := strconv.ParseFloat(sorted[j].Index, 64)
			return i1 < i2
		})

		ctx.Data["Sorted"] = sorted

		ctx.Data["Movie"] = mov
		ctx.HTML(200, "editmov")
	})

	m.Get("/edit/:movie/:pos", func(ctx *macaron.Context) {
		ctx.Data["Title"] = "Edit " + ctx.Params("movie") + " - 4DG"
		mov, ok := models.GetMovie(ctx.Params("movie"))
		if !ok {
			ctx.Redirect("/404")
			return
		}
		ctx.Data["Movie"] = mov
		eff, ok := mov.Effects[ctx.Params("pos")]
		ctx.Data["Effect"] = eff
		ctx.Data["Pos"] = ctx.Params("pos")

		ctx.HTML(200, "editpos")
	})

	m.Post("/edit/:movie/:pos", func(ctx *macaron.Context) {

		mov, ok := models.GetMovie(ctx.Params("movie"))
		if !ok {
			ctx.Redirect("/404")
			return
		}
		len, err := strconv.Atoi(ctx.Query("length"))
		if err != nil {
			return
		}

		mov.Effects[ctx.Params("pos")].EffectLength = len
		mov.Effects[ctx.Params("pos")].EffectName = ctx.Query("effect")

		models.SaveMovie(ctx.Params("movie"))

		ctx.Redirect("/edit/" + ctx.Params("movie"))
	})

	m.Get("/edit/:movie/delete/:pos", func(ctx *macaron.Context) {

		mov, ok := models.GetMovie(ctx.Params("movie"))
		if !ok {
			ctx.Redirect("/404")
			return
		}

		delete(mov.Effects, ctx.Params("pos"))

		models.SaveMovie(ctx.Params("movie"))

		ctx.Redirect("/edit/" + ctx.Params("movie"))
	})

	m.Get("/edit/:movie/new", func(ctx *macaron.Context) {
		ctx.Data["Title"] = "New effect for " + ctx.Params("movie") + " - 4DG"
		mov, ok := models.GetMovie(ctx.Params("movie"))
		if !ok {
			ctx.Redirect("/404")
			return
		}
		ctx.Data["Movie"] = mov

		ctx.HTML(200, "newpos")
	})

	m.Post("/edit/:movie/new", func(ctx *macaron.Context) {
		mov, ok := models.GetMovie(ctx.Params("movie"))
		if !ok {
			ctx.Redirect("/404")
			return
		}
		len, err := strconv.Atoi(ctx.Query("length"))
		if err != nil {
			return
		}

		ts := models.TimestampEffect{EffectName: ctx.Query("effect"), EffectLength: len}
		mov.Effects[ctx.Query("pos")] = &ts

		models.SaveMovie(ctx.Params("movie"))

		ctx.Redirect("/edit/" + ctx.Params("movie"))
	})

	m.Get("/new", func(ctx *macaron.Context) {
		ctx.Data["Title"] = "New movie script - 4DG"
		ctx.HTML(200, "new")
	})

	m.Post("/new", func(ctx *macaron.Context) {
		i, err := strconv.Atoi(ctx.Query("length"))
		if err != nil {
			return
		}
		newMovie := models.FDMovie{MovieName: ctx.Query("name"), MovieLength: i}
		newMovie.Effects = make(map[string]*models.TimestampEffect)

		lMovie := models.LoadedMovie{Movie: newMovie, Filename: ctx.Query("file") + ".toml"}
		models.LoadedMovies[models.LoadedMoviesSize] = lMovie
		models.LoadedMoviesSize++
		models.SaveMovie(ctx.Query("name"))
		ctx.Redirect("/edit/" + ctx.Query("name"))
	})

	// Run Macaron HTTP Server
	m.Run("0.0.0.0")
}

func init() {
	rootCmd.AddCommand(runCmd)

}
