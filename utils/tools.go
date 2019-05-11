package utils

import (
	"fmt"
	"io"
	"os"
	"strconv"
)

const (
	Blue   = "blue"
	Green  = "green"
	Yellow = "yellow"
	Red    = "red"
	Title  = "title"
	Info   = "Info"
)

var out io.Writer = os.Stdout // modified during testing

var Formats = map[string]string{
	Blue:      "\033[1;36m%s\r\n\033[0m]\r\n",
	Green:     "\033[1;32m%s\r\n\033[0m\r\n",
	Yellow:    "\033[1;33m%s\r\n\033[0m\r\n",
	Red:       "\033[1;31m%s\r\n\033[0m\r\n",
	Title:     "\033[30;42m%s\r\n\033[0m\r\n",
	Info:      "\033[32m%s\r\n\033[0m\r\n",
	"default": "\033[32m%s\r\n\033[0m\r\n",
}

func ColorPrintln(msg, colorInfo string) {

	var format string

	if val, ok := Formats[colorInfo]; ok {
		format = val
	} else {
		format = Formats["default"]
	}

	fmt.Fprintf(out, format, msg)
}

func ErrorPrintln(msg string, exit bool) {
	ColorPrintln(msg, Red)
	if exit {
		os.Exit(1)
	}

}

func SuccessPrintln(msg string) {
	ColorPrintln(msg, Green)
}

func FormatItem(value, maxValue int64) string {
	var valueStr string
	var maxStr string

	if value == 0 {
		valueStr = "-"
	} else {
		valueStr = strconv.Itoa(int(value))
	}

	if maxValue == 0 {
		maxStr = "-"
	} else {
		maxStr = strconv.Itoa(int(maxValue))
	}

	return fmt.Sprintf("%v(%v)", valueStr, maxStr)
}
