package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dmowcomber/dcc-go/rail"
	"github.com/dmowcomber/dcc-go/throttle"
	"github.com/dmowcomber/dcc-go/ui/api"
	"github.com/dmowcomber/dcc-go/ui/cli"
	"github.com/go-chi/chi"
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
	serialConfig := &serial.Config{
		Name:        device,
		Baud:        115200,
		ReadTimeout: 1 * time.Second,
	}
	serialWriter, err := serial.OpenPort(serialConfig)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("connected")

	port := 8080
	router := chi.NewRouter()
	addr := fmt.Sprintf(":%d", port)
	httpServer := &http.Server{
		Addr:         addr,
		Handler:      router,
		WriteTimeout: 500 * time.Millisecond,
		ReadTimeout:  1 * time.Second,
	}

	track := rail.NewTrack(serialWriter)
	apiServer := api.New(track, router, httpServer)

	throt := throttle.New(address, serialWriter)
	throttleCLI := cli.New(throt, track)
	go signalWatcher(track, apiServer, throttleCLI)

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
func signalWatcher(track *rail.Track, apiServer *api.API, throttleCLI *cli.CLI) {
	exitCode := 0
	defer func() {
		log.Println("powering off throttle")
		track.PowerOff()
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
