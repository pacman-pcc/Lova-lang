package main

import (
	"regexp"
	"strings"
	"unicode"
)

var (
	// reInter handles string interpolation:
	// Matches array element access like [arrName][index], [arrName][!], or standard variables {varName}.
	reInter   = regexp.MustCompile(`\[([a-zA-Z0-9_]+)\]\[([a-zA-Z0-9_!\*@]+)\]|\{([a-zA-Z0-9_\.\[\]\?]+)\}`)
	reWordAnd = regexp.MustCompile(`\band\b`)
	reWordOr  = regexp.MustCompile(`\bor\b`)
	reWordNot = regexp.MustCompile(`\bnot\b`)
)

// isIdentifier checks if a string is a valid variable identifier.
func isIdentifier(s string) bool {
	if len(s) == 0 {
		return false
	}
	for i, r := range s {
		if i == 0 {
			if !unicode.IsLetter(r) && r != '_' {
				return false
			}
		} else {
			if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' {
				return false
			}
		}
	}
	return true
}

// isDigit checks if a string consists entirely of digits.
func isDigit(s string) bool {
	if len(s) == 0 {
		return false
	}
	for _, r := range s {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

// parseStringInter translates interpolation patterns in strings into standard Bash syntax.
// Supports [arr][idx], [arr][!] (all elements), and {var}.
func parseStringInter(text string) string {
	return reInter.ReplaceAllStringFunc(text, func(m string) string {
		// Handle array indexing syntax: [name][num] or [name][!]
		if strings.HasPrefix(m, "[") {
			parts := strings.Split(m, "][")
			if len(parts) == 2 {
				arr := strings.TrimPrefix(parts[0], "[")
				idx := strings.TrimSuffix(parts[1], "]")
				arr = strings.ReplaceAll(arr, ".", "_")

				// Map '!' to '@' for Bash array expansion (all elements)
				if idx == "!" {
					idx = "@"
				}
				return "${" + arr + "[" + idx + "]}"
			}
		}

		// Handle standard variable interpolation: {var}
		match := reInter.FindStringSubmatch(m)
		if len(match) > 3 && match[3] != "" {
			varName := strings.ReplaceAll(strings.TrimSpace(match[3]), ".", "_")
			return "${" + varName + "}"
		}
		return m
	})
}

// parseCondition translates logical expressions and built-in checks into Bash-compatible syntax.
func parseCondition(cond string) string {
	cond = strings.TrimSpace(cond)
	if cond == "" {
		return ""
	}

	// File and directory check prefixes mapping to Bash test flags
	prefixes := []struct {
		prefix string
		flag   string
	}{
		{"is_file ", "-f"},
		{"is_dir ", "-d"},
		{"is_exist ", "-e"},
	}

	for _, p := range prefixes {
		if strings.HasPrefix(cond, p.prefix) {
			path := strings.TrimSpace(cond[len(p.prefix):])
			path = strings.Trim(path, "\"'")
			path = parseStringInter(path)
			return p.flag + " \"" + path + "\""
		}
	}

	// Replace logical operators with Bash equivalents
	cond = reWordAnd.ReplaceAllString(cond, "&&")
	cond = reWordOr.ReplaceAllString(cond, "||")
	cond = reWordNot.ReplaceAllString(cond, "!")

	words := strings.Fields(cond)
	processed := make([]string, 0, len(words))

	for _, w := range words {
		// Keep quoted strings intact
		if (strings.HasPrefix(w, "\"") && strings.HasSuffix(w, "\"")) ||
			(strings.HasPrefix(w, "'") && strings.HasSuffix(w, "'")) {
			processed = append(processed, w)
			continue
		}

		wClean := strings.ReplaceAll(w, ".", "_")

		switch wClean {
		case "true", "false", "&&", "||", "==", "!=", ">=", "<=", ">", "<", "!":
			processed = append(processed, w)
		default:
			if isDigit(wClean) {
				processed = append(processed, w)
			} else if isIdentifier(wClean) {
				processed = append(processed, "\"$"+wClean+"\"")
			} else if strings.HasPrefix(w, "$") {
				if !strings.HasPrefix(w, "\"$") {
					processed = append(processed, "\""+w+"\"")
				} else {
					processed = append(processed, w)
				}
			} else {
				processed = append(processed, w)
			}
		}
	}

	return strings.Join(processed, " ")
}

// ParseLine processes a single Token and translates it into the appropriate target Bash command structure.
func ParseLine(item Token) Token {
	kind, line := item.Type, item.Value

	if kind == "EMPTY" {
		return Token{Type: "EMPTY", Value: ""}
	}

	rawLine := strings.TrimSpace(line)

	// Handle labels (e.g., case patterns or goto targets)
	if strings.HasSuffix(rawLine, ":") && !strings.HasPrefix(rawLine, "proc") && !strings.HasPrefix(rawLine, "procloc") {
		label := strings.TrimSuffix(rawLine, ":")
		label = strings.TrimSpace(label)
		label = parseStringInter(label)
		if label == "_" {
			return Token{Type: "RAW", Value: "*)"}
		}
		return Token{Type: "RAW", Value: label + ")"}
	}

	// 1. Handle Array Declaration: arr [name] = [1, 2, 3] or arr [name] = ["A", "B"]
	if strings.HasPrefix(rawLine, "arr ") {
		content := strings.TrimSpace(rawLine[4:])
		if strings.Contains(content, "=") {
			parts := strings.SplitN(content, "=", 2)
			varRaw := strings.TrimSpace(parts[0])

			// Extract name from brackets if syntax uses arr [name] = [...]
			varName := strings.TrimPrefix(varRaw, "[")
			varName = strings.TrimSuffix(varName, "]")
			varName = strings.TrimSpace(varName)
			varName = strings.ReplaceAll(varName, ".", " ")

			valPart := strings.TrimSpace(parts[1])
			// Convert brackets [1, 2, 3] to Bash array syntax (1 2 3)
			valPart = strings.TrimPrefix(valPart, "[")
			valPart = strings.TrimSuffix(valPart, "]")

			// Interpolate strings inside array values if necessary
			valPart = parseStringInter(valPart)

			return Token{Type: "ASSIGN", Value: varName + "=(" + valPart + ")"}
		}
	}

	// 2. Handle Array Append Method: [name].append(val1, val2) -> name+=(val1 val2)
	if strings.Contains(rawLine, ".append(") && strings.HasSuffix(rawLine, ")") {
		idxOpen := strings.Index(rawLine, "[")
		idxClose := strings.Index(rawLine, "]")
		idxDot := strings.Index(rawLine, ".")

		if idxOpen == 0 && idxClose > idxOpen && idxDot == idxClose+1 {
			arrName := rawLine[idxOpen+1 : idxClose]
			arrName = strings.ReplaceAll(arrName, ".", " ")

			contentStart := strings.Index(rawLine, "(")
			contentEnd := strings.LastIndex(rawLine, ")")
			args := rawLine[contentStart+1 : contentEnd]
			args = parseStringInter(args)

			return Token{Type: "RAW", Value: arrName + "+=(" + args + ")"}
		}
	}

	// 3. Handle Array Delete Method: [name].delete(indexOrName) -> unset 'name[index]'
	if strings.Contains(rawLine, ".delete(") && strings.HasSuffix(rawLine, ")") {
		idxOpen := strings.Index(rawLine, "[")
		idxClose := strings.Index(rawLine, "]")
		idxDot := strings.Index(rawLine, ".")

		if idxOpen == 0 && idxClose > idxOpen && idxDot == idxClose+1 {
			arrName := rawLine[idxOpen+1 : idxClose]
			arrName = strings.ReplaceAll(arrName, ".", "_")

			contentStart := strings.Index(rawLine, "(")
			contentEnd := strings.LastIndex(rawLine, ")")
			indexArg := strings.TrimSpace(rawLine[contentStart+1 : contentEnd])
			indexArg = parseStringInter(indexArg)

			return Token{Type: "RAW", Value: "unset '" + arrName + "[" + indexArg + "]'"}
		}
	}

	// Handle arithmetic expressions with standard operators
	if strings.Contains(rawLine, "=") && (strings.Contains(rawLine, "+") || strings.Contains(rawLine, "-") || strings.Contains(rawLine, "*") || strings.Contains(rawLine, "/")) &&
		!strings.HasPrefix(rawLine, "proc ") && !strings.HasPrefix(rawLine, "procloc ") && !strings.HasPrefix(rawLine, "const ") {
		parts := strings.SplitN(rawLine, "=", 2)
		varName := strings.TrimSpace(parts[0])
		expr := strings.TrimSpace(parts[1])
		varName = strings.ReplaceAll(varName, ".", "_")
		expr = parseStringInter(expr)
		return Token{Type: "MATH", Value: varName + "=$((" + expr + "))"}
	}

	// Handle standard variable assignment
	if strings.HasPrefix(rawLine, "proc ") {
		content := parseStringInter(strings.TrimSpace(rawLine[5:]))
		if strings.Contains(content, "=") {
			parts := strings.SplitN(content, "=", 2)
			varName := strings.ReplaceAll(strings.TrimSpace(parts[0]), ".", "_")
			valStr := strings.TrimSpace(parts[1])
			return Token{Type: "ASSIGN", Value: varName + "=" + valStr}
		}
		return Token{Type: "ASSIGN", Value: content}
	}

	// Handle local variable assignment inside functions
	if strings.HasPrefix(rawLine, "procloc ") {
		content := parseStringInter(strings.TrimSpace(rawLine[8:]))
		if strings.Contains(content, "=") {
			parts := strings.SplitN(content, "=", 2)
			varName := strings.ReplaceAll(strings.TrimSpace(parts[0]), ".", "_")
			valStr := strings.TrimSpace(parts[1])
			return Token{Type: "LOCAL_ASSIGN", Value: "local " + varName + "=" + valStr}
		}
		return Token{Type: "LOCAL_ASSIGN", Value: "local " + content}
	}

	// Handle constant declarations
	if strings.HasPrefix(rawLine, "const ") {
		content := parseStringInter(strings.TrimSpace(rawLine[6:]))
		if strings.Contains(content, "=") {
			parts := strings.SplitN(content, "=", 2)
			varName := strings.ReplaceAll(strings.TrimSpace(parts[0]), ".", "_")
			valStr := strings.TrimSpace(parts[1])
			return Token{Type: "CONST", Value: "readonly " + varName + "=" + valStr}
		}
		return Token{Type: "CONST", Value: "readonly " + content}
	}

	// Handle variable deletion (unset)
	if strings.HasPrefix(rawLine, "del ") {
		varName := strings.ReplaceAll(strings.TrimSpace(rawLine[4:]), ".", "_")
		return Token{Type: "DEL", Value: "unset " + varName}
	}

	// Handle printing commands (print and printn)
	if strings.HasPrefix(rawLine, "printn ") || strings.HasPrefix(rawLine, "print ") {
		isNet := strings.HasPrefix(rawLine, "printn ")
		prefixLen := 6
		if isNet {
			prefixLen = 7
		}
		val := strings.TrimSpace(rawLine[prefixLen:])

		// If printing array elements or variables with formatting, parse interpolation first
		valInterp := parseStringInter(val)

		cmd := "PRINT"
		fmtStr := "%b"
		if isNet {
			cmd = "PRINTN"
			fmtStr = "%b\\n"
		}

		// Support passing array expansion directly without extra double quotes wrapping if it's an interpolation token
		if strings.HasPrefix(valInterp, "${") && strings.HasSuffix(valInterp, "}") {
			return Token{Type: cmd, Value: "printf \"" + fmtStr + "\" \"" + valInterp + "\""}
		}

		valTrimmed := strings.Trim(val, "\"'")
		finalVal := parseStringInter(valTrimmed)
		return Token{Type: cmd, Value: "printf \"" + fmtStr + "\" \"" + finalVal + "\""}
	}

	lineClean := parseStringInter(rawLine)

	// Handle conditional statements (if)
	if strings.HasPrefix(lineClean, "if ") && strings.HasSuffix(lineClean, "{") {
		cond := strings.TrimSpace(lineClean[3 : len(lineClean)-1])
		return Token{Type: "IF", Value: parseCondition(cond)}
	}

	checkLine := lineClean
	if strings.HasPrefix(checkLine, "}") {
		checkLine = strings.TrimSpace(checkLine[1:])
	}

	// Handle alternative branches (elif / else if)
	if (strings.HasPrefix(checkLine, "elseif ") || strings.HasPrefix(checkLine, "else if ")) && strings.HasSuffix(checkLine, "{") {
		idx := strings.Index(checkLine, "if ")
		cond := strings.TrimSpace(checkLine[idx+3 : len(checkLine)-1])
		return Token{Type: "ELIF", Value: parseCondition(cond)}
	}

	if checkLine == "else {" {
		return Token{Type: "ELSE", Value: ""}
	}

	// Handle loops (while)
	if strings.HasPrefix(lineClean, "while ") && strings.HasSuffix(lineClean, "{") {
		cond := strings.TrimSpace(lineClean[6 : len(lineClean)-1])
		return Token{Type: "WHILE", Value: parseCondition(cond)}
	}

	if strings.HasPrefix(lineClean, "for ") && strings.HasSuffix(lineClean, "{") {
			content := strings.TrimSpace(lineClean[3 : len(lineClean)-1])

			parts := strings.Split(content, " in ")
			if len(parts) == 2 {
				varName := strings.TrimSpace(parts[0])
				varName = strings.ReplaceAll(varName, ".", "_")

				iterable := strings.TrimSpace(parts[1])
				if strings.HasPrefix(iterable, "[") && strings.HasSuffix(iterable, "]") {
					arrName := strings.Trim(iterable, "[]")
					arrName = strings.ReplaceAll(arrName, ".", "_")
					iterable = "\"${" + arrName + "[@]}\""
				} else {
					iterable = parseStringInter(iterable)
				}

				return Token{Type: "FOR", Value: varName + " in " + iterable}
			}

			cExpr := content
			if !strings.HasPrefix(cExpr, "(") {
				cExpr = "((" + cExpr + "))"
			} else if strings.HasPrefix(cExpr, "(") && !strings.HasPrefix(cExpr, "((") {
				cExpr = "(" + cExpr + ")"
			}

			return Token{Type: "FOR", Value: cExpr}
		}

	// Handle until loops
	if strings.HasPrefix(lineClean, "fdo ") && strings.HasSuffix(lineClean, "{") {
		cond := strings.TrimSpace(lineClean[4 : len(lineClean)-1])
		return Token{Type: "UNTIL", Value: parseCondition(cond)}
	}

	// Handle case matching statements
	if strings.HasPrefix(lineClean, "case ") && strings.Contains(lineClean, "do {") {
		varStr := strings.ReplaceAll(lineClean, "case ", "")
		varStr = strings.ReplaceAll(varStr, "do {", "")
		varStr = strings.TrimSpace(varStr)
		varStr = strings.ReplaceAll(varStr, ".", "_")
		if !strings.HasPrefix(varStr, "$") && !strings.HasPrefix(varStr, "\"") {
			varStr = "\"$" + varStr + "\""
		}
		return Token{Type: "CASE_START", Value: varStr}
	}

	// Handle deferred execution (traps)
	if strings.HasPrefix(rawLine, "defer ") {
		innerCmd := strings.TrimSpace(rawLine[6:])
		innerToken := ParseLine(Token{Type: "SPECIAL", Value: innerCmd})
		return Token{Type: "DEFER", Value: innerToken.Value}
	}

	// Handle function declarations
	if strings.HasPrefix(lineClean, "function ") && strings.HasSuffix(lineClean, "{") {
		header := strings.TrimSpace(lineClean[9 : len(lineClean)-1])
		if idx := strings.Index(header, "("); idx != -1 {
			header = strings.TrimSpace(header[:idx])
		}
		return Token{Type: "FN_START", Value: header + " () {"}
	}

	// Handle return statements inside functions
	if strings.HasPrefix(lineClean, "return ") {
		code := strings.TrimSpace(lineClean[7:])
		return Token{Type: "RETURN", Value: "return " + code}
	}

	if lineClean == "over" {
		return Token{Type: "CASE_OVER", Value: ";;"}
	}

	if lineClean == "}" {
		return Token{Type: "BLOCK_END", Value: "}"}
	}

	return Token{Type: "RAW", Value: lineClean}
}
