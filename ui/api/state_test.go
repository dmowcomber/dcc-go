package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStateHandler(t *testing.T) {
	req := &http.Request{}
	rw := &fakeResponseWriter{}

	api := &API{}
	api.stateHandler(rw, req)

	assert.Equal(t, http.StatusOK, rw.statusCodeWritten)

	// expectedData := ``
	// data := bytes.NewBuffer(rw.bytesWritten).String()
	m := make(map[string]interface{})
	err := json.Unmarshal(rw.bytesWritten, m)
	assert.NoError(t, err)
	fmt.Println(string(rw.bytesWritten))
	assert.Equal(t, map[string]interface{}{}, m)
	// assert.Equal(t, expectedData, string(rw.bytesWritten))
	// assert.Equal(t, expectedData, data)
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
