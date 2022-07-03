package translator

import (
	"fmt"
	"strings"
)

var symbolCounter int

func Arithmetic(cmd string) string {
	out := new(strings.Builder)

	switch cmd {
	case "add":
		pop2(out)
		out.WriteString("D=D+M\n")
		out.WriteString("M=D\n")
	case "sub":
		pop2(out)
		out.WriteString("D=M-D\n")
		out.WriteString("M=D\n")
	case "neg":
		pop1(out)
		out.WriteString("M=-M\n")
		spForward(out)
	case "not":
		pop1(out)
		out.WriteString("M=!M\n")
		spForward(out)
	case "and":
		pop2(out)
		out.WriteString("D=D&M\n")
		out.WriteString("M=D\n")
	case "or":
		pop2(out)
		out.WriteString("D=D|M\n")
		out.WriteString("M=D\n")
	case "eq":
		condJump(out, cmd, "JEQ")
	case "lt":
		condJump(out, cmd, "JGT")
	case "gt":
		condJump(out, cmd, "JLT")
	}

	return out.String()
}

func pop1(out *strings.Builder) {
	out.WriteString("@SP\n")
	out.WriteString("M=M-1\n")
	out.WriteString("A=M\n")
	out.WriteString("D=M\n")
}
func pop2(out *strings.Builder) {
	pop1(out)
	out.WriteString("@SP\n")
	out.WriteString("A=M-1\n")
}

func spForward(out *strings.Builder) {
	out.WriteString("@SP\n")
	out.WriteString("M=M+1\n")
	out.WriteString("A=M\n")
}

func condJump(out *strings.Builder, cmd, asmCond string) {
	symbolCounter++

	pop2(out)
	out.WriteString("D=D-M\n")
	out.WriteString(fmt.Sprintf("@%s_%d\n", cmd, symbolCounter))
	out.WriteString(fmt.Sprintf("D;%s\n", asmCond))
	out.WriteString("@SP\n")
	out.WriteString("A=M-1\n")
	out.WriteString("M=0\n")
	out.WriteString(fmt.Sprintf("@%s_%d_END\n", cmd, symbolCounter))
	out.WriteString("D;JMP\n")
	out.WriteString(fmt.Sprintf("(%s_%d)\n", cmd, symbolCounter))
	out.WriteString("@SP\n")
	out.WriteString("A=M-1\n")
	out.WriteString("M=-1\n")
	out.WriteString(fmt.Sprintf("(%s_%d_END)\n", cmd, symbolCounter))
}

func PushPop(cmd string, segment, index string) string {
	out := new(strings.Builder)

	switch cmd {
	case "push":
		if segment == "constant" {
			out.WriteString("@" + index + "\n")
			out.WriteString("D=A\n")
			out.WriteString("@SP\n")
			out.WriteString("A=M\n")
			out.WriteString("M=D\n")
			out.WriteString("@SP\n")
			out.WriteString("M=M+1\n")
		}
	case "pop":
		pop1(out)
	}

	return out.String()
}

func EndLoop(out *strings.Builder) {
	out.WriteString("(END)\n")
	out.WriteString("@END\n")
	out.WriteString("0;JMP\n")
}
