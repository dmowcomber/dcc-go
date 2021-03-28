package throttle

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestThrottleSpeed(t *testing.T) {
	address := 3
	readerWriter := &fakeReaderWriter{}
	throt := New(address, readerWriter)

	expectedSpeed := 5
	throt.SetSpeed(expectedSpeed)
	// default direction is 1 for forward
	expectedDirection := 1

	expectedData := fmt.Sprintf("<t 1 %d %d %d>", address, expectedSpeed, expectedDirection)
	assert.Equal(t, expectedData, string(readerWriter.writtenBytes))

	throt.DirectionBackward()
	expectedDirection = 0 // 0 is backward
	expectedData = fmt.Sprintf("<t 1 %d %d %d>", address, expectedSpeed, expectedDirection)
	assert.Equal(t, expectedData, string(readerWriter.writtenBytes))

	throt.Stop()
	expectedSpeed = 0 // stop should set the speed to 0
	expectedData = fmt.Sprintf("<t 1 %d %d %d>", address, expectedSpeed, expectedDirection)
	assert.Equal(t, expectedData, string(readerWriter.writtenBytes))
}

func TestThrottleIndividualFunctions(t *testing.T) {
	address := 3
	readerWriter := &fakeReaderWriter{}
	throt := New(address, readerWriter)

	testCases := []struct {
		function     uint
		expectedData string
	}{
		{function: 1, expectedData: "<f 3 129>"},
		{function: 2, expectedData: "<f 3 130>"},
		{function: 3, expectedData: "<f 3 132>"},
		{function: 4, expectedData: "<f 3 136>"},
		{function: 5, expectedData: "<f 3 177>"},
		{function: 6, expectedData: "<f 3 178>"},
		{function: 7, expectedData: "<f 3 180>"},
		{function: 8, expectedData: "<f 3 184>"},
		{function: 9, expectedData: "<f 3 161>"},
		{function: 10, expectedData: "<f 3 162>"},
		{function: 11, expectedData: "<f 3 164>"},
		{function: 12, expectedData: "<f 3 168>"},
		{function: 13, expectedData: "<f 3 222 1>"},
		{function: 14, expectedData: "<f 3 222 2>"},
		{function: 15, expectedData: "<f 3 222 4>"},
		{function: 16, expectedData: "<f 3 222 8>"},
		{function: 17, expectedData: "<f 3 222 16>"},
		{function: 18, expectedData: "<f 3 222 32>"},
		{function: 19, expectedData: "<f 3 222 64>"},
		{function: 20, expectedData: "<f 3 222 128>"},
		{function: 21, expectedData: "<f 3 223 1>"},
		{function: 22, expectedData: "<f 3 223 2>"},
		{function: 23, expectedData: "<f 3 223 4>"},
		{function: 24, expectedData: "<f 3 223 8>"},
		{function: 25, expectedData: "<f 3 223 16>"},
		{function: 26, expectedData: "<f 3 223 32>"},
		{function: 27, expectedData: "<f 3 223 64>"},
		{function: 28, expectedData: "<f 3 223 128>"},
	}
	for _, testCase := range testCases {
		enabled, err := throt.ToggleFunction(testCase.function)
		assert.NoError(t, err)
		assert.True(t, enabled)
		assert.Equal(t, testCase.expectedData, string(readerWriter.writtenBytes), fmt.Sprintf("unexpected data for function: %d", testCase.function))

		// toggle off this function for the next test
		throt.ToggleFunction(testCase.function)
	}
}

func TestThrottleAllFunctions(t *testing.T) {
	address := 3
	readerWriter := &fakeReaderWriter{}
	throt := New(address, readerWriter)

	testCases := []struct {
		functions    []uint
		expectedData string
	}{
		{functions: []uint{1, 2, 3, 4}, expectedData: "<f 3 143>"},
		{functions: []uint{5, 6, 7, 8}, expectedData: "<f 3 191>"},
		{functions: []uint{9, 10, 11, 12}, expectedData: "<f 3 175>"},
		{functions: []uint{13, 14, 15, 16}, expectedData: "<f 3 222 15>"},
		{functions: []uint{17, 18, 19, 20}, expectedData: "<f 3 222 255>"},
		{functions: []uint{21, 22, 23, 24}, expectedData: "<f 3 223 15>"},
		{functions: []uint{25, 26, 27, 28}, expectedData: "<f 3 223 255>"},
	}
	for _, testCase := range testCases {
		// enable all functions in the slice then check the last written bytes
		for _, function := range testCase.functions {
			enabled, err := throt.ToggleFunction(function)
			assert.True(t, enabled)
			assert.NoError(t, err)
		}
		assert.Equal(t, testCase.expectedData, string(readerWriter.writtenBytes), fmt.Sprintf("unexpected data for function: %#v", testCase.functions))
	}
}

func TestThrottleInvalidFunction(t *testing.T) {
	throt := New(3, &fakeReaderWriter{})
	_, err := throt.ToggleFunction(29)
	assert.Error(t, err)
}

func TestThrottleWriteError(t *testing.T) {
	readerWriter := &fakeReaderWriter{
		writeErr: errors.New("failed to write"),
	}
	throt := New(3, readerWriter)
	err := throt.SetSpeed(10)
	assert.Error(t, err)
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
