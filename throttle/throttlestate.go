package throttle

import "fmt"

type State struct {
	Address   int             `json:"address"`
	Functions map[string]bool `json:"functions"`
	Speed     int             `json:"speed"`
	Direction int             `json:"direction"`
}

func (t *Throttle) State() State {
	t.mu.Lock()
	defer t.mu.Unlock()

	state := State{
		Address:   t.address,
		Speed:     t.speed,
		Direction: t.direction,
	}

	state.Functions = make(map[string]bool, len(t.functions))
	for function, enabled := range t.functions {
		functionName := fmt.Sprintf("%v", function)
		state.Functions[functionName] = enabled
	}

	return state
}
