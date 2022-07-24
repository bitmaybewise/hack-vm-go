package vm

import (
	"errors"
	"io"
	"os"
	"strings"

	"github.com/hack-vm-go/parser"
	"github.com/hack-vm-go/translator"
)

func Assemble(input *os.File) string {
	psr := parser.New(input, filenameWithoutExtension(input))
	out := new(strings.Builder)
	translator.Init(out)

	for {
		line, err := psr.ReadLine()
		if errors.Is(err, io.EOF) {
			break
		}
		if errors.Is(err, parser.IgnoredLine) {
			continue
		}

		out.WriteString("// " + line.Raw + "\n")
		out.WriteString(translator.ToAsm(line))
	}

	return out.String()
}

func filenameWithoutExtension(input *os.File) string {
	inputNameSplit := strings.Split(input.Name(), "/")
	name := inputNameSplit[len(inputNameSplit)-1]
	name = strings.Replace(name, ".vm", "", 1)
	return name
}
