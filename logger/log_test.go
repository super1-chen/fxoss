package logger

import (
	"testing"
	"log"
)
func TestMylogger(t *testing.T) {
	tests := []struct{debug bool; l *log.Logger}{
		{true, debugLogger},
		{false,  traceLogger},
	}
	for _, test := range tests {
		logger := Mylogger(test.debug)
		if logger != test.l {
			t.Errorf("want %v got %v", test.l, logger)
		}
	}

}