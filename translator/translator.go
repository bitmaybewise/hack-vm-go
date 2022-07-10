package translator

import (
	"fmt"
	"strconv"
	"strings"
)

const Temp0 = 5 // Temp0 starts at RAM[5]

var symbolCounter int

func ArithmeticLogic(cmd string) string {
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

func PushPop(cmd string, segment, index, filename string) string {
	out := new(strings.Builder)

	switch cmd {
	case "push":
		pushTo(out, segment, index, filename)
	case "pop":
		popTo(out, segment, index, filename)
	}

	return out.String()
}

func pushTo(out *strings.Builder, segment, index, filename string) {
	n, err := strconv.Atoi(index)
	if err != nil {
		panic(err)
	}
	push := func() {
		out.WriteString("@SP\n")
		out.WriteString("A=M\n")
		out.WriteString("M=D\n")
		out.WriteString("@SP\n")
		out.WriteString("M=M+1\n")
	}
	pop := func(at string) {
		out.WriteString("@" + at + "\n")
		for i := 0; i < n; i++ {
			out.WriteString("M=M+1\n")
		}
		out.WriteString("A=M\n")
		out.WriteString("D=M\n")
		if n > 0 {
			out.WriteString("@" + at + "\n")
		}
		for i := 0; i < n; i++ {
			out.WriteString("M=M-1\n")
		}
	}
	if segment == "constant" {
		out.WriteString("@" + index + "\n")
		out.WriteString("D=A\n")
		push()
	}
	if segment == "local" {
		pop("LCL")
		push()
	}
	if segment == "this" {
		pop("THIS")
		push()
	}
	if segment == "that" {
		pop("THAT")
		push()
	}
	if segment == "pointer" && n == 0 {
		out.WriteString("@THIS\n")
		out.WriteString("D=M\n")
		push()
	}
	if segment == "pointer" && n == 1 {
		out.WriteString("@THAT\n")
		out.WriteString("D=M\n")
		push()
	}
	if segment == "argument" {
		pop("ARG")
		push()
	}
	if segment == "temp" {
		out.WriteString(fmt.Sprintf("@R%d\n", Temp0+n))
		out.WriteString("D=M\n")
		push()
	}
	if segment == "static" {
		out.WriteString(fmt.Sprintf("@%s.%d\n", filename, n))
		out.WriteString("D=M\n")
		push()
	}
}

func popTo(out *strings.Builder, segment, index, filename string) {
	n, err := strconv.Atoi(index)
	if err != nil {
		panic(err)
	}
	pop := func(at string) {
		pop1(out)
		out.WriteString("@" + at + "\n")
		for i := 0; i < n; i++ {
			out.WriteString("M=M+1\n")
		}
		out.WriteString("A=M\n")
		out.WriteString("M=D\n")
		for i := 0; i < n; i++ {
			out.WriteString("@" + at + "\n")
			out.WriteString("M=M-1\n")
		}
	}
	if segment == "local" {
		pop1(out)
		out.WriteString("@LCL\n")
		out.WriteString("A=M\n")
		out.WriteString("M=D\n")
	}
	if segment == "argument" {
		pop("ARG")
	}
	if segment == "this" {
		pop("THIS")
	}
	if segment == "that" {
		pop("THAT")
	}
	if segment == "pointer" && n == 0 {
		pop1(out)
		out.WriteString("@THIS\n")
		out.WriteString("M=D\n")
	}
	if segment == "pointer" && n == 1 {
		pop1(out)
		out.WriteString("@THAT\n")
		out.WriteString("M=D\n")
	}
	if segment == "temp" {
		pop1(out)
		out.WriteString(fmt.Sprintf("@R%d\n", Temp0+n))
		out.WriteString("M=D\n")
	}
	if segment == "static" {
		pop1(out)
		out.WriteString(fmt.Sprintf("@%s.%d\n", filename, n))
		out.WriteString("M=D\n")
	}
}

func EndLoop(out *strings.Builder) {
	out.WriteString("(END)\n")
	out.WriteString("@END\n")
	out.WriteString("0;JMP\n")
}
