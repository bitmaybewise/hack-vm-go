package translator

import (
	"fmt"
	"strings"

	"github.com/hack-vm-go/parser"
)

const Temp0 = 5 // Temp0 starts at RAM[5]

type translator func() string

func ToAsm(line parser.ParsedLine) string {
	var commandMapping = map[string]translator{
		"add":     add,
		"sub":     sub,
		"neg":     neg,
		"not":     not,
		"and":     and,
		"or":      or,
		"eq":      eq,
		"lt":      lt,
		"gt":      gt,
		"push":    pushTo(line),
		"pop":     popTo(line),
		"label":   label(line),
		"goto":    goTo(line),
		"if-goto": ifGoTo(line),
	}

	fn, ok := commandMapping[line.CommandType()]
	if !ok {
		panic(fmt.Sprintf("command %q is missing implementation", line.CommandType()))
	}

	return fn()
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
