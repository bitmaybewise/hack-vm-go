package translator

import (
	"fmt"
	"strings"

	"github.com/hack-vm-go/parser"
)

func function(line parser.ParsedLine) translator {
	dZero := func(out *strings.Builder) {
		out.WriteString("@0\n")
		out.WriteString("D=A\n")
	}

	return func() string {
		out := new(strings.Builder)
		out.WriteString(fmt.Sprintf("(%s)\n", line.Segment()))
		for i := 0; i < line.Idx(); i++ {
			dZero(out)
			push(out)
		}
		return out.String()
	}
}

func functionReturn() translator {
	return func() string {
		out := new(strings.Builder)
		frameRef := func(n int) string {
			return fmt.Sprintf("@R%d\n", Temp0+n)
		}
		frameMinus := func(frame, n int) {
			out.WriteString(frameRef(frame))
			out.WriteString("D=M\n")
			for i := 0; i < n; i++ {
				out.WriteString("D=D-1\n")
			}
			out.WriteString("A=D\n")
			out.WriteString("D=M\n")
		}
		// frame = LCL
		out.WriteString("@LCL\n")
		out.WriteString("D=M\n")
		out.WriteString(frameRef(0))
		out.WriteString("M=D\n")
		// retAddr = *(frame-5)
		frameMinus(0, 5)
		out.WriteString(frameRef(1))
		out.WriteString("M=D\n")
		// *ARG = pop()
		pop1(out)
		out.WriteString("@ARG\n")
		out.WriteString("A=M\n")
		out.WriteString("M=D\n")
		// SP = ARG+1
		out.WriteString("@ARG\n")
		out.WriteString("D=M\n")
		out.WriteString("@SP\n")
		out.WriteString("M=D+1\n")
		// THAT = *(frame-1)
		frameMinus(0, 1)
		out.WriteString("@THAT\n")
		out.WriteString("M=D\n")
		// THIS = *(frame-2)
		frameMinus(0, 2)
		out.WriteString("@THIS\n")
		out.WriteString("M=D\n")
		// ARG = *(frame-3)
		frameMinus(0, 3)
		out.WriteString("@ARG\n")
		out.WriteString("M=D\n")
		// LCL = *(frame-4)
		frameMinus(0, 4)
		out.WriteString("@LCL\n")
		out.WriteString("M=D\n")
		// goto retAddr
		out.WriteString(frameRef(0))
		out.WriteString("A=M\n")
		out.WriteString("0;JMP\n")
		return out.String()
	}
}
