package rail

import (
	"errors"
	"io"
	"sort"
	"sync"
	"time"

	"github.com/dmowcomber/dcc-go/throttle"
	"github.com/prometheus/client_golang/prometheus"
)

var serialWriteLatency prometheus.Histogram

func init() {
	serialWriteLatency = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "serial_write_latency",
		Buckets: []float64{.005, .05, .1, .5, 1, 2, 3, 4, 5, 10, 30, 60},
	})
	prometheus.MustRegister(serialWriteLatency)
}

// Track controls power to the track and holds each of the throttles
type Track struct {
	serial io.ReadWriter
	mu     *sync.Mutex
	power  bool
	roster map[int]*throttle.Throttle
}

// New returns a new Track
func New() *Track {
	return &Track{
		mu:     &sync.Mutex{},
		roster: make(map[int]*throttle.Throttle),
	}
}

func (t *Track) SetWriter(serial io.ReadWriter) {
	t.serial = serial
}

// GetThrottle returns a throttle for a given address.
// It creates one if needed.
func (t *Track) GetThrottle(address int) *throttle.Throttle {
	t.mu.Lock()
	defer t.mu.Unlock()

	throt, ok := t.roster[address]
	if !ok {
		throt = throttle.New(address)
		throt.SetWriter(t.serial)
		t.roster[address] = throt
	}
	return throt
}

// GetAddresses returns a list of address that have a throttle
func (t *Track) GetAddresses() []int {
	if t == nil || t.roster == nil {
		return nil
	}
	addresses := make([]int, 0, len(t.roster))
	for address := range t.roster {
		addresses = append(addresses, address)
	}

	sort.Ints(addresses)
	return addresses
}

func (t *Track) IsPowerOn() bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.power
}

func (t *Track) PowerOn() error {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.power = true
	return t.writeString("<1>")
}

func (t *Track) PowerOff() error {
	t.mu.Lock()
	t.power = false
	t.mu.Unlock()

	// tell each throttle to turn off functions
	// so that when we turn things back on we don't have a horn blaring
	for _, throt := range t.roster {
		throt.Reset()
	}
	return t.writeString("<0>")
}

func (t *Track) PowerToggle() (bool, error) {
	t.mu.Lock()
	power := t.power
	t.mu.Unlock()

	if power {
		return false, t.PowerOff()
	}
	return true, t.PowerOn()
}

func (t *Track) writeString(s string) error {
	return t.write([]byte(s))
}

func (t *Track) write(data []byte) error {
	if t.serial == nil {
		return errors.New("serial writer has not been initialized")
	}
	start := time.Now()
	defer func() {
		serialWriteLatency.Observe(time.Since(start).Seconds())
	}()
	// log.Printf("writing data: %s\n", data)
	_, err := t.serial.Write(data)
	if err != nil {
		return err
	}

	return nil
}
