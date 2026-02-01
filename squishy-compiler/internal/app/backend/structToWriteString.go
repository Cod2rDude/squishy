package backend

import (
    "strings"

    "github.com/Cod2rDude/squishy/squishy-compiler/internal/app/language"
    "github.com/Cod2rDude/squishy/squishy-compiler/internal/app/types"
)

// Functions
func getWriteStringForArray(name string, _type *types.Type) string {
    if !_type.IsArray {
        return ""
    }

    if _type.Name == "bool" {
        if _type.ArraySize <= 0 {
            return format("cursor = writer.write_dynamicBoolArray(sharedBuffer, cursor, input.%s)", name)
        } else {
            return format("cursor = writer.write_boolArray(sharedBuffer, cursor, input.%s, %d)", name, _type.ArraySize)
        }
    }

    writeFunction := ""

    if language.DefaultTypes[_type.Name] {
        writeFunction = "writer.write_" + _type.Name
    } else {
        writeFunction = "write_" + _type.Name
    }

    if _type.ArraySize <= 0 {
        return format("cursor = writer.write_dynamicArray(sharedBuffer, cursor, input.%s, %s)", name, writeFunction)
    } else {
        return format("cursor = writer.write_array(sharedBuffer, cursor, input.%s, %s, %d)", name, writeFunction, _type.ArraySize)
    }
}

func getWriteStringForMap(name string, _type *types.Type) string {
    if !_type.IsMap {
        return ""
    }

    writeFunction := ""

    if language.DefaultTypes[_type.Name] {
        writeFunction = "writer.write_" + _type.Name
    } else {
        writeFunction = "write_" + _type.Name
    }

    shortMapString := ""

    if _type.IsShortMap {
        shortMapString = "s"
    }

    out := format("cursor = writer.write_%smap(sharedBuffer, cursor, input.%s, %s)", shortMapString, name, writeFunction)

    return out
}

func getWriteStringForADefaultType(name string, _type *types.Type) string {
    return format("cursor = writer.write_%s(sharedBuffer, cursor, input.%s)", _type.Name, name)
}

func getWriteStringForField(field *types.Field) string {
    if field.Type.IsArray {
        return getWriteStringForArray(field.Name, field.Type)
    } else if field.Type.IsMap {
        return getWriteStringForMap(field.Name, field.Type)
    }

    if language.DefaultTypes[field.Type.Name] {
        return getWriteStringForADefaultType(field.Name, field.Type)
    }

    return format("cursor = write_%s(cursor, input.%s)", field.Type.Name, field.Name)
}

// Public Functions
func StructToWriteString(_struct *types.Struct) []string {
    out := []string{}

    fields := _struct.Fields

    for _, val := range fields {
        out = append(out, getWriteStringForField(val))
    }

    return out
}

func StructListToWriteFunctions(list []string, structs map[string]*types.Struct) []string {
    out := []string{}

    for _, name := range list {
        out = append(out, format("function write_%s(cursor : number, input : %s) : number", name, name))
        out = append(out, "    "+strings.Join(StructToWriteString(structs[name]), "\n    "))
        out = append(out, "    return cursor")
        out = append(out, "end\n")
    }

    return out
}
