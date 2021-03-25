package roster

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoster(t *testing.T) {
	r := New(&fakeReadWriter{})
	throttle3 := r.GetThrottle(3)
	assert.NotNil(t, throttle3)
	throttle42 := r.GetThrottle(42)
	assert.NotNil(t, throttle42)

	addresses := r.GetAddresses()
	assert.Equal(t, []int{3, 42}, addresses)
}

type fakeReadWriter struct{}

func (f *fakeReadWriter) Read(p []byte) (n int, err error) {
	return 0, nil
}

func (f *fakeReadWriter) Write(p []byte) (n int, err error) {
	return 0, nil
}
