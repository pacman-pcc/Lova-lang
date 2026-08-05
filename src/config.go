package main

const BashHeaders = `#!/bin/bash
set -euo pipefail
IFS=$'\n\t'
`

var Operators = map[string]string{
	"==":  "-eq",
	"!=":  "-ne",
	">":   "-gt",
	"<":   "-lt",
	">=":  "-ge",
	"<=":  "-le",
	"and": "&&",
	"or":  "||",
}

var Keywords = map[string]bool{
	"proc":     true,
	"procloc":  true,
	"const":    true,
	"del":      true,
	"printn":   true,
	"print":    true,
	"if":       true,
	"elseif":   true,
	"else":     true,
	"while":    true,
	"fdo":      true,
	"case":     true,
	"over":     true,
	"function": true,
	"return":   true,
	"defer":    true,
}
