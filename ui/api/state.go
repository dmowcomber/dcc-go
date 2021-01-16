package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (a *API) stateHandler(w http.ResponseWriter, r *http.Request) {
	resp := &stateResponse{}
	for _, address := range a.roster.GetAddresses() {

	}
	data, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error":"%s"}`, err.Error())
		return
	}
	fmt.Println(resp)
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, data)
}

type stateResponse struct {
	Throttles []throttleState `json:"throttles"`
}

type throttleState struct {
	Address   int             `json:"address"`
	Power     bool            `json:"power"`
	Functions map[string]bool `json:"functions"`
	Speed     int             `json:"speed"`
	Direction int             `json:"direction"`
}
