package models

import (
	"fmt"
	"io/ioutil"
	"log"
	"time"

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
		//RunMovie(mov)
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

var runningMovie bool
var currentMovieTime MovieTime

// MovieTime represents the time since the movie started.
type MovieTime struct {
	sec int
	dec int
}

func GetMovie(movieName string) (FDMovie, bool) {
	for _, movie := range LoadedMovies {
		if movie.MovieName == movieName {
			return movie, true
		}
	}
	return FDMovie{}, false
}

// RunMovie will run the specified movie as long as there is no other movie running.
func RunMovie(movie FDMovie) {
	if !runningMovie {
		runningMovie = true
		currentMovieTime = MovieTime{0, 0}
		for {
			fmt.Println(formatMovieTime())
			effect, ok := movie.Effects[formatMovieTime()]
			if ok {
				for _, confEffect := range Conf.Effects {
					if confEffect.EffectName == effect.EffectName {
						RunEffectGPIO(confEffect, time.Duration(effect.EffectLength))
					}
				}
			}
			incrementTime()
			time.Sleep(100 * time.Millisecond)
		}

	}
}

func incrementTime() {
	if currentMovieTime.dec < 9 {
		currentMovieTime.dec++
	} else {
		currentMovieTime.dec = 0
		currentMovieTime.sec++
	}
}

// MillisecondsToPosition converts Milliseconds to position format.
func formatMovieTime() string {
	return fmt.Sprint(currentMovieTime.sec) + "." + fmt.Sprint(currentMovieTime.dec)
}

func getCurrentMS() int64 {
	return time.Now().UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond))
}
