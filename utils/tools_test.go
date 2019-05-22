package utils

import (
	"bytes"
	"fmt"
	"path"
	"strings"
	"testing"
	"os/exec"
	"os"
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
		c    color
		want string
	}{
		{Blue, "\033[1;36mhello\r\n\033[0m]\r\n"},
		{Green, "\033[1;32mhello\r\n\033[0m\r\n"},
		{Yellow, "\033[1;33mhello\r\n\033[0m\r\n"},
		{Red, "\033[1;31mhello\r\n\033[0m\r\n"},
		{Title, "\033[30;42mhello\r\n\033[0m\r\n"},
		{Info, "\033[32mhello\r\n\033[0m\r\n"},
		{DEFAULT, "\033[32mhello\r\n\033[0m\r\n"},
	}
	for _, test := range tests {
		descr := fmt.Sprintf("ColorPrintln(%q, %q)", "hello", test.c)

		out = new(bytes.Buffer) // captured output
		ColorPrintln("hello", test.c)
		got := out.(*bytes.Buffer).String()
		if got != test.want {
			t.Errorf("%s = %q, want %q", descr, got, test.want)
		}
	}
}

func TestErrorPrintlnExit(t *testing.T) {

	want := "\033[1;32mhello\r\n\033[0m\r\n"

	if os.Getenv("BE_CRASHER") == "1" {
		out = new(bytes.Buffer) // captured output

		ErrorPrintln("hello", true)

		got := out.(*bytes.Buffer).String()

		if got != want {
			t.Errorf("want = %q, got %q", want, got)
		}

		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestErrorPrintlnExit")
	cmd.Env = append(os.Environ(), "BE_CRASHER=1")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && e.Exited() {
		return
	}
	t.Fatalf("process TestErrorPrintlnExit with err %v, want exit status 1", err)
}



func TestSuccessPrintln(t *testing.T) {

	want := "\033[1;32mhello\r\n\033[0m\r\n"

	out = new(bytes.Buffer) // captured output

	SuccessPrintln("hello")

	got := out.(*bytes.Buffer).String()

	if got != want {
		t.Errorf("want = %q, got %q", want, got)
	}

}

func TestCreateFolder(t *testing.T){
	defer func(){
		os.RemoveAll("test")
	}()
	err := os.MkdirAll(path.Join("test", "test1"), os.ModePerm)
	if err != nil {
		t.Errorf("create test1 folder at setup failed %v", err)
		return
	}
	tests := []string{path.Join("test", "test1"), path.Join("test", "test2"), path.Join("test","test3")}
	for _, test := range tests{
		err := CreateFolder(test)
		if err != nil {
			t.Errorf("create folder %q failed: %v", test, err)
		}
	}
}

func TestPrintTable(t *testing.T){
	out = new(bytes.Buffer) // captured output
	headers := []string{"A", "B", "C"}
	content := [][]string{{"a1", "b1", "c1"}, {"a2", "b2", "c2"}}
	want := "+----+----+----+\n| A  | B  | C  |\n+----+----+----+\n| a1 | b1 | c1 |\n| a2 | b2 | c2 |\n+----+----+----+\n"
	PrintTable(headers, content)
	got := out.(*bytes.Buffer).String()
	if got != want {
		t.Errorf("test PrintTable error got != want")
	}
}


func TestSN2PortSuccess(t *testing.T) {
	tests := []struct{sn, port string}{
		{"cas0530001", "40001"},
		{"cas0510002",  "20002"},
		{"cas0510302", "20302"},
	}
	for _, test := range tests {
		got, err := SN2Port(test.sn)
		if err != nil {
			t.Errorf("func SN2Port failed, %v", err)
			continue
		}
		if got != test.port {
			t.Errorf("got %s != want %s", got, test.port)
		}
	}

}

func TestSN2PortError(t *testing.T) {
	_, err := SN2Port("cDs0510302")
	if err == nil {
		t.Errorf("want err but err == nil")
		return
	}
	if ! strings.Contains(err.Error(), "illegal cds sn") {
		t.Errorf("get err %v", err)
	}
}