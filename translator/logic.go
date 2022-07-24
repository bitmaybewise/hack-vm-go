package translator

import (
	"fmt"
	"strings"
)

func neg() string {
	out := new(strings.Builder)
	pop1(out)
	out.WriteString("M=-M\n")
	stackPtrForward(out)
	return out.String()
}

func not() string {
	out := new(strings.Builder)
	pop1(out)
	out.WriteString("M=!M\n")
	stackPtrForward(out)
	return out.String()
}

func and() string {
	out := new(strings.Builder)
	pop2(out)
	out.WriteString("D=D&M\n")
	out.WriteString("M=D\n")
	return out.String()
}

func or() string {
	out := new(strings.Builder)
	pop2(out)
	out.WriteString("D=D|M\n")
	out.WriteString("M=D\n")
	return out.String()
}

func eq() string {
	out := new(strings.Builder)
	condJump(out, "eq", "JEQ")
	return out.String()
}

func lt() string {
	out := new(strings.Builder)
	condJump(out, "lt", "JLT")
	return out.String()
}

func gt() string {
	out := new(strings.Builder)
	condJump(out, "gt", "JGT")
	return out.String()
}

var symbolCounter int

func condJump(out *strings.Builder, cmd, asmCond string) {
	symbolCounter++

	pop2(out)
	out.WriteString("D=M-D\n")
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

func stackPtrForward(out *strings.Builder) {
	out.WriteString("@SP\n")
	out.WriteString("M=M+1\n")
	out.WriteString("A=M\n")
}
