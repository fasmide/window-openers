package main

import (
	"github.com/fasmide/window-openers/manchester"
	"periph.io/x/periph/host/bcm283x"
)

func main() {
	// configure pin where the 433 radio is located
	manchester.Pin = bcm283x.GPIO21

	// Open window
	manchester.ShipPayload([]byte{177, 167, 68, 228, 115, 136, 77}, 12)

	// profit
}
