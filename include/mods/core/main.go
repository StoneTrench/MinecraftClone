package main

import "fmt"

func init() {
	log_info("Logging from mod!")
	log_info(fmt.Sprintf("Received api version major: %d", api_version()))
}

//go:wasmimport env log
func log_info(msg string)

//go:wasmimport env api_version
func api_version() uint64

func main() {}
