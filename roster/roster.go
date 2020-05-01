package roster

import (
	"io"
	"sync"

	"github.com/dmowcomber/dcc-go/throttle"
)

// Roster generates Throttles by address
type Roster struct {
	serial    io.ReadWriter
	mu        *sync.Mutex
	throttles map[int]*throttle.Throttle
}

// New returns a new Roster
func New(serial io.ReadWriter) *Roster {
	return &Roster{
		serial:    serial,
		mu:        &sync.Mutex{},
		throttles: make(map[int]*throttle.Throttle),
	}
}

// GetThrottle returns a throttle for a given address.
// It creates one if needed.
func (r *Roster) GetThrottle(address int) *throttle.Throttle {
	r.mu.Lock()
	defer r.mu.Unlock()

	throt, ok := r.throttles[address]
	if !ok {
		throt = throttle.New(address, r.serial)
		r.throttles[address] = throt
	}
	return throt
}

// GetAddresses returns a list of address that have throttles
func (r *Roster) GetAddresses() []int {
	addresses := make([]int, 0, len(r.throttles))
	for address := range r.throttles {
		addresses = append(addresses, address)
	}
	return addresses
}
