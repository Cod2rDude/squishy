/*
    ******************************************************************************
    * @file     : squishy/squishy-compiler/internal/app/frontend/frontend.go
    * @author   : Cod2rDude
    * @date     : January 20 2026
    * @lastEdit : January 28 2025 @ 15:38
    * @brief    : Squishy IDL Compiler Frontend.
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

package frontend

import (
    "github.com/Cod2rDude/squishy/squishy-compiler/internal/app/frontend/lexer"
    "github.com/Cod2rDude/squishy/squishy-compiler/internal/app/frontend/parser"
    "github.com/Cod2rDude/squishy/squishy-compiler/internal/app/types"
    "github.com/Cod2rDude/squishy/squishy-compiler/internal/cli/ui"
    "github.com/Cod2rDude/squishy/squishy-compiler/internal/config"
    "github.com/Cod2rDude/squishy/squishy-compiler/internal/errors"
    "github.com/Cod2rDude/squishy/squishy-compiler/internal/file"
)

// Public Structs

/*
    @object Frontend

    @privatevariables
    *   @privatevariable path : string ;; Source file dir/path.
    *   @privatevariable myLexer : lexer.Lexer ;; Lexer.
    *   @privatevariable myParser : parser.Parser ;; Parser.
    @publicvariables
    *   @publicvariable Result : types.Scheme ;; Result of lexing and parsing
    @publicmethods
    *   @publicmethod Work
    *   @publicmethod Debug
    @brief Frontend of Squishy IDL Compiler.
*/
type Frontend struct {
    path     string
    myLexer  *lexer.Lexer
    myParser *parser.Parser
    Result   *types.Scheme
}

// Constructor
func New(path string) (*Frontend, *errors.StackError) {
    _, err := file.IsAValidFile(path)
    if err != nil {
        return nil, err
    }

    _, err2 := file.HasAnyValidExtension(path, config.DefaultExpectedFileExtensions)
    if err2 != nil {
        return nil, err2
    }

    myLexer := lexer.New("")
    myParser := parser.New(myLexer)

    return &Frontend{
        path:     path,
        myLexer:  myLexer,
        myParser: myParser,
        Result: &types.Scheme{
            Structs: make(map[string]*types.Struct),
            Exports: "",
        },
    }, nil
}

func NewFromString(input string) *Frontend {
    myLexer := lexer.New("")
    myParser := parser.New(myLexer)

    return &Frontend{
        path:     "",
        myLexer:  myLexer,
        myParser: myParser,
        Result: &types.Scheme{
            Structs: make(map[string]*types.Struct),
            Exports: "",
        },
    }
}

// Public Methods
func (frontend *Frontend) WorkFromString(input string) *errors.StackError {
    frontend.myLexer = lexer.New(input)
    err1 := frontend.myLexer.Scan()
    if err1 != nil {
        return err1
    }
    frontend.myParser = parser.New(frontend.myLexer)
    err2 := frontend.myParser.Parse()
    if err2 != nil {
        return err2
    }

    frontend.Result = &frontend.myParser.Result

    return nil
}

func (frontend *Frontend) Work() *errors.StackError {
    fileContents, err := file.FileToString(frontend.path)
    if err != nil {
        return err
    }

    frontend.myLexer = lexer.New(fileContents)
    err2 := frontend.myLexer.Scan()
    if err2 != nil {
        return err2
    }
    frontend.myParser = parser.New(frontend.myLexer)
    err3 := frontend.myParser.Parse()
    if err3 != nil {
        return err3
    }

    frontend.Result = &frontend.myParser.Result

    return nil
}

func (frontend *Frontend) Debug() {
    ui.Log(config.MASTER, "info", "FRONTEND DEBUG START")
    frontend.myLexer.Print()
    frontend.myParser.Print()
    ui.Log(config.MASTER, "info", "FRONTEND DEBUG END")
}
