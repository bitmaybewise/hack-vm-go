package translator

import (
	"fmt"
	"strings"

	"github.com/hack-vm-go/parser"
)

func popTo(line parser.ParsedLine) translator {
	return func() string {
		out := new(strings.Builder)

		pop := func(at string) {
			pop1(out)
			out.WriteString("@" + at + "\n")
			for i := 0; i < line.Idx(); i++ {
				out.WriteString("M=M+1\n")
			}
			out.WriteString("A=M\n")
			out.WriteString("M=D\n")
			for i := 0; i < line.Idx(); i++ {
				out.WriteString("@" + at + "\n")
				out.WriteString("M=M-1\n")
			}
		}

		switch {
		case line.Segment() == "local":
			pop1(out)
			out.WriteString("@LCL\n")
			out.WriteString("A=M\n")
			out.WriteString("M=D\n")
		case line.Segment() == "argument":
			pop("ARG")
		case line.Segment() == "this":
			pop("THIS")
		case line.Segment() == "that":
			pop("THAT")
		case line.Segment() == "pointer" && line.Idx() == 0:
			pop1(out)
			out.WriteString("@THIS\n")
			out.WriteString("M=D\n")
		case line.Segment() == "pointer" && line.Idx() == 1:
			pop1(out)
			out.WriteString("@THAT\n")
			out.WriteString("M=D\n")
		case line.Segment() == "temp":
			pop1(out)
			out.WriteString(fmt.Sprintf("@R%d\n", Temp0+line.Idx()))
			out.WriteString("M=D\n")
		case line.Segment() == "static":
			pop1(out)
			out.WriteString(fmt.Sprintf("@%s.%d\n", line.Filename, line.Idx()))
			out.WriteString("M=D\n")
		}

		return out.String()
	}
}
