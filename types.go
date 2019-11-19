package vcd

import (
	"fmt"
	"strconv"
	"strings"
)

var supportedTypes = []string{"vector", "wire", "real", "string"}

var stringToType = map[string] VcdMarshall {
	"string": VcdStringType{},
	"real" : VcdRealType{},
	"vector" : VcdVectorType{bitDepth:8, maxVal:255},
}

type VcdDataType struct {
	VariableName string
	VariableType string
	BitDepth     int
	identifier   string
	marshal      VcdMarshall
}

func NewVariable(name string, variableType string, depth int) VcdDataType {
	return VcdDataType{VariableName: name, VariableType: variableType, BitDepth: depth}
}

type VcdMarshall interface {
	format(value string) (string, error)
}

type VcdRealType struct{}

func (t VcdRealType) format(value string) (string, error) {
	return fmt.Sprintf("r%s", value), nil
}

type VcdVectorType struct {
	bitDepth int
	maxVal   int
}

func (t VcdVectorType) format(value string) (string, error) {
	if value == "x" || value == "z" {
		return "b" + value, nil
	} else if num, err := strconv.Atoi(value); err == nil {
		if num > t.maxVal || num < -t.maxVal {
			return "bz", fmt.Errorf("vector is larger %d than bitdepth allows 2^%d=%d", num, t.bitDepth, t.maxVal)
		} else {
			return fmt.Sprintf("b%b", num), nil
		}
	} else {
		return "bz", fmt.Errorf("value %s is not a number, z, or x", value)
	}
}

type VcdStringType struct{}

func (t VcdStringType) format(value string) (string, error) {
	return fmt.Sprintf("s%s", strings.Replace(value, " ", "\\040", -1)), nil
}
