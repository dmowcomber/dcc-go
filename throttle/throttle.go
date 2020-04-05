package throttle

import (
	"errors"
	"fmt"
	"io"
	"log"
)

type Throttle struct {
	serial  io.ReadWriter
	address int

	power     bool
	functions map[int]bool
}

func New(address int, serial io.ReadWriter) *Throttle {
	return &Throttle{
		address:   address,
		serial:    serial,
		functions: make(map[int]bool),
	}
}
func (t *Throttle) PowerToggle() error {
	if t.power {
		return t.PowerOff()
	}
	return t.PowerOn()
}

func (t *Throttle) PowerOn() error {
	t.power = true
	return t.write([]byte("<1>"))
}

func (t *Throttle) PowerOff() error {
	t.power = false
	return t.write([]byte("<0>"))
}

func (t *Throttle) ToggleFunction(f int) error {
	if f != 8 {
		return errors.New("not implemented yet")
	}

	t.functions[f] = !t.functions[f]
	functionIsOn := boolToInt(t.functions[f])

	functionValue := 176 + (functionIsOn * 8)
	s := fmt.Sprintf("<f %d %d>", t.address, functionValue)
	return t.writeString(s)
}

func boolToInt(b bool) int {
	if b {
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
