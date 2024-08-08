package api

import (
	"net/http"
	"testing"

	"github.com/dmowcomber/dcc-go/rail"
	"github.com/stretchr/testify/assert"
)

func TestStateHandlerDefaultState(t *testing.T) {
	req := &http.Request{}
	rw := &fakeResponseWriter{}
	track := rail.New()
	track.SetWriter(&noopReaderWriter{})
	// create new addresses
	track.GetThrottle(3)
	track.GetThrottle(42)

	api := &API{
		track: track,
	}
	api.stateHandler(rw, req)

	assert.Equal(t, http.StatusOK, rw.statusCodeWritten)
	expectedJSON := `{
		"power": false,
		"throttles":
		{
			"3": {
				"address":3,
				"functions":{},
				"speed":0,
				"direction":1
			},
			"42": {
				"address":42,
				"functions":{},
				"speed":0,
				"direction":1
			}
		}
	}`
	assert.JSONEq(t, expectedJSON, string(rw.bytesWritten))
}

func TestStateHandler(t *testing.T) {
	req := &http.Request{}
	rw := &fakeResponseWriter{}
	track := rail.New()
	track.SetWriter(&noopReaderWriter{})
	api := &API{
		track: track,
	}
	track.PowerOn()

	throt3 := track.GetThrottle(3)
	throt3.ToggleFunction(4)
	throt3.ToggleFunction(28)
	throt3.SetSpeed(6)
	throt3.DirectionForward()

	throt42 := track.GetThrottle(42)
	throt42.DirectionBackward()

	api.stateHandler(rw, req)

	assert.Equal(t, http.StatusOK, rw.statusCodeWritten)
	expectedJSON := `{
		"power": true,
		"throttles":
		{
			"3": {
				"address":3,
				"functions": {
					"4": true,
					"28": true
				},
				"speed":6,
				"direction":1
			},
			"42": {
				"address":42,
				"functions":{},
				"speed":0,
				"direction":0
			}
		}
	}`
	assert.JSONEq(t, expectedJSON, string(rw.bytesWritten))
}

func TestStateHandlerNoThrottles(t *testing.T) {
	req := &http.Request{}
	rw := &fakeResponseWriter{}
	track := rail.New()
	track.SetWriter(&noopReaderWriter{})
	api := &API{
		track: track,
	}
	api.stateHandler(rw, req)

	assert.Equal(t, http.StatusOK, rw.statusCodeWritten)
	expectedJSON := `{
		"power": false,
		"throttles":{}
	}`
	assert.JSONEq(t, expectedJSON, string(rw.bytesWritten))
}

type fakeResponseWriter struct {
	respHeader     http.Header
	respWriteCount int
	respWriteErr   error

	bytesWritten      []byte
	statusCodeWritten int
}

func (f *fakeResponseWriter) Header() http.Header {
	return f.respHeader
}

func (f *fakeResponseWriter) Write(b []byte) (int, error) {
	f.bytesWritten = b
	return f.respWriteCount, f.respWriteErr
}

func (f *fakeResponseWriter) WriteHeader(statusCode int) {
	f.statusCodeWritten = statusCode
}

type noopReaderWriter struct{}

func (n *noopReaderWriter) Read(p []byte) (int, error) {
	return 0, nil
}

func (n *noopReaderWriter) Write(p []byte) (int, error) {
	return 0, nil
}
