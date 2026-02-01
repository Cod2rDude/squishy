package errors

import (
	stdErrors "errors"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/Cod2rDude/squishy/squishy-compiler/internal/cli/color"
)

// Variables
var stackSkipCount int = 3

// Functions
func captureStack() []runtime.Frame {
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(stackSkipCount, pcs[:])

	frames := runtime.CallersFrames(pcs[:n])
	var cleanFrames []runtime.Frame

	for {
		frame, more := frames.Next()

		if !strings.HasPrefix(frame.Function, "runtime.") {
			cleanFrames = append(cleanFrames, frame)
		}

		if !more {
			break
		}
	}
	return cleanFrames
}

// Public Structs
type StackError struct {
	Errs  map[int]error
	Code  int
	Stack []runtime.Frame
}

// Constructor
func New(code int, args ...any) *StackError {
	if code > len(errorCodeToString) || code < 0 {
		code = UnknownError
	}

	return &StackError{
		Errs:  map[int]error{0: stdErrors.New(fmt.Sprintf(errorCodeToString[code], args...))},
		Code:  code,
		Stack: captureStack(),
	}
}

// Private Methods
func (stackError *StackError) printErrors() {
	doesHaveMoreThan1Error := len(stackError.Errs) > 1

	for i, err := range stackError.Errs {
		index := ""

		if doesHaveMoreThan1Error {
			index = color.Paint(color.Bold, strconv.Itoa(i)) + ": "
		}

		outputString := index + color.Paint(color.Red, "[ERROR] ") + err.Error()
		fmt.Println(outputString)
	}
}

// Public Methods
func (stackError *StackError) Add(code int, args ...any) {
	if code > len(errorCodeToString) || code < 0 {
		code = UnknownError
	}

	stackError.Errs[len(stackError.Errs)] = stdErrors.New(fmt.Sprintf(errorCodeToString[code], args...))
}

func (stackError *StackError) Format(debug bool) string {
    var sb strings.Builder
    doesHaveMoreThan1Error := len(stackError.Errs) > 1

    for i := 0; i < len(stackError.Errs); i++ {
        err := stackError.Errs[i]
        index := ""

        if doesHaveMoreThan1Error {
            index = color.Paint(color.Bold, strconv.Itoa(i)) + ": "
        }

        sb.WriteString(index + color.Paint(color.Red, "[ERROR] ") + err.Error() + "\n")
    }

    if debug {
        sb.WriteString("\n" + color.Paint(color.Blue, "[STACK]") + "\n")
        for _, frame := range stackError.Stack {
            line := fmt.Sprintf("\t%s ~ %s:%d\n", frame.Function, frame.File, frame.Line)
            sb.WriteString(line)
        }
    }

    return sb.String()
}

func (stackError *StackError) Error() string {
    return stackError.Format(true)
}

func (stackError *StackError) Throw(verb rune, exit bool) {
    switch verb {
    case 's':
        fmt.Print(stackError.Format(false))
    case 'd':
        fmt.Print(stackError.Format(true))
    }

    if exit {
        if runtime.GOARCH == "wasm" && runtime.GOOS == "js" {
            panic("WASM_EXIT") 
        }
        os.Exit(stackError.Code)
    }
}