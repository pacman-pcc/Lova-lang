package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const Version = "0.5-beta"

const GREEN = "\033[0;32m"
const CYAN = "\033[0;36m"
const RED = "\033[0;31m"
const PURPLE = "\033[0;35m"
const NC = "\033[0m"

func translateFile(ozPath string) string {
	if !strings.HasSuffix(ozPath, ".lova") {
		fmt.Printf("%sLOVA: File %s not found extension .lova%s\n", RED, ozPath, NC)
		return ""
	}

	shPath := ozPath[:len(ozPath)-5] + ".sh"

	fmt.Printf("%sLOVA:%s Translate: %s :: %s..\n", PURPLE, NC, ozPath, shPath)
	start_time := time.Now()

	code, err := os.ReadFile(ozPath)
	if err != nil {
		fmt.Printf("%sLOVA: Error build :: %v%s\n", RED, err, NC)
		return ""
	}

	rawLines := Lex(string(code))

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

		if target == "help" || target == "-h" {
			fmt.Println("-r/run => Running program in Lova-lang.")
			fmt.Println("-v/ver/version => Check version Lova.")
			fmt.Println("lova <file.lova> => Compile Lova not run.")
			fmt.Println("lova => Compile all file *.lova in dir.")
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
				fmt.Printf("%sLOVA: %s not found. or help menu%s\n", RED, target, NC)
			}

		} else if fileExists(target) {
			translateFile(target)
		} else {
			fmt.Printf("%sLOVA: %s not found. or help menu%s\n", RED, target, NC)
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
