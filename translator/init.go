package translator

import "strings"

func Init(out *strings.Builder) {
	out.WriteString("// SP = 256\n")
	out.WriteString("@256\n")
	out.WriteString("D=A\n")
	out.WriteString("@SP\n")
	out.WriteString("M=D\n")

	out.WriteString("// call Sys.init\n")
	call(out, "Sys.init", 0)
	out.WriteString("(END)\n")
	out.WriteString("@END\n")
	out.WriteString("0;JMP\n")
}
