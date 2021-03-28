package roster

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTrackThrottles(t *testing.T) {
	track := New(&fakeReaderWriter{})
	throttle3 := track.GetThrottle(3)
	assert.NotNil(t, throttle3)
	throttle42 := track.GetThrottle(42)
	assert.NotNil(t, throttle42)

	addresses := track.GetAddresses()
	assert.Equal(t, []int{3, 42}, addresses)
}

func TestTrackPower(t *testing.T) {
	readerWriter := &fakeReaderWriter{}
	track := New(readerWriter)

	track.PowerOn()
	assert.Equal(t, "<1>", string(readerWriter.writtenBytes))

	track.PowerOff()
	assert.Equal(t, "<0>", string(readerWriter.writtenBytes))

	track.PowerToggle()
	assert.Equal(t, "<1>", string(readerWriter.writtenBytes))

	track.PowerToggle()
	assert.Equal(t, "<0>", string(readerWriter.writtenBytes))
}

type fakeReaderWriter struct {
	writeErr     error
	writtenBytes []byte
}

func (f *fakeReaderWriter) Read(b []byte) (int, error) {
	return 0, nil
}

func (f *fakeReaderWriter) Write(b []byte) (int, error) {
	f.writtenBytes = b
	return 0, f.writeErr
}
