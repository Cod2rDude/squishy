/*
    ******************************************************************************
    * @file     : squishy/squishy-compiler/internal/app/middleend/middleend.go
    * @author   : Cod2rDude
    * @date     : January 23 2026
    * @lastEdit : January 26 2025 @ 09:24
    * @brief    : Squishy IDL Compiler Middleend.
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

package middleend

import (
    "fmt"
    "strings"

    "github.com/Cod2rDude/squishy/squishy-compiler/internal/app/language"
    "github.com/Cod2rDude/squishy/squishy-compiler/internal/app/types"
    "github.com/Cod2rDude/squishy/squishy-compiler/internal/cli/ui"
    "github.com/Cod2rDude/squishy/squishy-compiler/internal/config"
    "github.com/Cod2rDude/squishy/squishy-compiler/internal/errors"
)

// Function
func getSizeForField(field *types.Field) int {
    size := 10 // Operators ':' ';' and 2 spaces before and after ':' also 4 spaces at start, EOL
    size += len(field.Name)
    size += len(field.Type.Name)

    if field.Type.IsArray || field.Type.IsMap { // { [number] : type } { [string] : type }
        size += 15
    }

    return size
}

func getFieldTypeString(field *types.Field) string {
    out := "    " + field.Name + " : "

    typeName := field.Type.Name

    if _, isDefault := language.DefaultTypes[typeName]; isDefault {
        typeName = language.DefaultTypesToRobloxTypes[typeName]
    }

    if field.Type.IsArray {
        out += "{ [number] : " + typeName + " }"
    } else if field.Type.IsMap {
        out += "{ [string] : " + typeName + " }"
    } else {
        out += typeName
    }

    out += ";\n"

    return out
}

// Public Structs

/*
    @object Middleend

    @privatevariables
    *   @privatevariable sortedStructs : []string ;; List of structs sorted in dependency-reference connection.
    *   @privatevariable scheme : *types.Scheme ;; Pointer to scheme created by frontend.
    *   @privatevariable exportBuilder : strings.Builder ;; String builder for lua export type.
    *   @privatevariable typeBuilder : strings.Builder ;; String builder for lua type.
    @privatemethods
    *   @privatemethod noteStructsToCareAbout
    *   @privatemethod sortStructs
    *   @privatemethod writeType
    *   @privatemethod writeTypes
    *   @privatemethod writeExport
    @publicmethods
    *   @publicmethod Work
    *   @publicmethod GetResults
    *   @publicmethod GetSortedStructs
    *   @publicmethod Debug
    @brief Middleend of Squishy IDL Compiler.
*/
type Middleend struct {
    sortedStructs []string
    scheme        *types.Scheme
    exportBuilder strings.Builder
    typeBuilder   strings.Builder
}

// Constructor
func New(scheme *types.Scheme) *Middleend {
    return &Middleend{
        sortedStructs: []string{},
        scheme:        scheme,
        exportBuilder: strings.Builder{},
        typeBuilder:   strings.Builder{},
    }
}

// Private Methods
func (middleend *Middleend) noteStructsToCareAbout() []string {
    notedStructs := []string{}

    for name, _struct := range middleend.scheme.Structs {
        if name == middleend.scheme.Exports || !_struct.EverReferenced {
            continue
        }

        notedStructs = append(notedStructs, name)
    }

    return notedStructs
}

func (middleend *Middleend) sortStructs() { // Not likely
    // Sort structs in order for later conversion to types.
    // Order means their order in file, its based on who requires who.
    // If a struct needs another struct, its below that another struct.
    // Topological Sorting??
    // This must run after semantic analysis (cyclic dependency check)

    notedStructs := middleend.noteStructsToCareAbout()

    visited := make(map[string]bool)
    sortedStructs := make([]string, 0, len(notedStructs))

    var visit func(name string)
    visit = func(name string) {
        if visited[name] {
            return
        }

        if _struct, exists := middleend.scheme.Structs[name]; exists {
            for dependencyName := range _struct.OtherStructReferences {
                visit(dependencyName)
            }
        }

        visited[name] = true
        sortedStructs = append(sortedStructs, name)
    }

    for _, name := range notedStructs {
        visit(name)
    }

    middleend.sortedStructs = sortedStructs
}

func (middleend *Middleend) writeType(name string) {
    fetchedStruct, _ := middleend.scheme.Structs[name]

    expectedSize := 12 //type  = {}\n
    expectedSize += len(name)
    expectedSize += 2 // EOL

    for _, field := range fetchedStruct.Fields {
        expectedSize += getSizeForField(field)
    }

    middleend.typeBuilder.Grow(expectedSize)

    middleend.typeBuilder.WriteString(fmt.Sprintf("type %s = {\n", fetchedStruct.Name))

    for _, field := range fetchedStruct.Fields {
        middleend.typeBuilder.WriteString(getFieldTypeString(field))
    }

    middleend.typeBuilder.WriteString("}\n")
}

func (middleend *Middleend) writeTypes() *errors.StackError {
    for i, name := range middleend.sortedStructs {
        middleend.writeType(name);

        if i != len(middleend.sortedStructs)-1 {
            middleend.typeBuilder.WriteString("\n")
        }
    }

    return nil
}

func (middleend *Middleend) writeExport()  {
    exportStruct := middleend.scheme.Structs[middleend.scheme.Exports]

    expectedSize := 19 //export type  = {}\n
    expectedSize += len(exportStruct.Name)
    expectedSize += 2 // EOL

    for _, field := range exportStruct.Fields {
        expectedSize += getSizeForField(field)
    }

    middleend.exportBuilder.Grow(expectedSize)

    middleend.exportBuilder.WriteString(fmt.Sprintf("export type %s = {\n", exportStruct.Name))

    for _, field := range exportStruct.Fields {
        middleend.exportBuilder.WriteString(getFieldTypeString(field))
    }

    middleend.exportBuilder.WriteString("}\n")
}

// Public Methods
func (middleend *Middleend) Work() {
    middleend.sortStructs()
    middleend.writeTypes()
    middleend.writeExport();
}

func (middleend *Middleend) GetResults() (string, string) {
    return middleend.typeBuilder.String(), middleend.exportBuilder.String()
}

func (middleend *Middleend) GetSortedStructs() []string {
    return middleend.sortedStructs
}

func (middleend *Middleend) Debug() {
    exportString := middleend.exportBuilder.String()
    exportStringByte := []byte(exportString)
    exportStringByte = exportStringByte[:len(exportStringByte)-1]
    exportString = strings.ReplaceAll(string(exportStringByte), "\n", "\n"+strings.Repeat(" ", config.APPRENTICE+7))

    typeString := middleend.typeBuilder.String()
    typeStringByte := []byte(typeString)
    typeStringByte = typeStringByte[:len(typeStringByte)-1]
    typeString = strings.ReplaceAll(string(typeStringByte), "\n", "\n"+strings.Repeat(" ", config.APPRENTICE+7))

    ui.Log(config.MASTER, "info", "MIDDLEEND DEBUG START")
    ui.Log(config.FELLOWCRAFT, "info", fmt.Sprintf("Sorted order of structs: '%s'", strings.Join(middleend.sortedStructs, ", ")))
    ui.Log(config.FELLOWCRAFT, "info", "EXPORT TYPE OUTPUT")
    ui.Log(config.APPRENTICE, "info", exportString)
    ui.Log(config.FELLOWCRAFT, "info", "TYPE OUTPUT")
    ui.Log(config.APPRENTICE, "info", typeString)
    ui.Log(config.MASTER, "info", "MIDDLEEND DEBUG END")
}
