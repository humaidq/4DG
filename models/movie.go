package models

import (
	"io/ioutil"
	"log"

	"github.com/BurntSushi/toml"
)

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

var (
	// LoadedMovies An array of FDMovies
	LoadedMovies [128]FDMovie
)

// Initialize Loads the movies and configuration from file.
func Initialize() {
	LoadConfig()
	scriptsDir, err := ioutil.ReadDir("./scripts")
	if err != nil {
		log.Fatal(err)
	}

	var i int
	for _, f := range scriptsDir {
		mov, err := Decode("./scripts/" + f.Name())
		if err != nil {
			log.Fatal("Unable to "+f.Name()+": ", err)
		}
		LoadedMovies[i] = mov
		i++
	}
}

// Decode Converts a movie script file to an FDMovie struct
func Decode(fileName string) (FDMovie, error) {
	var mov FDMovie
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		return FDMovie{}, err
	}
	if _, err := toml.Decode(string(b), &mov); err != nil {
		return FDMovie{}, err
	}
	return mov, nil
}
