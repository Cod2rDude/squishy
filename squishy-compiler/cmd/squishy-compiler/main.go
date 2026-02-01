package main

import (
	"bufio"
	"flag"
	"os"
	"strings"

	"fmt"

	app "github.com/Cod2rDude/squishy/squishy-compiler/internal/app"
	"github.com/Cod2rDude/squishy/squishy-compiler/internal/cli/ui"
	"github.com/Cod2rDude/squishy/squishy-compiler/internal/config"
)

// Private Functions
func getArgs() (string, bool, string) {
    reader := bufio.NewReader(os.Stdin)

    ui.Log(0, "info", "Enter input file path,")
    fmt.Print(">> ")
    text, _ := reader.ReadString('\n')
    inputFile := strings.TrimSpace(text)

    ui.Log(0, "info", "Enter output path,")
    fmt.Print(">> ")
    text, _ = reader.ReadString('\n')
    outputPath := strings.TrimSpace(text)

    if outputPath == "" {
        outputPath = "./"
    }

    debugEnabled := false

    ui.Log(0, "info", "Enable debug mode? (y/N)")
    fmt.Print(">> ")
    text, _ = reader.ReadString('\n')
    text = strings.TrimSpace(strings.ToLower(text))
    if text == "y" || text == "yes" {
        debugEnabled = true
    }

    fmt.Println()

    return inputFile, debugEnabled, outputPath
}

// Main
func main() {
    ui.Startup()
    ui.Log(0, "info", fmt.Sprintf("Welcome to 'squishy-compiler' Version %s!", config.Version))
    fmt.Println()

    var outputPath string
    var debug bool
    var inputFile string

    flag.StringVar(&outputPath, "o", "./", "Output path.")
    flag.BoolVar(&debug, "debug", false, "Enable debug mode.")

    if len(os.Args) < 2 {
        inputFile, debug, outputPath = getArgs()
    } else {
        flag.Parse()
        args := flag.Args()
        if len(args) < 1 {
            ui.Log(0, "error", "No input file specified?")
            os.Exit(-1)
        }
        inputFile = args[0]
    }

    ui.Log(0, "warning", fmt.Sprintf("Input file: %s, Output path: %s, Debug Mode: %t", inputFile, outputPath, debug))
    fmt.Println("")

    if err := app.App(inputFile, outputPath, debug); err != nil {
        err.Throw('d', true)
    }
}