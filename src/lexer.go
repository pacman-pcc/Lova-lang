package main

import (
	"strings"
)

type Token struct {
	Type  string
	Value string
}

func Tokenize(line string) []Token {
	line = strings.TrimSpace(line)

	if line == "" || strings.HasPrefix(line, "//") {
		return []Token{{"EMPTY", ""}}
	}

	if idx := strings.Index(line, " //"); idx != -1 {
		line = strings.TrimSpace(line[:idx])
	}

	return []Token{{"SPECIAL", line}}
}

func Lex(code string) []Token {
	lines := strings.Split(code, "\n")
	var parsedLines []Token

	for _, line := range lines {
		parsedLines = append(parsedLines, Tokenize(line)...)
	}

	return parsedLines
}
