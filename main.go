package main

import (
	"flag"
	"log"
	"time"

	"github.com/dmowcomber/dcc-go/throttle"
	"github.com/eiannone/keyboard"
	"github.com/tarm/serial"
)

const (
	defaultAddress = 3
	defaultDevice  = "/dev/ttyACM0"
)

func main() {
	var address int
	flag.IntVar(&address, "address", defaultAddress, "The address of the locomotive")
	var device string
	flag.StringVar(&device, "device", defaultDevice, "The usb device of your DCC system")
	flag.Parse()

	log.Printf("connecting to %s\n", device)
	log.Printf("locomotive address %d\n", address)
	serialConfig := &serial.Config{Name: device, Baud: 115200}
	port, err := serial.OpenPort(serialConfig)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("connected")

	throttle := throttle.New(address, port)
	defer func() {
		log.Println("powering off")
		time.Sleep(300 * time.Millisecond)
		throttle.Stop()
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
		case 's':
			t.Stop()
		case '0':
			t.ToggleFunction(0)
		case '1':
			t.ToggleFunction(1)
		case '2', 'h':
			t.ToggleFunction(2)
		case '3':
			t.ToggleFunction(3)
		case '4':
			t.ToggleFunction(4)
		case '5':
			t.ToggleFunction(5)
		case '6':
			t.ToggleFunction(6)
		case '7':
			t.ToggleFunction(7)
		case '8', 'm':
			t.ToggleFunction(8)
		case '9':
			t.ToggleFunction(9)
		default:
			if key == keyboard.KeyArrowUp {
				t.ThrottleUp()
			} else if key == keyboard.KeyArrowDown {
				t.ThrottleDown()
			} else if key == keyboard.KeyArrowRight {
				t.DirectionForward()
			} else if key == keyboard.KeyArrowLeft {
				t.DirectionBackward()
			} else if key == keyboard.KeyCtrlC {
				return nil
			} else {
				// TODO: print usage
			}
		}
	}
}
