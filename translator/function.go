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

var callCounter int

func call(out *strings.Builder, fn string, idx int) {
	ref := generateLabel(fn)
	pushLabel(out, ref)

	pop("LCL", out)
	push(out)

	pop("ARG", out)
	push(out)

	pop("THIS", out)
	push(out)

	pop("THAT", out)
	push(out)

	repositionsARG(out, idx)
	repositionsLCL(out)
	gotoFunction(out, fn)

	out.WriteString(fmt.Sprintf("(%s)\n", ref))
}

func pop(ref string, out *strings.Builder) {
	out.WriteString(fmt.Sprintf("@%s\n", ref))
	out.WriteString("D=M\n")
}

func repositionsARG(out *strings.Builder, n int) {
	out.WriteString("@SP\n")
	out.WriteString("D=M\n")
	out.WriteString("D=D-1\n")
	out.WriteString("D=D-1\n")
	out.WriteString("D=D-1\n")
	out.WriteString("D=D-1\n")
	out.WriteString("D=D-1\n")
	for i := 0; i < n; i++ {
		out.WriteString("D=D-1\n")
	}
	out.WriteString("@ARG\n")
	out.WriteString("M=D\n")
}

func repositionsLCL(out *strings.Builder) {
	out.WriteString("@SP\n")
	out.WriteString("D=M\n")
	out.WriteString("@LCL\n")
	out.WriteString("M=D\n")
}

func gotoFunction(out *strings.Builder, fn string) {
	out.WriteString(fmt.Sprintf("@%s\n", fn))
	out.WriteString("0;JMP\n")
}

func generateLabel(fn string) string {
	callCounter++
	return fmt.Sprintf("%s$ret.%d", fn, callCounter)
}

func pushLabel(out *strings.Builder, label string) {
	out.WriteString(fmt.Sprintf("@%s\n", label))
	out.WriteString("D=A\n")
	push(out)
}

func functionCall(line parser.ParsedLine) translator {
	return func() string {
		out := new(strings.Builder)
		call(out, line.Segment(), line.Idx())
		return out.String()
	}
}

func functionReturn() translator {
	return func() string {
		out := new(strings.Builder)
		frameMinus := func(frame string, n int) {
			out.WriteString(frame + "\n")
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
		out.WriteString("@R13\n")
		out.WriteString("M=D\n")
		// retAddr = *(frame-5)
		frameMinus("@R13", 5)
		out.WriteString("@R14\n")
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
		frameMinus("@R13", 1)
		out.WriteString("@THAT\n")
		out.WriteString("M=D\n")
		// THIS = *(frame-2)
		frameMinus("@R13", 2)
		out.WriteString("@THIS\n")
		out.WriteString("M=D\n")
		// ARG = *(frame-3)
		frameMinus("@R13", 3)
		out.WriteString("@ARG\n")
		out.WriteString("M=D\n")
		// LCL = *(frame-4)
		frameMinus("@R13", 4)
		out.WriteString("@LCL\n")
		out.WriteString("M=D\n")
		// goto retAddr
		out.WriteString("@R14\n")
		out.WriteString("A=M\n")
		out.WriteString("0;JMP\n")
		return out.String()
	}
}
