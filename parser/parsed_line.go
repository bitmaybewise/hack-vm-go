package parser

import (
	"fmt"
	"strconv"
)

type ParsedLine struct {
	Raw      string
	Filename string
	Args     []string
	idx      int
}

func (pl ParsedLine) Id() string {
	return fmt.Sprintf("%s.%s", pl.Filename, pl.Segment())
}

func (pl ParsedLine) CommandType() string {
	return pl.Args[0]
}

func (pl ParsedLine) Segment() string {
	return pl.Args[1]
}

func (pl ParsedLine) Idx() int {
	if pl.idx == -1 {
		n, err := strconv.Atoi(pl.Args[2])
		if err != nil {
			panic(err)
		}
		pl.idx = n
	}

	return pl.idx
}
