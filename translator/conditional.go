package translator

import (
	"fmt"
	"strings"

	"github.com/hack-vm-go/parser"
)

func label(line parser.ParsedLine) translator {
	return func() string {
		return fmt.Sprintf("(%s)\n", line.Id())
	}
}

func goTo(line parser.ParsedLine) translator {
	return func() string {
		out := new(strings.Builder)
		out.WriteString(fmt.Sprintf("@%s\n", line.Id()))
		out.WriteString("0;JMP\n")
		return out.String()
	}
}

func ifGoTo(line parser.ParsedLine) translator {
	return func() string {
		out := new(strings.Builder)
		out.WriteString("@SP\n")
		out.WriteString("M=M-1\n")
		out.WriteString("A=M\n")
		out.WriteString("D=M\n")
		out.WriteString(fmt.Sprintf("@%s\n", line.Id()))
		out.WriteString("D;JLT\n")
		return out.String()
	}
}
