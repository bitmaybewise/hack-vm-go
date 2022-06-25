package translator

import (
	"strings"

	"github.com/hack-vm-go/parser"
)

func Arithmetic(cmd parser.Command) string {
	var out strings.Builder
	out.WriteString("@SP\n")
	out.WriteString("M=M-1\n")
	out.WriteString("A=M\n")
	out.WriteString("D=M\n")
	out.WriteString("@SP\n")
	out.WriteString("A=M-1\n")
	out.WriteString("D=D+M\n")
	out.WriteString("M=D\n")
	return out.String()
}

func PushPop(cmd parser.Command, segment, index string) string {
	switch cmd {
	case parser.C_PUSH:
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
	case parser.C_POP:
	}
	return ""
}

func EndLoop() string {
	var out strings.Builder
	out.WriteString("(END)\n")
	out.WriteString("\t@END\n")
	out.WriteString("\t0;JMP\n")
	return out.String()
}
