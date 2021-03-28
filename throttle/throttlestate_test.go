package throttle

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestThrottleStateDefaultValues(t *testing.T) {
	throt := &Throttle{}

	state := throt.State()
	expectedState := State{
		Functions: map[string]bool{},
	}
	assert.Equal(t, expectedState, state)

	jsonBytes, err := json.Marshal(state)
	require.NoError(t, err)
	expectedJSON := `{
		"address":0,
		"functions":{},
		"speed":0,
		"direction":0
	}`
	assert.JSONEq(t, expectedJSON, string(jsonBytes))
}

func TestThrottleState(t *testing.T) {
	throt := &Throttle{
		address: 4,
		functions: map[uint]bool{
			uint(5):  true,
			uint(11): true,
		},
		speed:     5,
		direction: 1,
	}

	state := throt.State()
	expectedState := State{
		Address: 4,
		Functions: map[string]bool{
			"5":  true,
			"11": true,
		},
		Speed:     5,
		Direction: 1,
	}
	assert.Equal(t, expectedState, state)

	jsonBytes, err := json.Marshal(state)
	require.NoError(t, err)
	expectedJSON := `{
		"address":4,
		"functions":{
			"11":true,
			"5":true
		},
		"speed":5,
		"direction":1
	}`
	assert.JSONEq(t, expectedJSON, string(jsonBytes))
}
