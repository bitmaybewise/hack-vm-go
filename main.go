package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/hack-vm-go/vm"
)

func main() {
	var filename, dirname string
	flag.StringVar(&filename, "f", "", "the filename of the vm source file")
	flag.StringVar(&dirname, "d", "", "the directory of the vm source files")
	flag.Parse()
	if filename == "" && dirname == "" {
		panic("filename/directory is missing")
	}

	var vmFile *os.File
	defer vmFile.Close()

	if filename != "" {
		fmt.Printf("input:\t%s\n", filename)
		vmFile = openVMFile(filename)
		asm := vm.Assemble(vmFile)
		writeToAsmFile(filename, asm)
	}
	if dirname != "" {
		dirname = strings.TrimSuffix(dirname, "/")
		fmt.Printf("input:\t%s\n", dirname)
		vmFile, filename = openVMDirFiles(dirname)
		asm := vm.Assemble(vmFile)
		filename = fmt.Sprintf("%s/%s.vm", dirname, filename)
		writeToAsmFile(filename, asm)
	}
}

func openVMFile(filename string) *os.File {
	inputFile, err := os.Open(filename)
	panicsOnErrorf("error opening file\n", err)
	return inputFile
}

func openVMDirFiles(dirname string) (*os.File, string) {
	entries, err := os.ReadDir(dirname)
	panicsOnErrorf("error reading directory\n", err)

	splitDir := strings.Split(dirname, "/")
	tmpFilename := splitDir[len(splitDir)-1]
	tmpFile, err := os.CreateTemp("", tmpFilename)
	panicsOnError(err)
	defer os.Remove(tmpFile.Name())

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if strings.HasSuffix(entry.Name(), "vm") {
			vmFile := openVMFile(dirname + "/" + entry.Name())
			data, err := io.ReadAll(vmFile)
			panicsOnErrorf(fmt.Sprintf("could not read file <%s>\n", entry.Name()), err)
			_, err = tmpFile.Write(data)
			panicsOnError(err)
		}
	}

	_, err = tmpFile.Seek(0, 0)
	panicsOnError(err)
	return tmpFile, tmpFilename
}

func writeToAsmFile(filename string, content string) {
	outputFilename := strings.Replace(filename, ".vm", ".asm", 1)
	fmt.Printf("output:\t%s\n", outputFilename)

	err := os.WriteFile(outputFilename, []byte(content), 0666)
	panicsOnError(err)
}

func panicsOnError(err error) {
	if err != nil {
		panic(err)
	}
}

func panicsOnErrorf(msg string, err error) {
	if err != nil {
		panic(fmt.Sprintf("%s: <%s>", msg, err))
	}
}
