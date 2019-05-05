package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"periph.io/x/periph/conn/gpio/gpioreg"
	"periph.io/x/periph/host/bcm283x"

	"github.com/fasmide/window-openers/manchester"
	"github.com/fasmide/window-openers/window"
	"github.com/spf13/viper"
)

func init() {
	// Defaults
	viper.SetDefault("pin", "GPIO21")

	// look for configuration at
	// /etc/window-openers.yml
	// $HOME/window-openers.yml
	// current work directory/window-openers.yml
	viper.SetConfigName("window-openers")
	viper.AddConfigPath("/etc/")
	viper.AddConfigPath("$HOME/")
	viper.AddConfigPath(".")

	// use YML format
	viper.SetConfigType("yml")

	// find and read configuration
	// hint if no configuration was found - the user could simply
	// be using environment variables for everything
	err := viper.ReadInConfig()
	if err != nil {
		log.Printf("no configuration could be read: %s", err)
	}

	viper.AutomaticEnv()
}

func main() {
	if len(os.Args) != 3 {
		log.Fatalf("please provide window-name and action\nE.g. window-openers Bedroom open")
	}

	// configure pin where the 433 radio is located
	pin, ok := gpioreg.ByName(viper.GetString("pin")).(*bcm283x.Pin)
	if !ok {
		log.Fatalf("unable to cast %s pin to bcm283x type - not using a raspberry pi?", viper.GetString("pin"))
	}
	manchester.Pin = pin

	// find window in question from configuration
	// and unmarshal into window
	windowName := os.Args[1]
	windowID := viper.GetString(
		fmt.Sprintf("windows.%s.ID", windowName),
	)
	if windowID == "" {
		log.Fatalf("no configuration found key windows.%s.ID", windowName)
	}

	// the window ID is now presented as decimal bytes space seperated
	// e.g. 123 163 121 ...
	// we must turn this into a byte slice
	fields := strings.Fields(windowID)
	if len(fields) != 5 {
		log.Fatalf("%s id should be excatly 5 bytes", windowName)
	}

	id := make([]byte, 5)
	for i, decimal := range fields {
		b, err := strconv.Atoi(decimal)
		if err != nil {
			log.Fatalf("could not convert %s into byte slice", windowID)
		}
		id[i] = byte(b)
	}

	// initialize window
	window := window.Window{
		ID:    id,
		Radio: manchester.ShipPayload,
	}

	err := window.VerifyID()
	if err != nil {
		log.Fatalf("unusable id: %s", err)
	}

	action := os.Args[2]

	if action == "open" {
		err = window.Open()
		if err != nil {
			log.Fatalf("could not open window: %s", err)
		}
		return
	}

	if action == "close" {
		err = window.Close()
		if err != nil {
			log.Fatalf("could not close window: %s", err)
		}
		return
	}

	if action == "stop" {
		err = window.Stop()
		if err != nil {
			log.Fatalf("could not close window: %s", err)
		}
		return
	}

	if action == "pair" {
		err = window.Pair()
		if err != nil {
			log.Fatalf("could not Pair with window: %s", err)
		}
		return
	}

	// if we reached this far - the user did not specify a valid action
	log.Fatalf("%s is not a valid action, actions are open, close, stop and pair", action)

}
