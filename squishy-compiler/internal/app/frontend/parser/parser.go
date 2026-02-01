/*
    ******************************************************************************
    * @file     : squishy/squishy-compiler/internal/app/frontend/parser/parser.go
    * @author   : Cod2rDude
    * @date     : January 20 2026
    * @lastEdit : January 28 2025 @ 15:38
    * @brief    : Squishy IDL Compiler Parser.
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

package parser

import (
    "fmt"
    "strconv"
    "strings"

    "github.com/Cod2rDude/squishy/squishy-compiler/internal/app/frontend/lexer"
    "github.com/Cod2rDude/squishy/squishy-compiler/internal/app/language"
    "github.com/Cod2rDude/squishy/squishy-compiler/internal/app/types"
    "github.com/Cod2rDude/squishy/squishy-compiler/internal/cli/ui"
    "github.com/Cod2rDude/squishy/squishy-compiler/internal/config"
    "github.com/Cod2rDude/squishy/squishy-compiler/internal/errors"
    "github.com/Cod2rDude/squishy/squishy-compiler/internal/util"
)

// Functions
func getConcatenatedNames(fields []*types.Field, indexes []int) string {
    parts := make([]string, 0, len(indexes))

    for _, idx := range indexes {
        if idx >= 0 && idx < len(fields) {
            f := fields[idx]
            msg := fmt.Sprintf("%s:%s", f.Name, f.Type.Name)
            parts = append(parts, msg)
        }
    }

    return strings.Join(parts, ", ")
}

// Public Structs

/*
    @object Parser

    @privatevariables
    *   @privatevariable myLexer : *lexer.Lexer ;; Pointer to lexer.
    @publicvariables
    *   @publicvariable Result : types.Scheme ;; Result of parsing
    @privatemethods
    *   @privatemethod printStructFields
    *   @privatemethod getFieldTypeDescription
    *   @privatemethod printSingleStruct
    *   @privatemethod isTokenAValidType
    *   @privatemethod parseMap
    *   @privatemethod parseArray
    *   @privatemethod parseType
    *   @privatemethod parseField
    *   @privatemethod parseFields
    *   @privatemethod parseStructs
    *   @privatemethod parseExports
    *   @privatemethod checkPath
    *   @privatemethod detectCycles
    *   @privatemethod semanticAnalyze
    @publicmethods
    *   @publicmethod Parse
    *   @publicmethod Print
    @brief A custom lexer for Squishy IDL.
*/
type Parser struct {
    myLexer *lexer.Lexer
    Result  types.Scheme
}

// Private Methods
func (parser *Parser) printStructFields(fields []*types.Field) {
    ui.Log(config.UPPERCLASS, "info", "Fields:")

    for index, field := range fields {
        ui.Log(config.MIDCLASS, "info", fmt.Sprintf("%d'th Field", index+1))

        desc := parser.getFieldTypeDescription(field.Type)
        ui.Log(config.BOTTOMCLASS, "info", desc)
    }
}

func (parser *Parser) getFieldTypeDescription(t *types.Type) string {
    if t.IsArray {
        return fmt.Sprintf("Type: %s, Array, Dynamic: %t", t.Name, t.ArraySize != -1)
    }
    if t.IsMap {
        return fmt.Sprintf("Type: %s, Map", t.Name)
    }
    if t.IsReferenceToAnotherStruct {
        return fmt.Sprintf("Type: %s, Reference To Another Struct", t.Name)
    }
    return fmt.Sprintf("Type: %s", t.Name)
}

func (parser *Parser) printSingleStruct(s *types.Struct) {
    ui.Log(config.ROYAL, "info", "Struct "+s.Name+" At '"+s.Reference+"'")
    ui.Log(config.UPPERCLASS, "info", fmt.Sprintf("Field count: %d", len(s.Fields)))
    ui.Log(config.UPPERCLASS, "info", fmt.Sprintf("Ever Referenced: '%t'", s.EverReferenced))

    parser.printStructFields(s.Fields)

    if len(s.ReferencedBy) > 0 {
        ui.Log(config.UPPERCLASS, "info", "Structs That Referenced This Struct: ")
        for name, count := range s.ReferencedBy {
            ui.Log(config.BOTTOMCLASS, "info", fmt.Sprintf("Referenced by '%s', '%d' times.", name, count))
        }
    }

    if len(s.OtherStructReferences) > 0 {
        ui.Log(config.UPPERCLASS, "info", "Other Struct References:")
        for name, locations := range s.OtherStructReferences {
            ui.Log(config.MIDCLASS, "info", fmt.Sprintf("Reference '%s', Times: %d", name, len(locations)))
        }
    }
}

func (parser *Parser) isTokenAValidType(token *types.Token) *errors.StackError {
    if _, err := util.IsAValidName(token.Value); err != nil {
        return err
    }
    if token.Is != types.TypeToken {
        return errors.New(errors.ExpectedAValidType, token.RealPosition, token.Value)
    }

    return nil
}

func (parser *Parser) parseMap(_type *types.Type) *errors.StackError {
    _type.IsMap = true
    token1 := parser.myLexer.GetAtCursor()
    nextToken := parser.myLexer.Next()

    if nextToken.Value == "}" {
        return errors.New(errors.NoTypeSpecifiedForMap, token1.RealPosition)
    }

    if parser.myLexer.Next().Value != "}" {
        return errors.New(errors.CurlyBraceNotClosed, token1.RealPosition)
    }

    nextNextToken := parser.myLexer.Next() // Oh god

    if nextNextToken.Value != "map" && nextNextToken.Value != "smap" {
        return errors.New(errors.ExpectedMapDefinition, token1.RealPosition)
    }

    if err := parser.isTokenAValidType(nextToken); err != nil {
        return err
    }

    _type.IsShortMap = nextNextToken.Value == "smap"
    _type.Name = nextToken.Value

    return nil
}

func (parser *Parser) parseArray(_type *types.Type) *errors.StackError {
    _type.IsArray = true
    token1 := parser.myLexer.GetAtCursor()
    nextToken := parser.myLexer.Next()

    switch nextToken.Is {
    case types.IntToken:
        arraySize, err := strconv.Atoi(nextToken.Value)

        if err != nil {
            return errors.New(errors.UnknownError, err.Error())
        }

        _type.ArraySize = arraySize

        nextToken := parser.myLexer.Next()
        if nextToken.Value != "]" {
            return errors.New(errors.BracketNotClosed, token1.RealPosition, nextToken.Value, nextToken.RealPosition)
        }

        nextToken = parser.myLexer.Next()
        if err := parser.isTokenAValidType(nextToken); err != nil {
            return err
        }
        _type.Name = nextToken.Value
    default:
        if nextToken.Value != "]" {
            return errors.New(errors.BracketNotClosed, token1.RealPosition, nextToken.Value, nextToken.RealPosition)
        }
        nextToken = parser.myLexer.Next()
        if err := parser.isTokenAValidType(nextToken); err != nil {
            return err
        }
        _type.Name = nextToken.Value
    }

    return nil
}

func (parser *Parser) parseType() (types.Type, *errors.StackError) {
    _type := types.Type{
        Name:                       "",
        IsArray:                    false,
        ArraySize:                  -1,
        IsMap:                      false,
        IsShortMap:                 false,
        IsReferenceToAnotherStruct: false,
    }

    token1 := parser.myLexer.GetAtCursor()

    switch token1.Value {
    case "{": // Map
        if err := parser.parseMap(&_type); err != nil {
            return _type, err
        }
    case "[": // Array
        if err := parser.parseArray(&_type); err != nil {
            return _type, err
        }
    default: // Normal Type
        if err := parser.isTokenAValidType(token1); err != nil {
            return _type, err
        }
        _type.Name = token1.Value
    }

    if language.DefaultTypes[_type.Name] != true {
        _type.IsReferenceToAnotherStruct = true
    }

    return _type, nil
}

func (parser *Parser) parseField(at int) (int, types.Field, *errors.StackError) {
    field := types.Field{
        Name: "",
        Type: &types.Type{
            Name:                       "",
            IsArray:                    false,
            ArraySize:                  -1,
            IsMap:                      false,
            IsReferenceToAnotherStruct: false,
        },
    }

    parser.myLexer.JumpCursorAhead(at)
    token := parser.myLexer.GetAtCursor()

    name := parser.myLexer.LookAtFront()
    if _, err := util.IsAValidName(name.Value); err != nil {
        return 0, types.Field{}, err
    }
    if name.Is != types.FieldNameToken {
        return 0, types.Field{}, errors.New(errors.ExpectedNameForField, token.RealPosition)
    }

    parser.myLexer.StepCursorForward(2)
    _type, err := parser.parseType()
    if err != nil {
        return parser.myLexer.Cursor, field, err
    }

    field.Type = &_type
    field.Name = name.Value

    return parser.myLexer.Cursor, field, nil
}

func (parser *Parser) parseFields(_struct *types.Struct, start int) *errors.StackError {
    fieldNames := map[string]bool{}

    for _, tokenIndex := range parser.myLexer.FieldReferences {
        if tokenIndex < start {
            continue
        }

        for k := parser.myLexer.Cursor + 1; k < tokenIndex; k++ {
            token := parser.myLexer.TokenList[k]
            if token.Is != types.CommentToken {
                return errors.New(errors.UnexpectedTokenAfterField, token.Value, token.RealPosition)
            }
        }

        endedAt, field, err := parser.parseField(tokenIndex)
        if err != nil {
            return err
        }

        if fieldNames[field.Name] {
            return errors.New(errors.AnotherFieldWithSameNameExists, _struct.Name, field.Name)
        }

        fieldNames[field.Name] = true

        if field.Type.IsReferenceToAnotherStruct {
            _struct.OtherStructReferences[field.Type.Name] = append(_struct.OtherStructReferences[field.Type.Name], len(_struct.Fields))
        }

        _struct.Fields = append(_struct.Fields, &field)

        parser.myLexer.JumpCursorAhead(endedAt)
        if parser.myLexer.LookAtFront().Value == "}" {
            return nil
        }

        if parser.myLexer.LookAtFront().Value != "field" {
            return errors.New(errors.ExpectedFieldAfterAnotherField, field.Name, parser.myLexer.LookAtFront().RealPosition, _struct.Name)
        }
    }

    return nil
}

func (parser *Parser) parseStructs() *errors.StackError {
    for _, tokenIndex := range parser.myLexer.StructReferences {
        parser.myLexer.JumpCursorAhead(tokenIndex)
        token := parser.myLexer.GetAtCursor()

        name := parser.myLexer.LookAtFront()
        if _, err := util.IsAValidName(name.Value); err != nil {
            return err
        }
        if name.Is != types.StructNameToken {
            return errors.New(errors.ExpectedNameForStruct, token.RealPosition)
        }

        if language.DefaultTypes[name.Value] == true {
            return errors.New(errors.InvalidStructNaming, token.RealPosition, name.Value)
        }

        parser.myLexer.StepCursorForward(2)

        if tok := parser.myLexer.GetAtCursor(); tok.Value != "{" {
            return errors.New(errors.StructShouldStartWithCurlyBrace, token.RealPosition, tok.Value)
        }

        _struct := types.Struct{
            Reference:             token.RealPosition,
            Name:                  name.Value,
            Fields:                []*types.Field{},
            OtherStructReferences: make(map[string][]int),
            EverReferenced:        false,
            ReferencedBy:          make(map[string]int),
        }

        err := parser.parseFields(&_struct, tokenIndex)
        if err != nil {
            return err
        }

        val, found := parser.Result.Structs[_struct.Name]

        if found {
            return errors.New(errors.AnotherStructWithSameNameExists, _struct.Name, val.Reference, _struct.Reference)
        }

        if len(_struct.Fields) == 0 {
            return errors.New(errors.AStructMustHaveAtleast1Field, _struct.Name)
        }

        parser.Result.Structs[_struct.Name] = &_struct

        parser.myLexer.StepCursorForward(1)
        nextToken := parser.myLexer.LookAtFront()

        switch nextToken.Value {
        case "struct":
        case "exports":
        default:
            return errors.New(errors.UnexpectedTokenAfterStruct, nextToken.Value, nextToken.RealPosition)
        }
    }

    return nil
}

func (parser *Parser) parseExports() *errors.StackError {
    exportToken := parser.myLexer.TokenList[parser.myLexer.ExportReferences[0]]
    parser.myLexer.JumpCursorAhead(parser.myLexer.ExportReferences[0] + 1)
    exportNameToken := parser.myLexer.GetAtCursor()

    if _, err := util.IsAValidName(exportNameToken.Value); err != nil {
        return err
    }
    if exportNameToken.Is != types.ExportNameToken {
        return errors.New(errors.ExpectedNameForExport, exportToken.RealPosition)
    }

    _, found := parser.Result.Structs[exportNameToken.Value]

    if !found {
        return errors.New(errors.DidntFoundAStructToExport, exportNameToken.Value)
    }

    parser.Result.Exports = exportNameToken.Value

    return nil
}

func (parser *Parser) checkPath(currentName string, path []string, visited map[string]bool) *errors.StackError {
    if visited[currentName] {
        path = append(path, currentName)

        if len(visited) == 1 {
            return errors.New(errors.AStructCantReferenceItself, currentName, getConcatenatedNames(
                parser.Result.Structs[currentName].Fields, parser.Result.Structs[currentName].OtherStructReferences[currentName],
            ))
        }

        pathStr := strings.Join(path, " -> ")

        return errors.New(errors.CyclicReference, pathStr)
    }

    currentStruct, exists := parser.Result.Structs[currentName]
    if !exists {
        return errors.New(errors.UnknownType, currentName, path[0])
    }

    visited[currentName] = true
    path = append(path, currentName)

    for neighborName := range currentStruct.OtherStructReferences {
        if err := parser.checkPath(neighborName, path, visited); err != nil {
            return err
        }
    }

    delete(visited, currentName)

    return nil
}

func (parser *Parser) detectCycles() *errors.StackError {
    for name := range parser.Result.Structs {
        path := []string{}
        visited := make(map[string]bool)

        if err := parser.checkPath(name, path, visited); err != nil {
            return err
        }
    }

    return nil
}

func (parser *Parser) semanticAnalyze() *errors.StackError {
    for currentName, currentStruct := range parser.Result.Structs { // 3?
        for referencedName := range currentStruct.OtherStructReferences {
            targetStruct, found := parser.Result.Structs[referencedName]

            if !found {
                continue
            }

            targetStruct.EverReferenced = true
            targetStruct.ReferencedBy[currentName]++
        }
    }

    exportStruct := parser.Result.Structs[parser.Result.Exports]
    if exportStruct.EverReferenced {
        keys := make([]string, 0, len(exportStruct.ReferencedBy))
        for k := range exportStruct.ReferencedBy {
            keys = append(keys, k)
        }

        return errors.New(errors.ExportStructCantBeReferenced, parser.Result.Exports, parser.Result.Exports, strings.Join(keys, ", "))
    }

    if err := parser.detectCycles(); err != nil {
        return err
    }

    return nil
}

// Constructor
func New(myLexer *lexer.Lexer) *Parser {
    return &Parser{
        myLexer: myLexer,
        Result: types.Scheme{
            Exports: "",
            Structs: make(map[string]*types.Struct),
        },
    }
}

// Public Methods
func (parser *Parser) Parse() *errors.StackError {
    if parser.myLexer.Length() == 0 {
        return errors.New(errors.NotTokenized)
    }

    if len(parser.myLexer.ExportReferences) > 1 ||
        len(parser.myLexer.ExportReferences) == 0 ||
        parser.myLexer.TokenList[parser.myLexer.ExportReferences[0]+1].Is != types.ExportNameToken {
        return errors.New(errors.Expected1Export, len(parser.myLexer.ExportReferences))
    }

    if len(parser.myLexer.StructReferences) == 0 {
        return errors.New(errors.ExpectedStructs)
    }

    if len(parser.myLexer.FieldReferences) == 0 {
        return errors.New(errors.Expected1Field)
    }

    if err1 := parser.parseStructs(); err1 != nil {
        return err1
    }

    if err2 := parser.parseExports(); err2 != nil {
        return err2
    }

    if err3 := parser.semanticAnalyze(); err3 != nil {
        return err3
    }

    return nil
}

func (parser *Parser) Print() {
    ui.Log(config.FELLOWCRAFT, "info", "Printing parser results.")
    ui.Log(config.APPRENTICE, "info", fmt.Sprintf("Exports: %s", parser.Result.Exports))
    ui.Log(config.APPRENTICE, "info", fmt.Sprintf("Struct Count: %d", len(parser.Result.Structs)))
    ui.Log(config.APPRENTICE, "info", "STRUCTS")

    for _, _struct := range parser.Result.Structs {
        parser.printSingleStruct(_struct)
    }

    ui.Log(config.FELLOWCRAFT, "info", "Finished printing parsing results.")
}
