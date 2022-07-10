package translator

import "strings"

func add() string {
	out := new(strings.Builder)
	pop2(out)
	out.WriteString("D=D+M\n")
	out.WriteString("M=D\n")
	return out.String()
}

func sub() string {
	out := new(strings.Builder)
	pop2(out)
	out.WriteString("D=M-D\n")
	out.WriteString("M=D\n")
	return out.String()
}
