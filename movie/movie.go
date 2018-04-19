package movie

import "container/list"

// FDMovie holds all the information of a 4D movie script.
type FDMovie struct {
	MovieName   string                     `toml:"movie_name"`
	MovieLength int                        `toml:"movie_length"` // in minutes
	Effects     map[string]TimestampEffect `toml:"effects"`
}

// TimestampEffect holds a specific effect at a specific timestamp.
type TimestampEffect struct {
	EffectName   string `toml:"effect_name"`
	EffectLength int    `toml:"length_ms"` // in milliseconds
}

// Effect is used to give pins a user-friendly name.
type Effect struct {
	EffectName string
	pins       []int
}

var (
	// LoadedEffects A list of Effects
	LoadedEffects list.List
	// LoadedMovies A list of FDMovie
	LoadedMovies list.List
)

// Initialize Loads the movies and configuration from file.
func Initialize() {

}
