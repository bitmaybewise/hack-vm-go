package parser

import (
	"bufio"
	"errors"
	"io"
	"strings"
)

var ArithmeticLogicOp = []string{
	"add", "sub", "neg", "eq", "lt", "gt", "and", "or", "not",
}

var (
	IgnoredLine = errors.New("ignored line")

	EmptyLine = ParsedLine{}
)

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

	line = strings.ReplaceAll(line, "\r", "")
	line = strings.ReplaceAll(line, "\n", "")
	line = strings.ReplaceAll(line, "\t", "")
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
