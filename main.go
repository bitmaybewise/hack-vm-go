package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/hack-vm-go/vm"
)

func main() {
	var filename string
	flag.StringVar(&filename, "f", "", "the filename of the vm source file")
	flag.Parse()
	if filename == "" {
		panic("filename is missing")
	}
	fmt.Printf("input:\t%s\n", filename)

	vmFile := openVMFile(filename)
	defer vmFile.Close()
	asm := vm.Assemble(vmFile)
	writeToAsmFile(filename, asm)
}

func openVMFile(filename string) *os.File {
	inputFile, err := os.Open(filename)
	if err != nil {
		panic(fmt.Sprintf("error opening file <%s>\n", err))
	}
	return inputFile
}

func writeToAsmFile(filename string, content string) {
	outputFilename := strings.Replace(filename, ".vm", ".asm", 1)
	fmt.Printf("output:\t%s\n", outputFilename)

	err := os.WriteFile(outputFilename, []byte(content), 0666)
	if err != nil {
		panic(err)
	}
}
