package main

import (
	"log"
	"time"

	throttle "./throttle" // TODO: should be "github.com/dmowcomber/dcc-go/throttle"
	"github.com/eiannone/keyboard"
	"github.com/tarm/serial"
)

func main() {
	// TODO: allow these to be overriden as args
	address := 3
	portName := "/dev/ttyACM0"

	log.Printf("connecting to %s\n", portName)
	serialConfig := &serial.Config{Name: portName, Baud: 115200}
	port, err := serial.OpenPort(serialConfig)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("connected")

	throttle := throttle.New(address, port)
	defer func() {
		log.Println("powering off")
		time.Sleep(300 * time.Millisecond)
		throttle.PowerOff()
	}()

	userInput(throttle)
}

func userInput(t *throttle.Throttle) error {
	// TODO: print usage

	err := keyboard.Open()
	if err != nil {
		return err
	}
	defer keyboard.Close()

	for {
		char, key, err := keyboard.GetKey()
		if err != nil {
			return err
		}

		log.Printf("%q, %#v\n", char, key)

		switch char {
		case 'p':
			t.PowerToggle()
		case 'm':
			t.ToggleFunction(8)
		default:
		}

		if key == keyboard.KeyArrowUp {
			// do things
		}
		if key == keyboard.KeyCtrlC {
			return nil
		}
	}
}
