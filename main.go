package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/dmowcomber/dcc-go/throttle"
	"github.com/dmowcomber/dcc-go/ui/api"
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

	throt := throttle.New(address, serialWriter)

	throttleCLI := cli.New(throt)

	port := 8080
	apiServer := api.New(throt, port)

	go signalWatcher(throt, apiServer, throttleCLI)

	// run the cli
	go throttleCLI.Run()

	// run api server
	err = apiServer.Run()
	if err != nil {
		log.Fatalf("unable to start the api: %q", err.Error())
	}
}

// signalWatcher waits for a signal (control-c or kill -9).
// on SIGINT or SIGTERM it shuts everything down.
func signalWatcher(throt *throttle.Throttle, apiServer *api.API, throttleCLI *cli.CLI) {
	exitCode := 0
	defer func() {
		log.Println("powering off throttle")
		throt.PowerOff()
		os.Exit(exitCode)
	}()

	// wait for signal
	shutdownSignal := make(chan os.Signal, 1)
	signal.Notify(shutdownSignal, syscall.SIGINT, syscall.SIGTERM)
	sig := <-shutdownSignal

	log.Printf("received signal to shutdown: %s\n", sig)
	err := apiServer.Shutdown()
	if err != nil {
		log.Printf("failed to stop api: %q\n", err.Error())
		exitCode = 1
		return
	}
	log.Println("api successfully stopped")
}
