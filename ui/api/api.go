package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/dmowcomber/dcc-go/throttle"
)

type API struct {
	throt *throttle.Throttle

	httpServer *http.Server
}

func New(throt *throttle.Throttle, port int) *API {
	addr := fmt.Sprintf(":%d", port)
	httpServer := &http.Server{
		Addr: addr,
	}
	return &API{
		throt:      throt,
		httpServer: httpServer,
	}
}

func (a *API) Run() error {
	http.HandleFunc("/power", a.powerHandler)
	http.HandleFunc("/function", a.functionHandler)
	http.HandleFunc("/speed", a.speedDirectionHandler)
	http.HandleFunc("/stop", a.stopHandler)
	http.Handle("/", http.FileServer(http.Dir("./ui/web/assets")))
	err := a.httpServer.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("unable to run http server: %s", err.Error())
	}
	return nil
}

type powerResponse struct {
	Power bool `json:"power"`
}

type functionsResponse struct {
	Functions map[string]bool `json:"functions"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func (a *API) powerHandler(w http.ResponseWriter, r *http.Request) {
	power, err := a.throt.PowerToggle()
	if err != nil {
		log.Printf("unable to toggle power: %q", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, fmt.Sprintf("faile to toggle power: %q", err.Error()))
		return
	}

	powerResp := &powerResponse{Power: power}
	data, err := json.Marshal(powerResp)
	if err != nil {
		w.WriteHeader(http.StatusNoContent)
		log.Printf("failed to marshal response: %q", err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(data))
}

func (a *API) functionHandler(w http.ResponseWriter, r *http.Request) {
	functionParam := r.URL.Query().Get("function")
	function, err := strconv.ParseUint(functionParam, 10, 32)
	if err != nil {
		log.Printf("function must be a number: %q", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "function must be a number: %q", err.Error())
		return
	}

	funcValue, err := a.throt.ToggleFunction(uint(function))
	if err != nil {
		log.Printf("failed to toggle function: %q", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, fmt.Sprintf("failed to toggle function: %q", err.Error()))
		return
	}

	funcResp := &functionsResponse{Functions: map[string]bool{
		functionParam: funcValue,
	}}
	data, err := json.Marshal(funcResp)
	if err != nil {
		w.WriteHeader(http.StatusNoContent)
		log.Printf("failed to marshal response: %q", err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(data))
}

// TODO: use json requests data instead of parama. simpler code

func (a *API) speedDirectionHandler(w http.ResponseWriter, r *http.Request) {
	speedParam := r.URL.Query().Get("speed")
	speed, err := strconv.Atoi(speedParam)
	if err != nil {
		log.Printf("speed must be a number: %q", err.Error())
		// TODO: write json error
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "speed must be a number: %q", err.Error())
		return
	}

	forwardParam := r.URL.Query().Get("forward")
	forward, err := strconv.ParseBool(forwardParam)
	if err != nil {
		log.Printf("forward must be a boolean: %q", err.Error())
		// TODO: write json error
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "forward must be a boolean: %q", err.Error())
		return
	}

	fmt.Fprintf(w, fmt.Sprintf("setting speed to: %d, forward direction: %v", speed, forward))
	if forward {
		a.throt.DirectionForward()
	} else {
		a.throt.DirectionBackward()
	}

	a.throt.SetSpeed(speed)
}

func (a *API) stopHandler(w http.ResponseWriter, r *http.Request) {
	err := a.throt.Stop()
	if err != nil {
		errorResp := &errorResponse{
			Error: err.Error(),
		}
		data, _ := json.Marshal(errorResp)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, data)
		return
	}
}

func (a *API) Shutdown() error {
	return a.httpServer.Close()
}
