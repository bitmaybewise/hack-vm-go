package parser

import "github.com/hack-vm-go/goutils"

type ParsedLine struct {
	Raw  string
	Args []string
}

func (pl ParsedLine) IsArithmeticOrLogic() bool {
	return goutils.Includes(ArithmeticLogicOp, pl.Args[0])
}

func (pl ParsedLine) IsPushPop() bool {
	return pl.Args[0] == "push" || pl.Args[0] == "pop"
}

func (pl ParsedLine) Segment() string {
	return pl.Args[1]
}

func (pl ParsedLine) Index() string {
	return pl.Args[2]
}

func (pl ParsedLine) CommandType() string {
	return pl.Args[0]
}
