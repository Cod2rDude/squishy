/*
    ******************************************************************************
    * @file     : squishy/squishy-compiler/internal/app/app.go
    * @author   : Cod2rDude
    * @date     : January 20 2026
    * @lastEdit : January 29 2025 @ 20:43
    * @brief    : Squishy IDL Compiler Application Entry Point.
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

package app

import (
	"fmt"
	"path/filepath"
	"time"

	be "github.com/Cod2rDude/squishy/squishy-compiler/internal/app/backend"
	fe "github.com/Cod2rDude/squishy/squishy-compiler/internal/app/frontend"
	me "github.com/Cod2rDude/squishy/squishy-compiler/internal/app/middleend"
	"github.com/Cod2rDude/squishy/squishy-compiler/internal/cli/ui"
	"github.com/Cod2rDude/squishy/squishy-compiler/internal/config"
	"github.com/Cod2rDude/squishy/squishy-compiler/internal/errors"
	"github.com/Cod2rDude/squishy/squishy-compiler/internal/file"
)

// Public Functions
func App(inputFile string, outputDirectory string, debug bool) *errors.StackError {
    if _, err := file.IsAValidFile(inputFile); err != nil {
        return err
    }

    if _, err := file.IsADirectory(outputDirectory); err != nil {
        return err
    }

    ui.Log(config.GRANDMASTER, "info", "Compiling.....")

    startTime := time.Now()

    frontend, err := fe.New(inputFile)
    if err != nil {
        return err
    }

    err2 := frontend.Work()
    if err2 != nil {
        return err2
    }

    middleend := me.New(frontend.Result)
    middleend.Work()

    typeString, exportString := middleend.GetResults()

    backend, err4 := be.New(outputDirectory, frontend.Result, exportString, typeString, middleend.GetSortedStructs())
    if err4 != nil {
        return err4
    }
    err5 := backend.Work()
    if err5 != nil {
        return err5
    }

    totalTime := time.Since(startTime) * time.Millisecond / 10000
    
    fmt.Println("")
    ui.Log(config.GRANDMASTER, "info", "Done compiling! Took " + fmt.Sprintf("%s", totalTime))
    ui.Log(config.GRANDMASTER, "info", "Output is at '"+filepath.Join(outputDirectory, frontend.Result.Exports + ".luau")+"'")

    fmt.Println("")

    if debug {
        ui.Log(config.GRANDMASTER, "info", "STARTING DEBUG MODE")
        frontend.Debug()
        middleend.Debug()
        backend.Debug()
        ui.Log(config.GRANDMASTER, "info", "END OF DEBUG MODE")
    }

    return nil
}