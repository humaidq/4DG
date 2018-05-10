package models

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/BurntSushi/toml"
)

// FDMovie holds all the information of a 4D movie script.
type FDMovie struct {
	MovieName   string                      `toml:"movie_name"`
	MovieLength int                         `toml:"movie_length"` // in minutes
	Effects     map[string]*TimestampEffect `toml:"effects"`
}

// TimestampEffect holds a specific effect at a specific timestamp.
type TimestampEffect struct {
	EffectName   string `toml:"effect_name"`
	EffectLength int    `toml:"length_ms"` // in milliseconds
}

// LoadedMovie holds a Movie struct and the name of the file it is loaded from.
type LoadedMovie struct {
	Movie    FDMovie
	Filename string
}

var (
	// LoadedMovies An array of FDMovies
	LoadedMovies     [128]LoadedMovie
	LoadedMoviesSize int // todo replace with function
)

// Initialize Loads the movies and configuration from file.
func Initialize() {
	LoadConfig()
	scriptsDir, err := ioutil.ReadDir("./scripts")
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range scriptsDir {
		mov, err := Decode("./scripts/" + f.Name())
		if err != nil {
			log.Fatal("Unable to "+f.Name()+": ", err)
		}
		if mov.Effects == nil {
			mov.Effects = make(map[string]*TimestampEffect)
		}
		LoadedMovies[LoadedMoviesSize] = LoadedMovie{mov, f.Name()}
		//RunMovie(mov)
		LoadedMoviesSize++
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

// GetMovie using it's name from the loaded movies map.
func GetMovie(movieName string) (FDMovie, bool) {
	for _, lmovie := range LoadedMovies {
		if lmovie.Movie.MovieName == movieName {
			return lmovie.Movie, true
		}
	}
	return FDMovie{}, false
}

// GetLoadedMovie using it's name from.
func GetLoadedMovie(movieName string) (LoadedMovie, bool) {
	for _, lmovie := range LoadedMovies {
		if lmovie.Movie.MovieName == movieName {
			return lmovie, true
		}
	}
	return LoadedMovie{}, false
}

// SaveMovie to scripts folder.
func SaveMovie(movieName string) {
	mov, ok := GetLoadedMovie(movieName)
	if !ok {
		return
	}

	buf := new(bytes.Buffer)
	if err := toml.NewEncoder(buf).Encode(mov.Movie); err != nil {
		log.Fatal(err)
	}
	ioutil.WriteFile("scripts/"+mov.Filename, buf.Bytes(), 0644)
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
	if currentMovieTime.dec == 0 {
		return fmt.Sprint(currentMovieTime.sec)
	}
	return fmt.Sprint(currentMovieTime.sec) + "." + fmt.Sprint(currentMovieTime.dec)
}

func getCurrentMS() int64 {
	return time.Now().UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond))
}
