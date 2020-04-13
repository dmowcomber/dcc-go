package throttle

import (
	"fmt"
	"io"
	"log"
	"sync"
)

const maxSpeed = 127

type Throttle struct {
	serial  io.ReadWriter
	address int

	mu        sync.Mutex
	power     bool
	functions map[uint]bool
	speed     int
	direction int
}

func New(address int, serial io.ReadWriter) *Throttle {
	return &Throttle{
		address:   address,
		serial:    serial,
		functions: make(map[uint]bool),
		direction: 1, // start in forward direction
	}
}

func (t *Throttle) Reset() error {
	t.mu.Lock()
	defer t.mu.Unlock()

	for k, _ := range t.functions {
		t.functions[k] = false
	}

	t.speed = 0
	t.direction = 1
	return t.writeSpeedAndDirection()
}

func (t *Throttle) PowerToggle() error {
	t.mu.Lock()
	power := t.power
	t.mu.Unlock()

	if power {
		return t.PowerOff()
	}
	return t.PowerOn()
}

func (t *Throttle) PowerOn() error {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.power = true
	return t.writeString("<1>")
}

func (t *Throttle) PowerOff() error {
	t.mu.Lock()
	t.power = false
	t.mu.Unlock()

	t.Reset()
	return t.writeString("<0>")
}

func (t *Throttle) DirectionForward() error {
	t.mu.Lock()
	t.mu.Unlock()
	t.direction = 1
	return t.writeSpeedAndDirection()
}

func (t *Throttle) DirectionBackward() error {
	t.mu.Lock()
	t.mu.Unlock()
	t.direction = 0
	return t.writeSpeedAndDirection()
}

func (t *Throttle) ThrottleDown() error {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.speed -= 1
	if t.speed < 0 {
		t.speed = 0
	}
	return t.writeSpeedAndDirection()
}

func (t *Throttle) ThrottleUp() error {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.speed += 1
	if t.speed > maxSpeed {
		t.speed = maxSpeed
	}
	return t.writeSpeedAndDirection()
}

func (t *Throttle) Stop() error {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.speed = 0
	return t.writeSpeedAndDirection()
}

func (t *Throttle) writeSpeedAndDirection() error {
	throttlestring := fmt.Sprintf("<t 1 %d %d %d>", t.address, t.speed, t.direction)
	log.Printf("setting speed to %d and direction to %d\n", t.speed, t.direction)
	return t.writeString(throttlestring)
}

func (t *Throttle) ToggleFunction(f uint) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.functions[f] = !t.functions[f]
	var functionData uint

	// handle functions 12 and below
	if f <= 12 {
		if f <= 4 {
			functionData = 128 + t.getFunctionValue(0)*16 + t.getFunctionValue(1)*1 + t.getFunctionValue(2)*2 + t.getFunctionValue(3)*4 + t.getFunctionValue(4)*8
		} else if f <= 8 {
			functionData = 176 + t.getFunctionValue(5)*1 + t.getFunctionValue(6)*2 + t.getFunctionValue(7)*4 + t.getFunctionValue(8)*8
		} else if f <= 12 {
			functionData = 160 + t.getFunctionValue(9)*1 + t.getFunctionValue(10)*2 + t.getFunctionValue(11)*4 + t.getFunctionValue(12)*8
		}
		s := fmt.Sprintf("<f %d %d>", t.address, functionData)
		return t.writeString(s)
	}

	// handle remaining functions
	var functionPrefix uint
	if f <= 20 {
		functionPrefix = 222
		for i := uint(0); i < 8; i++ {
			functionData += t.getFunctionValue(13+i) << i
		}
	} else if f <= 28 {
		functionPrefix = 223
		for i := uint(0); i < 8; i++ {
			functionData += t.getFunctionValue(21+i) << i
		}
	} else {
		return fmt.Errorf("unknown function %d", f)
	}
	s := fmt.Sprintf("<f %d %d %d>", t.address, functionPrefix, functionData)
	return t.writeString(s)

}

func (t *Throttle) getFunctionValue(f uint) uint {
	if t.functions[f] {
		return 1
	}
	return 0
}

func (t *Throttle) writeString(s string) error {
	return t.write([]byte(s))
}

func (t *Throttle) write(data []byte) error {
	log.Printf("writing data: %s\n", data)
	_, err := t.serial.Write(data)
	if err != nil {
		return err
	}

	return nil
}
