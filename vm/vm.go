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
		if line.IsArithmeticOrLogic() {
			asm = translator.ArithmeticLogic(line.CommandType())
		}
		if line.IsPushPop() {
			asm = translator.PushPop(line.CommandType(), line.Segment(), line.Index(), filename(input))
		}

		out.WriteString("// " + line.Raw + "\n")
		out.WriteString(asm)
	}

	translator.EndLoop(out)
	return out.String()
}

func filename(input *os.File) string {
	inputNameSplit := strings.Split(input.Name(), "/")
	name := inputNameSplit[len(inputNameSplit)-1]
	name = strings.Replace(name, ".vm", "", 1)
	return name
}
