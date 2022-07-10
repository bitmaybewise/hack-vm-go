package translator

import "strings"

func EndLoop(out *strings.Builder) {
	out.WriteString("(END)\n")
	out.WriteString("@END\n")
	out.WriteString("0;JMP\n")
}
