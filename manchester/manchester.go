package manchester

import (
	"fmt"
	"log"
	"time"

	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/host"
)

func init() {
	// Load all the drivers:
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}

}

// Pin defines the hardware pin which we are sending bits on
var Pin interface {
	FastOut(gpio.Level)
	Out(gpio.Level) error
}

var state gpio.Level

// ShipPayload sends what i believe is Manchester encoding
// at about 1200 baud - i just used my hands on these functions
// until their output looked like a reference trace on a scope :)
func ShipPayload(input []byte, repeat int) error {

	// ensure we start off low
	state = gpio.Low

	// Out pulls the pin low - but it also ensures this
	// pin is configured as an output
	err := Pin.Out(gpio.Low)
	if err != nil {
		return fmt.Errorf("could not configure and pull pin low: %s", err)
	}

	// 430us is the lowest possible pulse we need
	t := time.NewTicker(430 * time.Microsecond)

	// every time we start shipping a payload
	// we should start off by sending the 85 byte
	// also known as 01010101. I dont know why but it
	// could have something todo with giving the remote end
	// a chance of locking to our clock
	sendByte(byte(85), t.C)

	for index := 0; index < repeat; index++ {
		for _, b := range input {
			sendByte(b, t.C)
		}

		// wait between repeats must stay low
		Pin.FastOut(gpio.Low)
		state = gpio.Low

		// wait 3 miliseconds between package repeats
		// do this by skipping 7 "clock cycles"
		// as this should be around 3ms
		for skip := 0; skip < 7; skip++ {
			<-t.C
		}
	}

	<-t.C
	t.Stop()

	// we should not return from this function unless
	// we are in low state
	Pin.FastOut(gpio.Low)
	state = gpio.Low

	return nil
}

func sendByte(input byte, clock <-chan time.Time) {
	var bit byte
	for index := 0; index < 8; index++ {

		// read LSB
		bit = input & 1

		// move the rest over
		input = input >> 1

		if bit == 1 {
			// high bits must transition High -> Low
			// if we are low - we must go high first
			if !state {
				Pin.FastOut(gpio.High)
				state = gpio.High
			}
			<-clock
			Pin.FastOut(gpio.Low)
			state = gpio.Low
			<-clock
			continue
		}

		// low bits must transition Low -> High
		// flip state if we are High
		if state {
			Pin.FastOut(gpio.Low)
			state = gpio.Low
		}
		<-clock
		Pin.FastOut(gpio.High)
		state = gpio.High
		<-clock
	}
}
