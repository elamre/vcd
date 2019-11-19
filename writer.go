package vcd

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"time"
)

type VcdError struct {
	msg string
}

func (error *VcdError) Error() string {
	return error.msg
}

// Valid timescale numbers.
var supportedTimescale = []int{1, 10, 100}

// Valid timescale units.
var supportedTimescaleUnit = []string{"s", "ms", "us", "ns", "ps", "fs"}

type VcdWriter struct {
	loadedFile          *os.File
	buffered            *bufio.Writer
	variableDefiner     int
	stringIdentifierMap map[string]VcdDataType
	previousTime        uint64
}

func New(filename string, timeScale string) (VcdWriter, error) {
	f, err := os.Create(filename)
	writer := VcdWriter{
		loadedFile:          f,
		buffered:            bufio.NewWriter(f),
		variableDefiner:     33,
		stringIdentifierMap: make(map[string]VcdDataType),
		previousTime:        0,
	}
	if err == nil {
		dat := time.Now().Format("01-02-2006 15:04:05")
		check2(writer.buffered.WriteString("$date\n\t" + dat + "\n$end\n"))
		check2(writer.buffered.WriteString("$timescale " + timeScale + " $end\n"))
		check(writer.buffered.Flush())
	}
	return writer, err
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func initVariable(variable *VcdDataType, identifier string) error {
	if !stringInSlice(variable.VariableType, supportedTypes) {
		return fmt.Errorf("unsupported data type: \"%s\" supported types: %v", variable.VariableType, supportedTypes)
	}
	variable.identifier = identifier
	switch variable.VariableType {
	case "real":
		variable.marshal = VcdRealType{}
	case "wire":
		variable.marshal = VcdVectorType{bitDepth: variable.BitDepth, maxVal: 2 << (variable.BitDepth - 1)}
	case "vector":
		variable.marshal = VcdVectorType{bitDepth: variable.BitDepth, maxVal: 2 << (variable.BitDepth - 1)}
	case "string":
		variable.marshal = VcdStringType{}
	default:
		return fmt.Errorf("not implemented datatype: \"%s\"", variable.VariableType)
	}
	return nil
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func check2(nums int, e error) {
	if e != nil {
		panic(e)
	}
}

func (vcd *VcdWriter) RegisterVariables(module string, variables ...VcdDataType) (map[string]VcdDataType, error) {
	check2(vcd.buffered.WriteString("$scope module " + module + " $end\n"))
	for _, variable := range variables {
		check(initVariable(&variable, string(vcd.variableDefiner)))

		vcd.variableDefiner = vcd.variableDefiner + 1
		response := fmt.Sprintf("%s %d %s %s", variable.VariableType, variable.BitDepth, variable.identifier, variable.VariableName)
		vcd.stringIdentifierMap[variable.VariableName] = variable
		check2(vcd.buffered.WriteString("$var " + response + " $end\n"))
	}
	check2(vcd.buffered.WriteString("$upscope $end\n"))
	check2(vcd.buffered.WriteString("$enddefinitions $end\n"))
	return vcd.stringIdentifierMap, nil
}

func (vcd *VcdWriter) SetValue(time uint64, value string, variableName string) error {
	if time < vcd.previousTime {
		return fmt.Errorf("changing value from an earlier time: %d < %d", time, vcd.previousTime)
	}
	if time != vcd.previousTime {
		_, _ = vcd.buffered.WriteString("#" + strconv.FormatUint(time, 10) + "\n")
		vcd.previousTime = time
	}
	format, e := format(value)
	if e != nil {
		panic(e)
	}
	check2(vcd.buffered.WriteString(format + " " + vcd.stringIdentifierMap[variableName].identifier + "\n"))
	return e
}

func (vcd VcdWriter) SetComment(comment string) VcdWriter {
	check2(vcd.buffered.WriteString("$comment\n\t" + comment + "\n$end\n"))
	return vcd
}

func (vcd VcdWriter) SetVersion(version string) VcdWriter {
	check2(vcd.buffered.WriteString("$version\n\t" + version + "\n$end\n"))
	return vcd
}

func (vcd VcdWriter) Close() {
	check(vcd.buffered.Flush())
	check(vcd.loadedFile.Close())
}
