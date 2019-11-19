package vcd

import (
	"fmt"
	"os"
	"testing"
)

var testDirectory = "vcd_test/"
var testFile = "test"

func checkFormat(formatFunc VcdMarshall, value string, expectedResponse string) error {
	resp, e := formatFunc.format(value)
	if e != nil {
		panic(e)
	}
	if resp != expectedResponse {
		return fmt.Errorf("formatting malfunction expected %s got %s", expectedResponse, resp)
	}
	return nil
}

func checkT(t *testing.T, e error) {
	if e != nil {
		t.Fatal(e)
	}
}

func TestFormatTypes(t *testing.T) {
	t.Run("Testing real formatting", func(t *testing.T) {
		checkT(t, checkFormat(stringToType["real"], "1.24", "r1.24"))
	})
	t.Run("Testing simple string formatting", func(t *testing.T) {
		checkT(t, checkFormat(stringToType["string"], "test", "stest"))
	})
	t.Run("Testing vector formatting", func(t *testing.T) {
		checkT(t, checkFormat(stringToType["vector"], "80", "b1010000"))
	})
	t.Run("Testing large string formatting", func(t *testing.T) {
		checkT(t, checkFormat(stringToType["string"], "string with space", "sstring\\040with\\040space"))
	})

}

func TestCreate(t *testing.T) {
	t.Run("Creating file", func(t *testing.T) {
		writer, e := New(testDirectory+testFile, "10")
		checkT(t, e)
		defer writer.Close()
		_, e = os.Stat(testDirectory + testFile + ".vcd")
		checkT(t, e)
		writer.SetComment("Test").SetVersion("Current")
	})
}

func TestMain(m *testing.M) {
	if _, err := os.Stat(testDirectory); os.IsNotExist(err) {
		check(os.Mkdir(testDirectory, os.ModeDir))
	}
	code := m.Run()
	if e := os.RemoveAll(testDirectory); e != nil {
		panic(e)
	}
	os.Exit(code)
}
