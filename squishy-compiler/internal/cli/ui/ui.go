package ui

import (
	_ "embed"
	"fmt"
	"strings"

	"github.com/Cod2rDude/squishy/squishy-compiler/internal/cli/color"
)

// Variables

//go:embed assets/banner.txt
var banner string

// Public Functions
func Log(append int, option string, message string) {
	switch option {
	case "warning":
		fmt.Println(strings.Repeat(" ", append) + color.Paint(color.Orange, "[WARNING] ") + color.Paint(color.Reset, message))
	case "error":
		fmt.Println(strings.Repeat(" ", append) + color.Paint(color.Red, "[ERROR] ") + color.Paint(color.Reset, message))
	default:
		fmt.Println(strings.Repeat(" ", append) + color.Paint(color.Blue, "[INFO] ") + color.Paint(color.Reset, message))
	}
}

func Startup() {
	fmt.Println(color.Paint(color.Green, banner))
	fmt.Println(color.Paint(
		color.Green, "| Developed by: ") + color.Paint(color.Blue, "Cod2rDude") +
		color.Paint(color.Green, "                                                                                                    |"))
	fmt.Println(color.Paint(color.Green, "+----------------------------------------------------------------------------------------------------------------------------+"))
	fmt.Println("")
}
