package parser

import (
	"bufio"
	"errors"
	"io"
	"strings"

	"github.com/hack-vm-go/goutils"
)

type Command int

const (
	C_ARITHMETIC = Command(iota)
	C_PUSH
	C_POP
	C_LABEL
	C_GOTO
	C_IF
	C_FUNCTION
	C_RETURN
	C_CALL
)

var ArithmeticOp = []string{
	"add", "sub", "neg", "eq", "lt", "gt", "and", "or", "not",
}

var IgnoredLine = errors.New("ignored line")

type ParsedLine struct {
	Raw  string
	Args []string
}

func (pl ParsedLine) IsArithmetic() bool {
	return goutils.Includes(ArithmeticOp, pl.Args[0])
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

var EmptyLine = ParsedLine{}

type Parser struct {
	input *bufio.Reader
}

func (p *Parser) ReadLine() (ParsedLine, error) {
	line, err := p.input.ReadString('\n')
	if err != nil {
		return EmptyLine, err
	}

	// removing comment
	if strings.HasPrefix(line, "//") {
		return EmptyLine, IgnoredLine
	}
	commentFoundAt := strings.Index(line, "//")
	if commentFoundAt > 1 {
		line = line[:commentFoundAt-1]
	}

	line = strings.Replace(line, "\r", "", 1)
	line = strings.Replace(line, "\n", "", 1)
	line = strings.Trim(line, " ")
	if line == "" {
		return EmptyLine, IgnoredLine
	}

	args := strings.Split(line, " ")
	return ParsedLine{Raw: line, Args: args}, nil
}

func New(input io.Reader) Parser {
	reader := bufio.NewReader(input)
	return Parser{reader}
}
