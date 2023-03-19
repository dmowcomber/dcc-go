package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/dmowcomber/dcc-go/rail"
	"github.com/go-chi/chi"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type API struct {
	track      *rail.Track
	router     chi.Router
	httpServer *http.Server
}

func New(track *rail.Track, router chi.Router, httpServer *http.Server) *API {
	return &API{
		track:      track,
		router:     router,
		httpServer: httpServer,
	}
}

var inflightRequests *prometheus.GaugeVec

func init() {
	inflightRequests = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "inflight_http_reqeusts",
	}, []string{"endpoint"})
	prometheus.MustRegister(inflightRequests)

}

func metricsHandler(rw http.ResponseWriter, req *http.Request) {
	promhttp.Handler().ServeHTTP(rw, req)
}

func (a *API) Run() error {

	// api group
	a.router.Group(func(r chi.Router) {
		r.Use(func(next http.Handler) http.Handler {
			handler := func(w http.ResponseWriter, r *http.Request) {
				inflightRequests.WithLabelValues(r.URL.Path).Inc()
				defer inflightRequests.WithLabelValues(r.URL.Path).Dec()
				next.ServeHTTP(w, r)
			}
			return http.HandlerFunc(handler)
		})
		r.Get("/metrics", metricsHandler)
		r.Get("/{address:[0-9]+}/function", a.functionHandler)
		r.Get("/{address:[0-9]+}/speed", a.speedDirectionHandler)
		r.Get("/{address:[0-9]+}/stop", a.stopHandler)
		r.Get("/power", a.powerHandler)
		r.Get("/state", a.stateHandler)
	})

	// web ui routes
	assetsDir := http.Dir("./ui/web/assets/")
	a.router.Handle("/*", http.FileServer(assetsDir))

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
	power, err := a.track.PowerToggle()
	if err != nil {
		log.Printf("unable to toggle power: %q", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "faile to toggle power: %q", err.Error())
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
	address, err := parseAddress(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, err.Error())
		return
	}

	functionParam := r.URL.Query().Get("function")
	function, err := strconv.ParseUint(functionParam, 10, 32)
	if err != nil {
		log.Printf("function must be a number: %q", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "function must be a number: %q", err.Error())
		return
	}

	funcValue, err := a.track.GetThrottle(address).ToggleFunction(uint(function))
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
// TODO: address middleware

func parseAddress(r *http.Request) (int, error) {
	addressParam := chi.URLParam(r, "address")
	address, err := strconv.Atoi(addressParam)
	if err != nil {
		return 0, fmt.Errorf("address must be a number: %q", err.Error())
	}
	return address, nil
}

func (a *API) speedDirectionHandler(w http.ResponseWriter, r *http.Request) {
	address, err := parseAddress(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, err.Error())
		return
	}

	speedParam := r.URL.Query().Get("speed")
	speed, err := strconv.Atoi(speedParam)
	if err != nil {
		log.Printf("speed must be a number: %q", err.Error())
		// TODO: write json error
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "speed must be a number: %q", err.Error())
		return
	}

	// TODO: reuse a lot of this validation logic
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
		a.track.GetThrottle(address).DirectionForward()
	} else {
		a.track.GetThrottle(address).DirectionBackward()
	}

	a.track.GetThrottle(address).SetSpeed(speed)
}

func (a *API) stopHandler(w http.ResponseWriter, r *http.Request) {
	address, err := parseAddress(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, err.Error())
		return
	}

	err = a.track.GetThrottle(address).Stop()
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
	if a.httpServer != nil {
		return a.httpServer.Close()
	}
	return nil
}
