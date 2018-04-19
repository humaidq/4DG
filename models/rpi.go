package models

import (
	"fmt"

	rpio "github.com/stianeikeland/go-rpio"
)

func GPIOCheck() bool {
	err := rpio.Open()
	if err != nil {
		fmt.Println("Cannot access GPIO memory range!", err)
		return false
	}
	return true
}
