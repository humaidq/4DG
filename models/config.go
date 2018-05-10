package models

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/BurntSushi/toml"
	rpio "github.com/stianeikeland/go-rpio"
)

type Config struct {
	ActiveHigh bool
	Effects    []Effect `toml:"effects_labels"`
}

// Effect is used to give pins a user-friendly name.
type Effect struct {
	EffectName string `toml:"effect_label"`
	Pins       []int  `toml:"pins"`
}

// Conf A loaded config struct from file.
var Conf Config
var isRPI = true

// LoadConfig Loads configuration from file and stores it into variable Conf
func LoadConfig() {
	b, err := ioutil.ReadFile("config.toml")
	if err != nil {
		log.Fatal(err)
	}
	if _, err := toml.Decode(string(b), &Conf); err != nil {
		log.Fatal(err)
	}
	if GPIOCheck() {
		for _, effect := range Conf.Effects {
			for _, pin := range effect.Pins {
				rpin := rpio.Pin(pin)
				rpin.Output()
				pinOff(rpin)
			}
		}
	} else {
		fmt.Println("Will only run simulation mode.")
		isRPI = false
	}

}
