package models

import (
	"fmt"
	"log"
	"time"

	rpio "github.com/stianeikeland/go-rpio"
)

// GPIOCheck checks if there is GPIO access.
func GPIOCheck() bool {
	err := rpio.Open()
	if err != nil {
		fmt.Println("Cannot access GPIO memory range!", err)
		return false
	}
	return true
}

// RunEffectGPIO runs the effect in GPIO, toggles the pins according to the duration.
func RunEffectGPIO(effect Effect, length time.Duration) {
	for _, pin := range effect.Pins {
		go func(pin int) {
			pinOn(pin)
			time.Sleep(time.Millisecond * length)
			pinOff(pin)
		}(pin)
	}
}

func pinOff(pin int) {
	if !isRPI {
		fmt.Println("Simulation: Pin ", pin, " set to off")
		return
	}
	rpin, ok := loadedPins[pin]
	if !ok {
		log.Fatal("No pin ", pin, " in map!")
		return
	}
	if Conf.ActiveHigh {
		rpin.High()
	} else {
		rpin.Low()
	}
}

func pinOn(pin int) {
	if !isRPI {
		fmt.Println("Simulation: Pin ", pin, " set to on")
		return
	}
	rpin, ok := loadedPins[pin]
	if !ok {
		log.Fatal("No pin ", pin, " in map!")
		return
	}
	if Conf.ActiveHigh {
		rpin.Low()
	} else {
		rpin.High()
	}
}
