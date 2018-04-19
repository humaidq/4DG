package models

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Effects []Effect `toml:"effects_labels"`
}

// Effect is used to give pins a user-friendly name.
type Effect struct {
	EffectName string `toml:"effect_label"`
	Pins       []int  `toml:"pins"`
}

// Conf A loaded config struct from file.
var Conf Config
var isRPI bool

// LoadConfig Loads configuration from file and stores it into variable Conf
func LoadConfig() {
	b, err := ioutil.ReadFile("config.toml")
	if err != nil {
		log.Fatal(err)
	}
	if _, err := toml.Decode(string(b), &Conf); err != nil {
		log.Fatal(err)
	}
	if isRPI := GPIOCheck(); !isRPI {
		fmt.Println("Will only run simulation mode.")
	}
}
