package vcd

import (
	"fmt"
	"os"
	"strings"
)

var flagNames = []string{
	"highlight", // Highlight the trace item
	"hex",       // Hexadecimal data value representation
	"dec",       // Decimal data value representation
	"bin",       // Binary data value representation
	"oct",       // Octal data value representation
	"rjustify",  // Right-justify signal name/alias
	"invert",
	"reverse",
	"exclude",
	"blank",                // Used for blank, label, and/or analog height
	"signed",               // Signed (2's compliment) data representation
	"ascii",                // ASCII character representation
	"collapsed",            // Used for closed groups
	"ftranslated",          // Trace translated with filter file
	"ptranslated",          // Trace translated with filter process
	"analog_step",          // Show trace as discrete analog steps
	"analog_interpolated",  // Show trace as analog with interpolation
	"analog_blank_stretch", // Used to extend height of analog data
	"real",                 // Read (floating point) data value representation
	"analog_fullscale",     // Analog data scaled using full simulation time
	"zerofill",
	"onefill",
	"closed",
	"grp_begin", // Begin a group of signals
	"grp_end",   // End a group of signals
	"bingray",
	"graybin",
	"real2bits",
	"ttranslated",
	"popcnt",
	"fpdecshift",
}

var flagDecoder = map[string]uint32{}

func init() {
	for i, v := range flagNames {
		flagDecoder[v] = 1 << i
	}
}

type gtkw struct {
	file *os.File
}

type gtkMarshal interface {
	toString() string
	getFlags() string
}

type gtkwTrace struct {
	gtkMarshal
	name  string
	alias string
	flags [] string
}

func getEncodedFlags(flags [] string) string {
	tempFlag := uint32(0)
	for _, flag := range flags {
		tempFlag |= flagDecoder[flag]
	}
	return fmt.Sprintf("%x", tempFlag)
}

func (trace gtkwTrace) getFlags() string {
	return fmt.Sprintf("@%s\n", getEncodedFlags(trace.flags))
}

func (trace gtkwTrace) toString() string {
	return fmt.Sprintf("+{%s} %s\n", trace.alias, trace.name)
}

func Trace(name string, alias string, flags ...string) gtkwTrace {
	return gtkwTrace{
		name:  name,
		alias: alias,
		flags: flags,
	}
}

func Gtkw(filename string) gtkw {
	if !strings.HasSuffix(filename, ".gtkw") {
		filename = filename + ".gtkw"
	}
	f, err := os.Create(filename)
	check(err)
	return gtkw{file: f}
}

func (gtkw *gtkw) SetDumpfile(dumpfile string) {
	_, _ = gtkw.file.WriteString(fmt.Sprintf("[dumpfile] \"%s\"\n", dumpfile))
}

func (gtkw *gtkw) writeFlags(flags ...string) {
	tempFlag := uint32(0)
	for _, flag := range flags {
		tempFlag |= flagDecoder[flag]
	}
	_, _ = gtkw.file.WriteString(fmt.Sprintf("@%x\n", tempFlag))
}

func (gtkw *gtkw) Group(groupName string, closed bool, traces ...gtkMarshal) {
	// Start the group
	if closed {
		gtkw.writeFlags("grp_begin", "closed", "blank")
	} else {
		gtkw.writeFlags("grp_begin", "blank")
	}
	_, _ = gtkw.file.WriteString(fmt.Sprintf("-%s\n", groupName))
	prevFlag := ""
	for _, trace := range traces {
		flag := trace.getFlags()
		if flag != prevFlag {
			_, _ = gtkw.file.WriteString(flag)
			prevFlag = flag
		}
		_, _ = gtkw.file.WriteString(trace.toString())
	}
	// End the group
	if closed {
		gtkw.writeFlags("grp_end", "closed", "blank", "collapsed")
	} else {
		gtkw.writeFlags("grp_end", "blank")
	}
	_, _ = gtkw.file.WriteString(fmt.Sprintf("-%s\n", groupName))
}

func (gtkw *gtkw) Trace(traces ...gtkMarshal) {
	prevFlag := ""
	for _, trace := range traces {
		flag := trace.getFlags()
		if flag != prevFlag {
			_, _ = gtkw.file.WriteString(flag)
			prevFlag = flag
		}
		_, _ = gtkw.file.WriteString(trace.toString())
	}
}

func (gtkw *gtkw) Close() {
	_ = gtkw.file.Close()
}
