package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/hack-vm-go/translator"
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

	out := new(strings.Builder)
	translator.Init(out)

	if filename != "" {
		fmt.Printf("input:\t%s\n", filename)
		vmFile := openVMFile(filename)
		defer vmFile.Close()
		vm.Assemble(vmFile, out)
	}
	if dirname != "" {
		dirname = strings.TrimSuffix(dirname, "/")
		fmt.Printf("input:\t%s\n", dirname)

		for _, vmFile := range dirFiles(dirname) {
			defer vmFile.Close()
			vm.Assemble(vmFile, out)
		}

		splitDir := strings.Split(dirname, "/")
		filename = splitDir[len(splitDir)-1]
		filename = fmt.Sprintf("%s/%s.vm", dirname, filename)
	}

	writeToAsmFile(filename, out.String())
}

func openVMFile(filename string) *os.File {
	inputFile, err := os.Open(filename)
	panicsOnErrorf("error opening file\n", err)
	return inputFile
}

func dirFiles(dirname string) []*os.File {
	entries, err := os.ReadDir(dirname)
	panicsOnErrorf("error reading directory\n", err)

	splitDir := strings.Split(dirname, "/")
	tmpFilename := splitDir[len(splitDir)-1]
	tmpFile, err := os.CreateTemp("", tmpFilename)
	panicsOnError(err)
	defer os.Remove(tmpFile.Name())

	files := make([]*os.File, 0)
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if strings.HasSuffix(entry.Name(), "vm") {
			file := openVMFile(dirname + "/" + entry.Name())
			files = append(files, file)
		}
	}

	return files
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
