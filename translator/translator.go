package translator

import (
	"fmt"
	"strings"
)

var symbolCounter int

func Arithmetic(cmd string) string {
	var out strings.Builder

	switch cmd {
	case "add":
		out.WriteString(pop2())
		out.WriteString("D=D+M\n")
		out.WriteString("M=D\n")
	case "sub":
		out.WriteString(pop2())
		out.WriteString("D=M-D\n")
		out.WriteString("M=D\n")
	case "neg":
		out.WriteString(pop1())
		out.WriteString("@SP\n")
		out.WriteString("A=M\n")
		out.WriteString("M=-M\n")
		out.WriteString("@SP\n")
		out.WriteString("M=M+1\n")
		out.WriteString("A=M\n")
	case "not":
		out.WriteString(pop1())
		out.WriteString("@SP\n")
		out.WriteString("A=M\n")
		out.WriteString("M=!M\n")
		out.WriteString("@SP\n")
		out.WriteString("M=M+1\n")
		out.WriteString("A=M\n")
	case "and":
		out.WriteString(pop2())
		out.WriteString("D=D&M\n")
		out.WriteString("M=D\n")
	case "or":
		out.WriteString(pop2())
		out.WriteString("D=D|M\n")
		out.WriteString("M=D\n")
	case "eq":
		symbolCounter++
		out.WriteString(pop2())
		out.WriteString("D=D-M\n")
		out.WriteString(fmt.Sprintf("@EQ_%d\n", symbolCounter))
		out.WriteString("D;JEQ\n")
		out.WriteString("@SP\n")
		out.WriteString("A=M-1\n")
		out.WriteString("M=0\n")
		out.WriteString(fmt.Sprintf("@EQ_%d_END\n", symbolCounter))
		out.WriteString("D;JMP\n")
		out.WriteString(fmt.Sprintf("(EQ_%d)\n", symbolCounter))
		out.WriteString("@SP\n")
		out.WriteString("A=M-1\n")
		out.WriteString("M=-1\n")
		out.WriteString(fmt.Sprintf("(EQ_%d_END)\n", symbolCounter))
	case "lt":
		symbolCounter++
		out.WriteString(pop2())
		out.WriteString("D=D-M\n")
		out.WriteString(fmt.Sprintf("@LT_%d\n", symbolCounter))
		out.WriteString("D;JLE\n")
		out.WriteString("@SP\n")
		out.WriteString("A=M-1\n")
		out.WriteString("M=-1\n")
		out.WriteString(fmt.Sprintf("@LT_%d_END\n", symbolCounter))
		out.WriteString("D;JMP\n")
		out.WriteString(fmt.Sprintf("(LT_%d)\n", symbolCounter))
		out.WriteString("@SP\n")
		out.WriteString("A=M-1\n")
		out.WriteString("M=0\n")
		out.WriteString(fmt.Sprintf("(LT_%d_END)\n", symbolCounter))
	case "gt":
		symbolCounter++
		// out.WriteString(conditional(cmd, "JGT", false))
		out.WriteString(pop2())
		out.WriteString("D=D-M\n")
		out.WriteString(fmt.Sprintf("@GT_%d\n", symbolCounter))
		out.WriteString("D;JGE\n")
		out.WriteString("@SP\n")
		out.WriteString("A=M-1\n")
		out.WriteString("M=-1\n")
		out.WriteString(fmt.Sprintf("@GT_%d_END\n", symbolCounter))
		out.WriteString("D;JMP\n")
		out.WriteString(fmt.Sprintf("(GT_%d)\n", symbolCounter))
		out.WriteString("@SP\n")
		out.WriteString("A=M-1\n")
		out.WriteString("M=0\n")
		out.WriteString(fmt.Sprintf("(GT_%d_END)\n", symbolCounter))
	}

	return out.String()
}

func pop1() string {
	var out strings.Builder
	out.WriteString("@SP\n")
	out.WriteString("M=M-1\n")
	out.WriteString("A=M\n")
	out.WriteString("D=M\n")
	return out.String()
}
func pop2() string {
	var out strings.Builder
	out.WriteString(pop1())
	out.WriteString("@SP\n")
	out.WriteString("A=M-1\n")
	return out.String()
}

func conditional(cmd, asmCond string, invert bool) string {
	var out strings.Builder
	out.WriteString(pop2())
	out.WriteString("D=D-M\n")
	out.WriteString(fmt.Sprintf("@%s_%d\n", cmd, symbolCounter))
	out.WriteString(fmt.Sprintf("D;%s\n", asmCond))
	out.WriteString("@SP\n")
	out.WriteString("A=M-1\n")
	if invert {
		out.WriteString("M=-1\n")
	} else {
		out.WriteString("M=0\n")
	}
	out.WriteString(fmt.Sprintf("@%s_%d_END\n", cmd, symbolCounter))
	out.WriteString("D;JMP\n")
	out.WriteString(fmt.Sprintf("(%s_%d)\n", cmd, symbolCounter))
	out.WriteString("@SP\n")
	out.WriteString("A=M-1\n")
	if invert {
		out.WriteString("M=0\n")
	} else {
		out.WriteString("M=-1\n")
	}
	out.WriteString(fmt.Sprintf("(%s_%d_END)\n", cmd, symbolCounter))
	return out.String()
}

func PushPop(cmd string, segment, index string) string {
	switch cmd {
	case "push":
		if segment == "constant" {
			var out strings.Builder
			out.WriteString("@" + index + "\n")
			out.WriteString("D=A\n")
			out.WriteString("@SP\n")
			out.WriteString("A=M\n")
			out.WriteString("M=D\n")
			out.WriteString("@SP\n")
			out.WriteString("M=M+1\n")
			return out.String()
		}
	case "pop":
		return pop1()
	}

	return ""
}

func Setup() string {
	var out strings.Builder
	out.WriteString("@256\n")
	out.WriteString("D=A\n")
	out.WriteString("@SP\n")
	out.WriteString("M=D\n")
	out.WriteString("A=D\n")
	return out.String()
}
func EndLoop() string {
	var out strings.Builder
	out.WriteString("(END)\n")
	out.WriteString("@END\n")
	out.WriteString("0;JMP\n")
	return out.String()
}
