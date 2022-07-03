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

	var out strings.Builder

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
			asm = translator.PushPop(line.CommandType(), line.Segment(), line.Index())
		}

		out.WriteString("// " + line.Raw + "\n")
		out.WriteString(asm)
	}

	out.WriteString(translator.EndLoop())
	return out.String()
}
