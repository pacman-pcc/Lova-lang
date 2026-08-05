package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const Version = "1.2-stable"

const GREEN = "\033[0;32m"
const CYAN = "\033[0;36m"
const RED = "\033[0;31m"
const PURPLE = "\033[0;35m"
const NC = "\033[0m"

var embeddedLibs = map[string]string{
	"colors": `
proc RED = "\033[0;31m"
proc GREEN = "\033[0;32m"
proc BLUE = "\033[0;34m"
proc CYAN = "\033[0;36m"
proc PURPLE = "\033[0;35m"
proc YELLOW = "\033[0;33m"
proc NC = "\033[0m"
`,
	"os": `
proc OS_NAME = "$(uname -s)"
proc IS_LINUX = "$([[ $(uname -s) == 'Linux' ]] && echo true || echo false)"
proc IS_MAC = "$([[ $(uname -s) == 'Darwin' ]] && echo true || echo false)"
`,
}

func translateFile(ozPath string) string {
	if !strings.HasSuffix(ozPath, ".lova") {
		fmt.Printf("%sLOVA: File %s not found extension .lova%s\n", RED, ozPath, NC)
		return ""
	}

	shPath := ozPath[:len(ozPath)-5] + ".sh"

	fmt.Printf("%sLOVA:%s Translate: %s :: %s..\n", PURPLE, NC, ozPath, shPath)
	start_time := time.Now()

	var resolveIncludes func(path string, visited map[string]bool) (string, error)
	resolveIncludes = func(path string, visited map[string]bool) (string, error) {
		absPath, err := filepath.Abs(path)
		if err != nil {
			return "", err
		}
		if visited[absPath] {
			return "", nil // защита от циклического инклуда
		}
		visited[absPath] = true

		contentBytes, err := os.ReadFile(absPath)
		if err != nil {
			return "", err
		}

		baseDir := filepath.Dir(absPath)
		lines := strings.Split(string(contentBytes), "\n")
		var result []string

		for _, line := range lines {
			trimmed := strings.TrimSpace(line)
			if strings.HasPrefix(trimmed, "import ") {
				impName := strings.Trim(strings.TrimSpace(trimmed[6:]), "\"'")

				if libContent, ok := embeddedLibs[impName]; ok {
					result = append(result, libContent)
					continue
				}

				if !strings.HasSuffix(impName, ".lova") {
					impName += ".lova"
				}
				impPath := filepath.Join(baseDir, impName)

				includedContent, incErr := resolveIncludes(impPath, visited)
				if incErr != nil {
					return "", incErr
				}
				result = append(result, includedContent)
			} else {
				result = append(result, line)
			}
		}

		return strings.Join(result, "\n"), nil
	}

	visited := make(map[string]bool)
	fullCode, err := resolveIncludes(ozPath, visited)
	if err != nil {
		fmt.Printf("%sLOVA: Error import :: %v%s\n", RED, err, NC)
		return ""
	}

	rawLines := Lex(fullCode)

	parsedItems := make([]Token, 0, len(rawLines))
	for _, item := range rawLines {
		parsedItems = append(parsedItems, ParseLine(item))
	}

	bashCode := GenBash(parsedItems)

	err = os.WriteFile(shPath, []byte(bashCode), 0775)
	if err != nil {
		fmt.Printf("%sLOVA: Error build :: %v%s\n", RED, err, NC)
		return ""
	}

	over_time := time.Since(start_time)
	fmt.Printf("%sLOVA: Time compile: %v%s\n", CYAN, over_time, NC)
	return shPath
}

func main() {
	if len(os.Args) > 1 {
		target := os.Args[1]

		if target == "ver" || target == "version" || target == "-v" || target == "--version" {
			fmt.Printf("LOVA Compiler v%s\n", Version)
			return
		}

		if target == "run" || target == "-r" {
			target = "main.lova"
			if len(os.Args) > 2 {
				target = os.Args[2]
			}

			if fileExists(target) {
				shFile := translateFile(target)
				if shFile != "" {
					cmd := exec.Command("bash", shFile)
					cmd.Stdout = os.Stdout
					cmd.Stderr = os.Stderr
					cmd.Stdin = os.Stdin
					_ = cmd.Run()
				}
			} else {
				fmt.Printf("%sLOVA: %s not found.%s\n", RED, target, NC)
			}

		} else if fileExists(target) {
			translateFile(target)
		} else {
			fmt.Printf("%sLOVA: %s not found.%s\n", RED, target, NC)
		}

	} else {
		ozFiles, err := filepath.Glob("*.lova")
		if err != nil || len(ozFiles) == 0 {
			fmt.Printf("%sLOVA: *.lova files not found in directory%s\n", RED, NC)
			return
		}

		for _, ozFile := range ozFiles {
			translateFile(ozFile)
		}
	}
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
