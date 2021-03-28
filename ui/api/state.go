package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dmowcomber/dcc-go/throttle"
)

func (a *API) stateHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: add optional address param to reduce the payload size? or add a new path with address in it
	addrs := a.track.GetAddresses()
	resp := &stateResponse{
		Power:     a.track.IsPowerOn(),
		Throttles: make(map[int]throttle.State, len(addrs)),
	}

	for _, address := range addrs {
		throt := a.track.GetThrottle(address)
		resp.Throttles[address] = throt.State()
	}
	data, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error":"%s"}`, err.Error())
		return
	}
	fmt.Println(resp)
	w.WriteHeader(http.StatusOK)
	fmt.Printf("writing data: %s\n", string(data))
	w.Write(data)
}

type stateResponse struct {
	// Throttles map of Address to Throttle State
	// Throttles map[int]throttleState  `json:"throttles"`
	Throttles map[int]throttle.State `json:"throttles"`
	Power     bool                   `json:"power"`
}
