package vcd

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Both vector and wire result to vectortype
// TODO implement missing types
var supportedTypes = []string{"vector", "wire", "real", "string"}

// Not really an error, but prevents writing of empty strings when there was no change. This causes glitches
// TODO maybe eventually add a boolean return for every marshal instead of throwing and comparing errors
var duplicateErr = errors.New("duplicate")

var stringToType = map[string]vcdMarshall{
	"string": &VcdStringType{},
	"real":   VcdRealType{},
	"vector": VcdVectorType{bitDepth: 8, maxVal: 255},
}

type VcdDataType struct {
	VariableName string
	VariableType string
	BitDepth     int
	identifier   string
	marshal      vcdMarshall
}

// Creates a new variable which can be registered
// Panics when trying to create an unkown type
// Depth argument is used for vector types
// See supportedTypes for the supported types
func NewVariable(name string, variableType string, depth int) VcdDataType {
	if !stringInSlice(variableType, supportedTypes) {
		errorStr := fmt.Sprintf("unsupported type: %s\nUse one of the following: %v", variableType, supportedTypes)
		panic(errorStr)
	}
	return VcdDataType{VariableName: name, VariableType: variableType, BitDepth: depth}
}

// Marshaller that should be implemented by every type
type vcdMarshall interface {
	format(value string) (string, error)
}

// Defines real types such as 1.4 -3.4
type VcdRealType struct{}

func (t VcdRealType) format(value string) (string, error) {
	return fmt.Sprintf("r%s", value), nil
}

// Defines vector types such as 01010101
type VcdVectorType struct {
	bitDepth int
	maxVal   uint64
}

func (t VcdVectorType) format(value string) (string, error) {
	if value == "x" || value == "z" {
		return "b" + value, nil
	} else if num, err := strconv.ParseInt(value, 10, 64); err == nil {
		if uint64(num) > t.maxVal {
			return "bz", fmt.Errorf("vector is larger %d than bitdepth allows 2^%d=%d", num, t.bitDepth, t.maxVal)
		} else {
			return fmt.Sprintf("b%b", num), nil
		}
	} else {
		return "bz", fmt.Errorf("value %s is not a number, z, or x", value)
	}
}

// Defines string types
type VcdStringType struct {
	empty bool
}

func (t * VcdStringType) format(value string) (string, error) {
	if value == "" && t.empty{
		return "", duplicateErr
	}else if value == "" {
		t.empty = true
	}else{
		t.empty = false
	}
	return fmt.Sprintf("s%s", strings.Replace(value, " ", "\\040", -1)), nil
}
