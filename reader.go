package vcd

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

type ReadValue struct {
	Time  int64
	value interface{}
}

type VcdReader struct {
	loadedFile *os.File
	buffered   *bufio.Reader
	time       int64

	date              string
	timescale         string
	version           string
	comment           string
	identifierNameMap map[string]VcdDataType
}

func NewReader(filename string) (VcdReader, error) {
	reader := VcdReader{}
	f, err := os.Open(filename)
	reader.loadedFile = f
	reader.buffered = bufio.NewReader(reader.loadedFile)
	reader.identifierNameMap = nil
	return reader, err
}

func (reader VcdReader) GetIdentifiers() map[string]VcdDataType {
	return reader.identifierNameMap
}

func (reader *VcdReader) ParseHeader() {
	reader.identifierNameMap = make(map[string]VcdDataType)
	module := ""
	for ; ; {
		readStream := ""
		for ; !strings.HasSuffix(readStream, "$end\n"); {
			b, err := reader.buffered.ReadByte()
			check(err)
			readStream += string(b)
			if !strings.HasPrefix(readStream, "$") {
				check(reader.buffered.UnreadByte())
				// Todo log Warning
				return
			}
		}
		readStream = strings.ReplaceAll(readStream, "\n\r", " ")
		readStream = strings.ReplaceAll(readStream, "\r\n", " ")
		readStream = strings.ReplaceAll(readStream, "\n", " ")
		readStream = strings.ReplaceAll(readStream, "   ", " ")
		readStream = strings.ReplaceAll(readStream, "  ", " ")
		delim := strings.Split(readStream, " ")
		l := len(delim)
		switch delim[0] {
		case "$scope":
			module += delim[2] + "."
		case "$upscope":
			ss := strings.Split(module, ".")
			module = ""
			for _, up := range ss[:len(ss)-1] {
				module += up + "."
			}
		case "$comment":
			for _, ss := range delim[1 : l-2] {
				reader.comment += ss + " "
			}
		case "$var":
			var datType VcdDataType
			datType.VariableType = delim[1]
			datType.BitDepth = int(check2(strconv.ParseInt(delim[2], 10, 32)).(int64))
			datType.identifier = delim[3]
			datType.VariableName = module + delim[4]
			switch delim[1] {
			case "wire":
				datType.marshal = &VcdVectorType{}
			case "string":
				datType.marshal = &VcdStringType{}
			case "real":
				datType.marshal = VcdRealType{}
			default:
				panic("Unknown variable")
			}
			reader.identifierNameMap[datType.identifier] = datType
		case "$date":
			reader.date = delim[1] + " " + delim[2]
		case "$version":
			for _, ss := range delim[1 : l-2] {
				reader.version += ss
			}
		case "$timescale":
			reader.timescale = delim[1]
		case "$enddefinitions":
			break
		default:
			log.Printf("Unknown header: %s", delim[0])
		}
	}
}

func (reader *VcdReader) ReadAll() map[string][]ReadValue {
	if reader.identifierNameMap == nil {
		reader.ParseHeader()
	}
	retVal := make(map[string][]ReadValue)
	for k, _ := range reader.identifierNameMap {
		retVal[k] = make([]ReadValue, 0)
	}
	for ; ; {
		valid, time, identifier, value := reader.Next()
		if !valid {
			break
		}
		// TODO use some linkedlist
		retVal[identifier] = append(retVal[identifier], ReadValue{time, value})
	}
	for k, v := range reader.identifierNameMap {
		retVal[v.VariableName] = retVal[k]
		delete(retVal, k)
	}
	return retVal
}

// return: Completed, Time, Identifier, Value
func (reader *VcdReader) Next() (bool, int64, string, interface{}) {
	var line []byte
	var err error

	for ; ; {
		line, _, err = reader.buffered.ReadLine()
		if err != nil {
			return false, 0, "", ""
		}
		dec := strings.Split(string(line), " ")
		if len(dec) == 1 {
			time, err := strconv.ParseInt(string(line[1:]), 10, 64)
			if err != nil {
				panic(err)
			}
			reader.time = time
		} else {
			val, err := reader.identifierNameMap[dec[1]].marshal.parse(dec[0])
			if err != nil {
				// More info
				log.Printf("Error unmarshalling")
			} else {
				return true, reader.time, dec[1], val
			}
		}

	}

	return false, 0, "", ""
}
