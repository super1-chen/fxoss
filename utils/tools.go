package utils

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"
)

type color int

const (
	Blue    = color(iota)
	Green
	Yellow
	Red
	Title
	Info
	DEFAULT
)

var (
	out     io.Writer = os.Stdout // modified during testing
	formats           = map[color]string{
		Blue:    "\033[1;36m%s\r\n\033[0m]\r\n",
		Green:   "\033[1;32m%s\r\n\033[0m\r\n",
		Yellow:  "\033[1;33m%s\r\n\033[0m\r\n",
		Red:     "\033[1;31m%s\r\n\033[0m\r\n",
		Title:   "\033[30;42m%s\r\n\033[0m\r\n",
		Info:    "\033[32m%s\r\n\033[0m\r\n",
		DEFAULT: "\033[32m%s\r\n\033[0m\r\n",
	}
)

// ColorPrintln println message in different color
func ColorPrintln(msg string, c color) {

	var format string

	if val, ok := formats[c]; ok {
		format = val
	} else {
		format = formats[c]
	}

	fmt.Fprintf(out, format, msg)
}

// ErrorPrintln print message in color read
func ErrorPrintln(msg string, exit bool) {
	ColorPrintln(msg, Red)
	if exit {
		os.Exit(1)
	}

}

// SuccessPrintln print message in color green
func SuccessPrintln(msg string) {
	ColorPrintln(msg, Green)
}

// FormatItem format value max to string likes value(max)
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

// CreateFolder create folder by give name if the given folder is not exites
func CreateFolder(folder string) error {
	_, err := os.Stat(folder)
	if os.IsNotExist(err) {
		// todo need add a logger
		// log.info("create new log)
		if err := os.MkdirAll(folder, os.ModePerm); err != nil {
			return fmt.Errorf("create fold %s failed %v", folder, err)
		}
	}
	return nil
}

// PrintTable print ascii table
func PrintTable(headers []string, content [][]string) {
	table := tablewriter.NewWriter(out)
	table.SetHeader(headers)
	table.AppendBulk(content)
	table.Render()
}

// SN2Port converts sn 2 frpc port
func SN2Port(sn string) (port string, err error) {
	if strings.HasPrefix(strings.ToUpper(sn), "CAS053") {
		port = "40" + sn[len(sn)-3:]
	} else if strings.HasPrefix(strings.ToUpper(sn), "CAS051") {
		port = "20" + sn[len(sn)-3:]
	} else {
		return port, fmt.Errorf("illegal cds sn %s", sn)
	}
	return
}
