/*
    ******************************************************************************
    * @file     : squishy/squishy-compiler/internal/app/backend/backend.go
    * @author   : Cod2rDude
    * @date     : January 26 2026
    * @lastEdit : January 27 2025 @ 09:24
    * @brief    : Squishy IDL Compiler Backend.
    * @version  : 1.0.0
    ******************************************************************************
    * @attention
    *
    * Copyright © 2026 Axon Corporation.
    * All rights reserved.
    *
    * This software is licensed under terms that can be found in the LICENSE file
    * in the root directory of this software component.
    * If no LICENSE file comes with this software, it is provided AS-IS.
    *
    ******************************************************************************
*/

package backend

import (
	"bufio"
	_ "embed"
	"strings"
	"time"

	"github.com/Cod2rDude/squishy/squishy-compiler/internal/app/types"
	"github.com/Cod2rDude/squishy/squishy-compiler/internal/cli/ui"
	"github.com/Cod2rDude/squishy/squishy-compiler/internal/config"
	"github.com/Cod2rDude/squishy/squishy-compiler/internal/errors"
	"github.com/Cod2rDude/squishy/squishy-compiler/internal/file"
)

// Variables

//go:embed assets/template.luau
var template string

// Public Structs

/*
    @object Backend

    @privatevariables
    *   @privatevariable scheme : *types.Scheme ;; Pointer to scheme created by frontend.
    @privatemethods
    *   @privatemethod
    @publicmethods
    *   @publicmethod Work
    *   @publicmethod Debug
    @brief Backend of Squishy IDL Compiler.
*/
type Backend struct {
    outputPath    string
    exportString  string
    typeString    string
    scheme        *types.Scheme
    sortedStructs []string
}

// Constructor
func New(outputPath string, scheme *types.Scheme,
    exportString string, typeString string,
    sortedStructs []string,
) (*Backend, *errors.StackError) {
    if c, err := file.IsADirectory(outputPath); !c || err != nil {
        return nil, err
    }

    return &Backend{
        outputPath:    outputPath,
        scheme:        scheme,
        exportString:  exportString,
        typeString:    typeString,
        sortedStructs: sortedStructs,
    }, nil
}

func NewWithoutPath(scheme *types.Scheme, 
    exportString string, typeString string, sortedStructs []string) *Backend {
        return &Backend{
        outputPath:    "",
        scheme:        scheme,
        exportString:  exportString,
        typeString:    typeString,
        sortedStructs: sortedStructs,
    }
}

// Public Methods
func (backend *Backend) GetString() string {
    var lines []string
    scanner := bufio.NewScanner(strings.NewReader(template))

    for scanner.Scan() {
        lines = append(lines, scanner.Text())
    }

    now := time.Now()

    writeFunctions := StructListToWriteFunctions(backend.sortedStructs, backend.scheme.Structs)
    readFunctions := StructListToReadFunctions(backend.sortedStructs, backend.scheme.Structs)
    exportFunctionWriteBody := StructToWriteString(backend.scheme.Structs[backend.scheme.Exports])
    exportFunctionReadBody, exportFunctionReadReturn := StructToReadString(backend.scheme.Structs[backend.scheme.Exports])

    lines[7] = format("    * @file     : %s%s%s%s", backend.outputPath, "/", backend.scheme.Exports, ".luau")
    lines[8] = "    * @author   : squishy-compiler"
    lines[9] = format("    * @date     : %s", now.Format("January 2 2006"))
    lines[10] = format("    * @lastEdit : %s @ %s", now.Format("January 2 2006"), now.Format("15:04"))
    lines[11] = format("    * @brief    : Squishy IDL Compiler generated code for %s.", backend.scheme.Exports)
    lines[40] = format("function scheme.write(input : %s) : buffer?", backend.scheme.Exports)
    lines[48] = format("function scheme.read(buff : buffer) : %s?", backend.scheme.Exports)
    lines[51] = "    return " + strings.Join(strings.Split(exportFunctionReadReturn, ";"), ";\n            ")

    out := []string{}

    out = append(out, lines[0:28]...)
    out = append(out, backend.typeString)
    out = append(out, lines[29:33]...)
    out = append(out, writeFunctions...)
    out = append(out, readFunctions...)
    out = append(out, lines[34:38]...)
    out = append(out, backend.exportString)
    out = append(out, lines[39:42]...)
    out = append(out, " ")
    out = append(out, "    "+strings.Join(exportFunctionWriteBody, "\n    "))
    out = append(out, " ")
    out = append(out, lines[43:50]...)
    out = append(out, " ")
    out = append(out, "    " + strings.Join(exportFunctionReadBody, "\n    "))
    out = append(out, " ")
    out = append(out, lines[51:]...)

    finalOutput := strings.Join(out, "\n")

    return finalOutput
}

func (backend *Backend) Work() *errors.StackError {
    finalOutput := backend.GetString()

    if err := file.CreateAndWriteFile(backend.outputPath, backend.scheme.Exports+".luau", finalOutput); err != nil {
        return err
    }

    return nil
}

func (backend *Backend) Debug() {
    writeFunctions := StructListToWriteFunctions(backend.sortedStructs, backend.scheme.Structs)
    exportFunctionWriteBody := StructToWriteString(backend.scheme.Structs[backend.scheme.Exports])
    exportFunctionReadBody, exportFunctionReadReturn := StructToReadString(backend.scheme.Structs[backend.scheme.Exports])

    writeFunctions = writeFunctions[:len(writeFunctions)-1]

    ui.Log(config.MASTER, "info", "BACKEND DEBUG START")
    ui.Log(config.FELLOWCRAFT, "info", "Write Functions:")
    ui.Log(config.APPRENTICE, "info", strings.Join(writeFunctions, "\n"+strings.Repeat(" ", 13)))
    ui.Log(config.FELLOWCRAFT, "info", "Read Functions:")
    ui.Log(config.FELLOWCRAFT, "info", "Write Body:")
    ui.Log(config.APPRENTICE, "info", strings.Join(exportFunctionWriteBody, "\n"+strings.Repeat(" ", 13)))
    ui.Log(config.FELLOWCRAFT, "info", "Read Body:")
    ui.Log(config.APPRENTICE, "info", strings.Join(exportFunctionReadBody, "\n"+strings.Repeat(" ", 13)))
    ui.Log(config.FELLOWCRAFT, "info", "Read Return:")
    ui.Log(config.APPRENTICE, "info", exportFunctionReadReturn)
    ui.Log(config.MASTER, "info", "BACKEND DEBUG END")
}
