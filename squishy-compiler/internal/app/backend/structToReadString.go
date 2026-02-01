package backend

import (
    "fmt"
    "strings"

    "github.com/Cod2rDude/squishy/squishy-compiler/internal/app/language"
    "github.com/Cod2rDude/squishy/squishy-compiler/internal/app/types"
)

// Functions
func format(str string, a ...any) string {
    return fmt.Sprintf(str, a...)
}

func getReadStringForArray(name string, _type *types.Type) string {
    if !_type.IsArray {
        return ""
    }

    if _type.Name == "bool" {
        if _type.ArraySize <= 0 {
            return format("cursor, %s = reader.read_dynamicBoolArray(sharedBuffer, cursor)", name)
        } else {
            return format("cursor, %s = reader.read_boolArray(sharedBuffer, cursor, %d)", name, _type.ArraySize)
        }
    }

    readFunction := ""

    if language.DefaultTypes[_type.Name] {
        readFunction = "reader.read_" + _type.Name
    } else {
        readFunction = "read_" + _type.Name
    }

    if _type.ArraySize <= 0 {
        return format("cursor, %s = reader.read_dynamicArray(buff, cursor, %s)", name, readFunction)
    } else {
        return format("cursor, %s = reader.read_array(buff, cursor, %s, %d)", name, readFunction, _type.ArraySize, )
    }
}

func getReadStringForMap(name string, _type *types.Type) string {
    if !_type.IsMap {
        return ""
    }

    readFunction := ""

    if language.DefaultTypes[_type.Name] {
        readFunction = "reader.read_" + _type.Name
    } else {
        readFunction = "read_" + _type.Name
    }

    shortMapString := ""

    if _type.IsShortMap {
        shortMapString = "s"
    }
    
    out := format("cursor, %s = reader.read_%smap(buff, cursor, %s)", name, shortMapString, readFunction)

    return out
}

func getReadStringForADefaultType(name string, _type *types.Type) string {
    return format("cursor, %s = reader.read_%s(buff, cursor)", name, _type.Name)
}

func getReadStringForField(field *types.Field) string {
    if field.Type.IsArray {
        return getReadStringForArray(field.Name, field.Type)
    } else if field.Type.IsMap {
        return getReadStringForMap(field.Name, field.Type)
    }

    if language.DefaultTypes[field.Type.Name] {
        return getReadStringForADefaultType(field.Name, field.Type)
    }

    return format("cursor, %s = read_%s(buff, cursor)", field.Name, field.Type.Name)
}

// Public Functions
func StructToReadString(_struct *types.Struct) ([]string, string) {
    fields := _struct.Fields

    returnString := "{ "
    fieldNames := []string{}

    for _, field := range fields {
        fieldNames = append(fieldNames, field.Name)
        returnString = format("%s%s = %s; ", returnString, field.Name, field.Name)
    }

    variableString := "local " + strings.Join(fieldNames, ", ")

    returnString = returnString + "}"

    out := []string{variableString}

    for _, val := range fields {
        out = append(out, getReadStringForField(val))
    }

    return out, returnString
}

func StructListToReadFunctions(list []string, structs map[string]*types.Struct) []string {
    out := []string{}

    for _, name := range list {
        out = append(out, format("function read_%s(buff : buffer, cursor : number) : (number, %s)", name, name))
        body, returnString := StructToReadString(structs[name])
        out = append(out, "    "+strings.Join(body, "\n    "))
        out = append(out, "    return cursor, "+returnString)
        out = append(out, "end\n")
    }

    return out
}
