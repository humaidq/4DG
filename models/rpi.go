package models

import (
	"fmt"
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
			rpin := rpio.Pin(pin)
			pinOn(rpin)
			time.Sleep(time.Millisecond * length)
			pinOff(rpin)
		}(pin)
	}
}

func pinOff(pin rpio.Pin) {
	if !isRPI {
		fmt.Println("Simulation: Pin ", pin, " set to off")
		return
	}
	if Conf.ActiveHigh {
		pin.High()
	} else {
		pin.Low()
	}
}

func pinOn(pin rpio.Pin) {
	if !isRPI {
		fmt.Println("Simulation: Pin ", pin, " set to on")
		return
	}
	if Conf.ActiveHigh {
		pin.Low()
	} else {
		pin.High()
	}
}
