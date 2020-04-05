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
	functions map[int]bool
	speed     int
	direction int
}

func New(address int, serial io.ReadWriter) *Throttle {
	return &Throttle{
		address:   address,
		serial:    serial,
		functions: make(map[int]bool),
	}
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
	defer t.mu.Unlock()

	t.power = false
	return t.writeString("<0>")
}

func (t *Throttle) DirectionPositive() error {
	t.direction = 1
	throttlestring := fmt.Sprintf("<t 1 %d %d %d>", t.address, t.speed, t.direction)
	log.Printf("speed set to %d", t.speed)
	return t.writeString(throttlestring)
}

func (t *Throttle) DirectionZero() error {
	t.direction = 0
	throttlestring := fmt.Sprintf("<t 1 %d %d %d>", t.address, t.speed, t.direction)
	log.Printf("speed set to %d", t.speed)
	return t.writeString(throttlestring)
}

func (t *Throttle) ThrottleDown() error {
	t.speed -= 1
	if t.speed < 0 {
		t.speed = 0
	}
	throttlestring := fmt.Sprintf("<t 1 %d %d %d>", t.address, t.speed, t.direction)
	log.Printf("speed set to %d", t.speed)
	return t.writeString(throttlestring)
}

func (t *Throttle) ThrottleUp() error {
	t.speed += 1
	if t.speed > maxSpeed {
		t.speed = maxSpeed
	}
	throttlestring := fmt.Sprintf("<t 1 %d %d %d>", t.address, t.speed, t.direction)
	log.Printf("speed set to %d", t.speed)
	return t.writeString(throttlestring)
}

func (t *Throttle) ToggleFunction(f int) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.toggleFunctionValue(f)

	var functionData int
	if f <= 4 {
		functionData = 128 + t.getFunctionValue(0)*16 + t.getFunctionValue(1)*1 + t.getFunctionValue(2)*2 + t.getFunctionValue(3)*4 + t.getFunctionValue(4)*8
	} else if f <= 8 {
		functionData = 176 + t.getFunctionValue(5)*1 + t.getFunctionValue(6)*2 + t.getFunctionValue(7)*4 + t.getFunctionValue(8)*8
	} else {
		return fmt.Errorf("function %d not implemented yet", f)
	}

	s := fmt.Sprintf("<f %d %d>", t.address, functionData)
	return t.writeString(s)
}

func (t *Throttle) getFunctionValue(f int) int {
	if t.functions[f] {
		return 1
	}
	return 0
}

func (t *Throttle) toggleFunctionValue(f int) int {
	t.functions[f] = !t.functions[f]
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
