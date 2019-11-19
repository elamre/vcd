package main

import "github.com/elamre/vcd"

func main() {
	writer, e := vcd.New("example", "10ps")
	if e != nil {
		panic(e)
	}
	defer writer.Close()
	writer.SetVersion("1.0.0").SetComment("Example for VCD GO")
	_, e = writer.RegisterVariables("logic",
		vcd.NewVariable("miso", "wire", 8),
		vcd.NewVariable("mosi", "wire", 8),
		vcd.NewVariable("cs", "wire", 1),
		vcd.NewVariable("command", "string", 1))
	if e != nil {
		panic(e)
	}
	_ = writer.SetValue(0, "1", "cs")
	_ = writer.SetValue(0, "", "command")
	_ = writer.SetValue(0, "z", "mosi")
	_ = writer.SetValue(0, "z", "miso")
	_ = writer.SetValue(100, "0", "cs")
	_ = writer.SetValue(100, "String command", "command")
	_ = writer.SetValue(100, "80", "mosi")
	_ = writer.SetValue(100, "0", "miso")
	_ = writer.SetValue(200, "48", "mosi")
	_ = writer.SetValue(200, "43", "miso")
	_ = writer.SetValue(300, "90", "mosi")
	_ = writer.SetValue(300, "10", "miso")
	_ = writer.SetValue(400, "100", "mosi")
	_ = writer.SetValue(400, "90", "miso")
	_ = writer.SetValue(500, "1", "cs")
	_ = writer.SetValue(500, "z", "miso")
	_ = writer.SetValue(500, "z", "mosi")
	_ = writer.SetValue(500, "", "command")
	writer.SetTimestamp(600)
}
