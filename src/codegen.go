package main

import (
	"strings"
)

// GenBash compiles the sequence of parsed Tokens into a valid, executable Bash script string.
func GenBash(parsedItems []Token) string {
	out := []string{BashHeaders}
	var blockStack []string
	var deferStack []string

	for _, item := range parsedItems {
		kind, val := item.Type, item.Value

		switch kind {
		case "RAW_CMD", "EMPTY", "ASSIGN", "LOCAL_ASSIGN", "CONST", "DEL", "PRINTN", "PRINT", "RAW", "MATH":
			out = append(out, val)

		case "IF":
			cond := processCond(val)
			if strings.TrimSpace(cond) == "" {
				cond = "true"
			}
			out = append(out, "if [[ "+cond+" ]]; then")
			blockStack = append(blockStack, "if")

		case "ELIF":
			cond := processCond(val)
			if strings.TrimSpace(cond) == "" {
				cond = "true"
			}
			out = append(out, "elif [[ "+cond+" ]]; then")

		case "ELSE":
			out = append(out, "else")

		case "WHILE":
			out = append(out, "while [[ "+val+" ]]; do")
			blockStack = append(blockStack, "loop")

		case "UNTIL":
			out = append(out, "until [[ "+val+" ]]; do")
			blockStack = append(blockStack, "loop")

		// Handle for loop block headers
		case "FOR":
			out = append(out, "for "+val+"; do")
			blockStack = append(blockStack, "loop")

		case "DEFER":
			deferStack = append([]string{val}, deferStack...)

		case "CASE_START":
			if strings.HasPrefix(val, "\"$") {
				out = append(out, "case "+val+" in")
			} else {
				out = append(out, "case \"$"+val+"\" in")
			}
			blockStack = append(blockStack, "case")

		case "CASE_OVER":
			out = append(out, ";;")

		case "FN_START":
			out = append(out, val)
			blockStack = append(blockStack, "fn")

		case "RETURN":
			out = append(out, "    "+val)

		case "BLOCK_END":
			if len(blockStack) == 0 {
				out = append(out, "}")
				continue
			}

			lastBlock := blockStack[len(blockStack)-1]
			blockStack = blockStack[:len(blockStack)-1]

			switch lastBlock {
			case "if":
				out = append(out, "fi")
			case "loop":
				out = append(out, "done")
			case "case":
				out = append(out, "esac")
			case "fn":
				out = append(out, "}")
			}
		}
	}

	// Prepend deferred tasks as an EXIT trap if any were registered
	if len(deferStack) > 0 {
		deferString := strings.Join(deferStack, "; ")
		trapCmd := "\ntrap '" + deferString + "' EXIT\n"

		out = append(out[:1], append([]string{trapCmd}, out[1:]...)...)
	}

	return strings.Join(out, "\n")
}

// processCond normalizes condition expressions and replaces logical words with Bash operators.
func processCond(val string) string {
	val = strings.TrimSpace(val)
	val = strings.ReplaceAll(val, " and ", " && ")
	val = strings.ReplaceAll(val, " or ", " || ")

	words := strings.Fields(val)
	processed := make([]string, 0, len(words))

	for _, w := range words {
		if isIdentifier(w) && w != "true" && w != "false" {
			processed = append(processed, "\"$"+w+"\"")
		} else {
			processed = append(processed, w)
		}
	}

	return strings.Join(processed, " ")
}
