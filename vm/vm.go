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
	psr := parser.New(input)

	out := new(strings.Builder)

	for {
		line, err := psr.ReadLine()
		if errors.Is(err, io.EOF) {
			break
		}
		if errors.Is(err, parser.IgnoredLine) {
			continue
		}

		var asm string
		if line.IsArithmetic() {
			asm = translator.Arithmetic(line.CommandType())
		}
		if line.IsPushPop() {
			inputNameSplit := strings.Split(input.Name(), "/")
			filename := inputNameSplit[len(inputNameSplit)-1]
			filename = strings.Replace(filename, ".vm", "", 1)
			asm = translator.PushPop(line.CommandType(), line.Segment(), line.Index(), filename)
		}

		out.WriteString("// " + line.Raw + "\n")
		out.WriteString(asm)
	}

	translator.EndLoop(out)
	return out.String()
}
