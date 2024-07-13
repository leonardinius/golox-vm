package vmcompiler

import "github.com/leonardinius/goloxvm/internal/vmscanner"

type Parser struct {
	current  vmscanner.Token
	previous vmscanner.Token
	hadError bool
	// panicMode bool
}

func NewParser() Parser {
	return Parser{}
}
