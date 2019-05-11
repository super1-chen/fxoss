package utils

import (
	"bytes"
	"fmt"
	"testing"
	"os/exec"
	"github.com/bouk/monkey"
)

func TestFormatItem(t *testing.T) {
	tests := []struct {
		value, maxValue int64
		want            string
	}{
		{20, 30, "20(30)"},
		{20, 50, "20(50)"},
		{20, 0, "20(-)"},
		{0, 50, "-(50)"},
		{0, 0, "-(-)"},
	}
	for _, test := range tests {
		got := FormatItem(test.value, test.maxValue)
		if got != test.want {
			t.Errorf("value %d, max %d want %s got %s", test.value, test.maxValue, test.want, got)
		}
	}
}

func TestColorPrintln(t *testing.T) {
	var tests = []struct {
		colorInfo, want string
	}{
		{Blue, "\033[1;36mhello\r\n\033[0m]\r\n"},
		{Green, "\033[1;32mhello\r\n\033[0m\r\n"},
		{Yellow, "\033[1;33mhello\r\n\033[0m\r\n"},
		{Red, "\033[1;31mhello\r\n\033[0m\r\n"},
		{Title, "\033[30;42mhello\r\n\033[0m\r\n"},
		{Info, "\033[32mhello\r\n\033[0m\r\n"},
		{"default", "\033[32mhello\r\n\033[0m\r\n"},
	}
	for _, test := range tests {
		descr := fmt.Sprintf("ColorPrintln(%q, %q)", "hello", test.colorInfo)

		out = new(bytes.Buffer) // captured output
		ColorPrintln("hello", test.colorInfo)
		got := out.(*bytes.Buffer).String()
		if got != test.want {
			t.Errorf("%s = %q, want %q", descr, got, test.want)
		}
	}
}

func TestErrorPrintlnExit1(t *testing.T) {

}
