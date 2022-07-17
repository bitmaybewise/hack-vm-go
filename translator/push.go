package translator

import (
	"fmt"
	"strings"

	"github.com/hack-vm-go/parser"
)

func push(out *strings.Builder) {
	out.WriteString("@SP\n")
	out.WriteString("A=M\n")
	out.WriteString("M=D\n")
	out.WriteString("@SP\n")
	out.WriteString("M=M+1\n")
}

func pushTo(line parser.ParsedLine) translator {
	return func() string {
		out := new(strings.Builder)

		pop := func(at string) {
			out.WriteString("@" + at + "\n")
			for i := 0; i < line.Idx(); i++ {
				out.WriteString("M=M+1\n")
			}
			out.WriteString("A=M\n")
			out.WriteString("D=M\n")
			if line.Idx() > 0 {
				out.WriteString("@" + at + "\n")
			}
			for i := 0; i < line.Idx(); i++ {
				out.WriteString("M=M-1\n")
			}
		}

		switch {
		case line.Segment() == "constant":
			out.WriteString(fmt.Sprintf("@%d\n", line.Idx()))
			out.WriteString("D=A\n")
			push(out)
		case line.Segment() == "local":
			pop("LCL")
			push(out)
		case line.Segment() == "this":
			pop("THIS")
			push(out)
		case line.Segment() == "that":
			pop("THAT")
			push(out)
		case line.Segment() == "pointer" && line.Idx() == 0:
			out.WriteString("@THIS\n")
			out.WriteString("D=M\n")
			push(out)
		case line.Segment() == "pointer" && line.Idx() == 1:
			out.WriteString("@THAT\n")
			out.WriteString("D=M\n")
			push(out)
		case line.Segment() == "argument":
			pop("ARG")
			push(out)
		case line.Segment() == "temp":
			out.WriteString(fmt.Sprintf("@R%d\n", Temp0+line.Idx()))
			out.WriteString("D=M\n")
			push(out)
		case line.Segment() == "static":
			out.WriteString(fmt.Sprintf("@%s.%d\n", line.Filename, line.Idx()))
			out.WriteString("D=M\n")
			push(out)
		}

		return out.String()
	}
}
