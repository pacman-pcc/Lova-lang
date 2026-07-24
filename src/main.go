package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const Version = "0.3-beta"

func translateFile(ozPath string) string {
	if !strings.HasSuffix(ozPath, ".lova") {
		fmt.Printf("LOVA: File %s not found extension .lova\n", ozPath)
		return ""
	}

	shPath := ozPath[:len(ozPath)-5] + ".sh"

	fmt.Printf("Translate: %s :: %s..\n", ozPath, shPath)
	start_time := time.Now()

	code, err := os.ReadFile(ozPath)
	if err != nil {
		fmt.Printf("LOVA: Error build :: %v\n", err)
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
		fmt.Printf("LOVA: Error build :: %v\n", err)
		return ""
	}

	fmt.Println("LOVA: Ready ::")
	over_time := time.Since(start_time)
	fmt.Printf("LOVA: Time compile: %v\n", over_time)
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
				fmt.Printf("LOVA: %s not found.\n", target)
			}

		} else if fileExists(target) {
			translateFile(target)
		} else {
			fmt.Printf("LOVA: %s not found.\n", target)
		}


	} else {
		ozFiles, err := filepath.Glob("*.lova")
		if err != nil || len(ozFiles) == 0 {
			fmt.Println("LOVA: *.lova files not found in directory")
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
