package main

import (
	"flag"
	"log"
	"time"

	"github.com/dmowcomber/dcc-go/throttle"
	"github.com/dmowcomber/dcc-go/ui/cli"
	"github.com/tarm/serial"
)

const (
	defaultAddress = 3
	defaultDevice  = "/dev/ttyACM0"
)

func main() {
	var address int
	var device string
	flag.IntVar(&address, "address", defaultAddress, "The address of the locomotive")
	flag.StringVar(&device, "device", defaultDevice, "The usb device of your DCC system")
	flag.Parse()

	log.Printf("connecting to %s\n", device)
	log.Printf("locomotive address %d\n", address)
	serialConfig := &serial.Config{Name: device, Baud: 115200}
	serialWriter, err := serial.OpenPort(serialConfig)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("connected")

	throttle := throttle.New(address, serialWriter)
	defer func() {
		log.Println("powering off")
		time.Sleep(300 * time.Millisecond)
		throttle.Stop()
		throttle.PowerOff()
	}()

	// run the cli
	throttleCLI := cli.New(throttle)
	throttleCLI.Run()
}
