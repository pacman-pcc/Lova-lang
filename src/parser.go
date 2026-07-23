package main

import (
	"regexp"
	"strings"
	"unicode"
)

var (
	reInter   = regexp.MustCompile(`\{([a-zA-Z0-9_\.\[\]\?]+)\}`)
	reWordAnd = regexp.MustCompile(`\band\b`)
	reWordOr  = regexp.MustCompile(`\bor\b`)
	reWordNot = regexp.MustCompile(`\bnot\b`)
)

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

func parseStringInter(text string) string {
	return reInter.ReplaceAllStringFunc(text, func(m string) string {
		match := reInter.FindStringSubmatch(m)
		if len(match) > 1 {
			varName := strings.ReplaceAll(strings.TrimSpace(match[1]), ".", "_")
			return "${" + varName + "}"
		}
		return m
	})
}

func parseCondition(cond string) string {
	cond = strings.TrimSpace(cond)
	if cond == "" {
		return ""
	}

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

	cond = reWordAnd.ReplaceAllString(cond, "&&")
	cond = reWordOr.ReplaceAllString(cond, "||")
	cond = reWordNot.ReplaceAllString(cond, "!")

	words := strings.Fields(cond)
	processed := make([]string, 0, len(words))

	for _, w := range words {
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

func ParseLine(item Token) Token {
	kind, line := item.Type, item.Value

	if kind == "EMPTY" {
		return Token{Type: "EMPTY", Value: ""}
	}

	rawLine := strings.TrimSpace(line)

	// ПРИОРИТЕТ 1: Ветки case (например "init":, "build" | "compile": или _:)
	if strings.HasSuffix(rawLine, ":") && !strings.HasPrefix(rawLine, "proc") && !strings.HasPrefix(rawLine, "procloc") {
		label := strings.TrimSuffix(rawLine, ":")
		label = strings.TrimSpace(label)
		label = parseStringInter(label)
		if label == "_" {
			return Token{Type: "RAW", Value: "*)"}
		}
		return Token{Type: "RAW", Value: label + ")"}
	}

	// ПРИОРИТЕТ 2: Математика (например counter = counter + 1)
	if strings.Contains(rawLine, "=") && (strings.Contains(rawLine, "+") || strings.Contains(rawLine, "-") || strings.Contains(rawLine, "*") || strings.Contains(rawLine, "/")) &&
		!strings.HasPrefix(rawLine, "proc ") && !strings.HasPrefix(rawLine, "procloc ") && !strings.HasPrefix(rawLine, "const ") {
		parts := strings.SplitN(rawLine, "=", 2)
		varName := strings.TrimSpace(parts[0])
		expr := strings.TrimSpace(parts[1])
		varName = strings.ReplaceAll(varName, ".", "_")
		expr = parseStringInter(expr)
		return Token{Type: "MATH", Value: varName + "=$((" + expr + "))"}
	}

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

	if strings.HasPrefix(rawLine, "del ") {
		varName := strings.ReplaceAll(strings.TrimSpace(rawLine[4:]), ".", "_")
		return Token{Type: "DEL", Value: "unset " + varName}
	}

	if strings.HasPrefix(rawLine, "printn ") || strings.HasPrefix(rawLine, "print ") {
		isNet := strings.HasPrefix(rawLine, "printn ")
		prefixLen := 6
		if isNet {
			prefixLen = 7
		}
		val := strings.TrimSpace(rawLine[prefixLen:])
		val = strings.Trim(val, "\"'")
		val = parseStringInter(val)

		cmd := "PRINT"
		fmtStr := "%b"
		if isNet {
			cmd = "PRINTN"
			fmtStr = "%b\\n"
		}
		return Token{Type: cmd, Value: "printf \"" + fmtStr + "\" \"" + val + "\""}
	}

	lineClean := parseStringInter(rawLine)

	if strings.HasPrefix(lineClean, "if ") && strings.HasSuffix(lineClean, "{") {
		cond := strings.TrimSpace(lineClean[3 : len(lineClean)-1])
		return Token{Type: "IF", Value: parseCondition(cond)}
	}

	checkLine := lineClean
	if strings.HasPrefix(checkLine, "}") {
		checkLine = strings.TrimSpace(checkLine[1:])
	}

	if (strings.HasPrefix(checkLine, "elseif ") || strings.HasPrefix(checkLine, "else if ")) && strings.HasSuffix(checkLine, "{") {
		idx := strings.Index(checkLine, "if ")
		cond := strings.TrimSpace(checkLine[idx+3 : len(checkLine)-1])

		return Token{Type: "ELIF", Value: parseCondition(cond)}
	}

	if checkLine == "else {" {
		return Token{Type: "ELSE", Value: ""}
	}

	if strings.HasPrefix(lineClean, "while ") && strings.HasSuffix(lineClean, "{") {
		cond := strings.TrimSpace(lineClean[6 : len(lineClean)-1])
		return Token{Type: "WHILE", Value: parseCondition(cond)}
	}

	if strings.HasPrefix(lineClean, "fdo ") && strings.HasSuffix(lineClean, "{") {
		cond := strings.TrimSpace(lineClean[4 : len(lineClean)-1])
		return Token{Type: "UNTIL", Value: parseCondition(cond)}
	}

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

	if strings.HasPrefix(rawLine, "defer ") {
		innerCmd := strings.TrimSpace(rawLine[6:])
		innerToken := ParseLine(Token{Type: "SPECIAL", Value: innerCmd})
		return Token{Type: "DEFER", Value: innerToken.Value}
	}

	if strings.HasPrefix(lineClean, "function ") && strings.HasSuffix(lineClean, "{") {
		fnHeader := strings.TrimSpace(lineClean[9 : len(lineClean)-1])
		if !strings.HasSuffix(fnHeader, "()") {
			fnHeader = fnHeader + " ()"
		}
		return Token{Type: "FN_START", Value: fnHeader + " {"}
	}

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
