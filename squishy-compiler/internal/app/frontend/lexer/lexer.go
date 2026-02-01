/*
    ******************************************************************************
    * @file     : squishy/squishy-compiler/internal/app/frontend/lexer/lexer.go
    * @author   : Cod2rDude
    * @date     : January 20 2026
    * @lastEdit : January 28 2025 @ 15:36
    * @brief    : Squishy IDL Compiler Lexer.
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

package lexer

import (
    "fmt"
    "strconv"
    "strings"
    "text/scanner"

    "github.com/Cod2rDude/squishy/squishy-compiler/internal/app/language"
    "github.com/Cod2rDude/squishy/squishy-compiler/internal/app/types"
    "github.com/Cod2rDude/squishy/squishy-compiler/internal/cli/ui"
    "github.com/Cod2rDude/squishy/squishy-compiler/internal/config"
    "github.com/Cod2rDude/squishy/squishy-compiler/internal/errors"
)

// Variables
var tokenIsToString = map[int]string{
    0: "Unknown",
    1: "Keyword",
    2: "Type",
    3: "Struct Name",
    4: "Int",
    5: "Operator",
    6: "Invalid",
    7: "Field Name",
    8: "Export Name",
}

// Functions
func castTokenIsToString(token *types.Token) string {
    return tokenIsToString[token.Is]
}

// Public Structs

/*
    @object Lexer

    @privatevariables
    *   @privatevariable s : Scanner ;; Scanner.
    @publicvariables
    *   @publicvariable TokenList : []types.Token ;; List containing tokens.
    *   @publicvariable Cursor : int ;; Location of Cursor.
    *   @publicvariable StructReferences : []int ;; Location of struct references in @object:TokenList.
    *   @publicvariable FieldReferences : []int ;; Location of field references in @object:TokenList.
    *   @publicvariable ExportReferences : []int ;; Location of export references in @object:TokenList.
    @privatemethods
    *   @privatemethod analyzeAndCategorizeToken
    @publicmethods
    *   @publicmethod Scan
    *   @publicmethod GetAtCursor
    *   @publicmethod Next
    *   @publicmethod LookAtBack
    *   @publicmethod LookAtFront
    *   @publicmethod LookAhead
    *   @publicmethod ExpectAhead
    *   @publicmethod StepCursorForward
    *   @publicmethod StepCursorBackward
    *   @publicmethod JumpCursorAhead
    *   @publicmethod Expect
    *   @publicmethod Feed
    *   @publicmethod Length
    *   @publicmethod GetPiece
    *   @publicmethod Print
    @brief A custom lexer for Squishy IDL.
*/
type Lexer struct {
    s                scanner.Scanner
    TokenList        []types.Token
    Cursor           int
    StructReferences []int
    FieldReferences  []int
    ExportReferences []int
}

// Constructor
func New(input string) *Lexer {
    var s scanner.Scanner
    s.Init(strings.NewReader(input))

    s.Mode = scanner.ScanIdents | scanner.ScanInts | scanner.ScanComments
    TokenList := []types.Token{}

    return &Lexer{
        s:         s,
        TokenList: TokenList,
        Cursor:    0,
    }
}

// Private Methods
func (lexer *Lexer) analyzeAndCategorizeToken(tok rune, last *types.Token) (int, *errors.StackError) {
    is := types.UnknownToken
    position := lexer.s.Pos()
    text := lexer.s.TokenText()

    // Too nested.
    switch tok {
    case scanner.Ident:
        if language.Keywords[text] {
            is = types.KeywordToken
            switch text {
            case "struct":
                lexer.StructReferences = append(lexer.StructReferences, len(lexer.TokenList))
            case "field":
                lexer.FieldReferences = append(lexer.FieldReferences, len(lexer.TokenList))
            case "exports":
                lexer.ExportReferences = append(lexer.ExportReferences, len(lexer.TokenList))
            }
        } else if last != nil {
            switch last.Value {
            case "struct":
                is = types.StructNameToken
            case "field":
                is = types.FieldNameToken
            case "exports":
                is = types.ExportNameToken
            default:
                is = types.TypeToken
            }
        }
        // It's parsers job to make sure map and array definitions are correct.
        // Also parser must make sure a struct as a type exists.
        // But im actualy doing most of parser's job here.
    case scanner.Int:
        is = types.IntToken
    case scanner.Comment:
        return types.CommentToken, nil
    default:
        if !language.Operators[text] {
            return is, errors.New(errors.UnknownOperator, text, position.String())
        }

        is = types.OperatorToken
    }

    if last == nil && is != types.KeywordToken {
        return is, errors.New(errors.UnexpectedTokenAtStart, text)
    }

    return is, nil
}

// Public Methods
func (lexer *Lexer) Scan() *errors.StackError {
    var last *types.Token = nil
    for tok := lexer.s.Scan(); tok != scanner.EOF; tok = lexer.s.Scan() {
        is, err := lexer.analyzeAndCategorizeToken(tok, last)
        if err != nil {
            return err
        }

        if is == types.CommentToken {
            continue
        }

        position := lexer.s.Pos()
        text := lexer.s.TokenText()

        lexer.TokenList = append(lexer.TokenList, types.Token{
            Position:     len(lexer.TokenList),
            Is:           is,
            Value:        text,
            RealPosition: position.String(),
        })

        last = &lexer.TokenList[len(lexer.TokenList)-1]
    }
    return nil
}

func (lexer *Lexer) GetAtCursor() *types.Token {
    if lexer.Cursor >= len(lexer.TokenList) {
        return &types.Token{
            Is:           types.InvalidToken,
            Value:        "EOF",
            Position:     -1,
            RealPosition: "EOF",
        }
    }

    return &lexer.TokenList[lexer.Cursor]
}

func (lexer *Lexer) Next() *types.Token {
    if lexer.Cursor < len(lexer.TokenList) {
        lexer.Cursor++
    }

    out := lexer.GetAtCursor()

    return out
}

func (lexer *Lexer) LookAtBack() *types.Token {
    if lexer.Cursor > 0 {
        return &lexer.TokenList[lexer.Cursor-1]
    }

    return nil
}

func (lexer *Lexer) LookAtFront() *types.Token {
    if lexer.Cursor+1 < len(lexer.TokenList) {
        return &lexer.TokenList[lexer.Cursor+1]
    }

    return nil
}

func (lexer *Lexer) LookAhead(n int) *types.Token {
    target := lexer.Cursor + n

    if target >= 0 && target < len(lexer.TokenList) {
        return &lexer.TokenList[target]
    }

    return nil
}

func (lexer *Lexer) ExpectAhead(expected string) bool {
    for i := lexer.Cursor; i < len(lexer.TokenList); i++ {
        if lexer.TokenList[i].Value == expected {
            return true
        }
    }

    return false
}

func (lexer *Lexer) StepCursorForward(n int) {
    if n > 0 && lexer.Cursor+n <= len(lexer.TokenList) {
        lexer.Cursor += n
    }
}

func (lexer *Lexer) StepCursorBackward(n int) {
    if n > 0 && lexer.Cursor-n >= 0 {
        lexer.Cursor -= n
    }
}

func (lexer *Lexer) JumpCursorAhead(n int) {
    if n > lexer.Cursor && n <= len(lexer.TokenList) {
        lexer.Cursor = n
    }
}

func (lexer *Lexer) Expect(expected string) bool {
    token := lexer.GetAtCursor()
    return token.Is != types.InvalidToken && token.Value == expected
}

func (lexer *Lexer) Feed(input string) {
    lexer.s.Init(strings.NewReader(input))
    lexer.TokenList = []types.Token{}
    lexer.Cursor = 0
}

func (lexer *Lexer) Length() int {
    return len(lexer.TokenList)
}

func (lexer *Lexer) GetPiece(start int, length int) []types.Token {
    if start < 0 || start+length > len(lexer.TokenList) {
        return nil
    }

    end := start + length
    piece := make([]types.Token, length)
    copy(piece, lexer.TokenList[start:end])

    return piece
}

func (lexer *Lexer) Print() {
    ui.Log(config.FELLOWCRAFT, "info", "Printing tokens")
    for _, token := range lexer.TokenList {
        ui.Log(config.APPRENTICE, "info", fmt.Sprintf("Token #%d, Position: '%s', Is: '%s', Value: '%s'",
            token.Position, token.RealPosition, castTokenIsToString(&token), token.Value))
    }
    ui.Log(config.APPRENTICE, "info", "Count of struct references: "+strconv.Itoa(len(lexer.StructReferences)))
    ui.Log(config.APPRENTICE, "info", "Count of field references: "+strconv.Itoa(len(lexer.FieldReferences)))
    ui.Log(config.APPRENTICE, "info", "Count of export references: "+strconv.Itoa(len(lexer.ExportReferences)))
    if len(lexer.ExportReferences) > 1 || len(lexer.ExportReferences) == 0 {
        ui.Log(config.APPRENTICE, "warning", "Count of export references normally must be 1!")
    }
    ui.Log(config.FELLOWCRAFT, "info", "Finished printing tokens.")
}