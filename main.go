package main

import (
	"log"
	"time"

	"github.com/fasmide/window-openers/manchester"
	"github.com/fasmide/window-openers/window"
	"periph.io/x/periph/host/bcm283x"
)

func main() {
	// configure pin where the 433 radio is located
	manchester.Pin = bcm283x.GPIO21

	window := window.Window{
		ID:    []byte{177, 167, 68, 228, 115},
		Radio: manchester.ShipPayload,
	}

	err := window.VerifyID()
	if err != nil {
		log.Fatalf("unusable id: %s", err)
	}

	err = window.Open()
	if err != nil {
		log.Fatalf("could not open window: %s", err)
	}

	time.Sleep(time.Second * 5)

	err = window.Close()
	if err != nil {
		log.Fatalf("could not close window: %s", err)
	}

}
