package errors

// Private Variables
var errorCodeToString = map[int]string{
    UnknownError: "Given error code is not recognised. Consider checking! Given code: '%d'",
    InvalidNaming: "Given name '%s' is not in preferred format. Preferred format is: ^[a-zA-Z_][a-zA-Z0-9_]",
    UnknownVerb: "Given verb for throwing an error is unknown. Consider checking. Given verb: '%s'",
    EmptyError: "%s",

    //

    PathDoesntExist: "The specified path '%s' does not exist.",
    PathIsADirectoryNotAFile: "The specified path '%s' is a directory, not a file.",
    InvalidExtension: "The file '%s' does not have the required extension '%s'.",
    SourceFileIsntValid: "The source file '%s' for copying from is not valid.",
    DestinationDirectoryIsntValid: "The destination directory '%s' for copying to is not valid or not a directory at all.",
    
    //

    UnknownOperator: "The operator '%s' at '%s' is not recognised.",
    UnexpectedTokenAtStart: "First token at a file must be a keyword always but got '%s' instead.",

    NotTokenized: "The input source code has not been tokenized yet.",
    Expected1Export: "Expected 1 export statement but got %d instead.",
    ExpectedStructs: "Expected at least 1 struct definition but got none instead.",
    Expected1Field: "Expected at least 1 field in whole file but got none instead.",
    ExportStructCantBeReferenced: "The struct '%s' which is referenced for export cannot be referenced inside the file by another struct.\nBut struct '%s' which was referenced for export was referenced by '%s'.",
    UnknownType: "The type '%s' used in struct '%s' is not recognised. Check manual.",
    //TwoStructsCantCrossReference: "Two structs can not reference each other in any way. But there was a cross reference with following path",
    AStructCantReferenceItself: "A struct cannot reference itself directly or indirectly. But struct '%s' referenced itself in given fields '%s'.",
    DidntFoundAStructToExport: "Did not find any struct to export in the source file.",
    ExpectedNameForExport: "Expected a name for export statement at '%s' but got none instead.",
    UnexpectedTokenAfterStruct: "Got unexpected token '%s' after struct definition end at '%s'.",
    AStructMustHaveAtleast1Field: "A struct must have at least 1 field defined inside it. But struct '%s' has no fields defined.",
    AnotherStructWithSameNameExists: "Another struct with same name '%s' already exists. Struct names must be unique.",
    StructShouldStartWithCurlyBrace: "A struct definition should start with a curly brace '{' but at '%s' got '%s'.",
    ExpectedNameForStruct: "Expected a name for struct definition at '%s' but it was missing or either was not in preferred format.",
    ExpectedFieldAfterAnotherField: "Expected a field definition after field '%s' at '%s' in struct '%s' since struct was not closed yet but got none instead. Either you had a typo ",
    AnotherFieldWithSameNameExists: "Another field in struct '%s' with same name '%s' already exists in. Field names must be unique inside a struct.",
    ExpectedNameForField: "Expected a name for field definition at '%s' but it was missing or either was not in preferred format.",
    BracketNotClosed: "A opened bracket at '%s' was not closed for defining slice type. Got '%s' at '%s' instead of closing bracket ']'.",
    ExpectedMapDefinition: "Expected a map keyword at '%s' because there was '{type} before it.",
    CurlyBraceNotClosed: "A opened curly brace at '%s' was not closed for defining struct or map. Consider checking.",
    NoTypeSpecifiedForMap: "No type was specified after '{' at '%s'. If a type for a field starts with curly brace it is considered as a map type.",
    ExpectedAValidType: "Expected a valid type at '%s' but got '%s' instead. Consider checking manual for valid types.",
    InvalidStructNaming: "The struct defined at '%s' with name '%s' can not have that name since that name is a default type.",
    UnexpectedTokenAfterField: "Got unexpected token '%s' after field definition at '%s'.",
    CyclicReference: "Cyclic reference detected, path is: '%s'.",
}

// Public Constants
const (
    UnknownError int = iota
    InvalidNaming
    UnknownVerb
    EmptyError

    //

    PathDoesntExist
    PathIsADirectoryNotAFile
    InvalidExtension
    SourceFileIsntValid
    DestinationDirectoryIsntValid

    //

    UnknownOperator
    UnexpectedTokenAtStart

    NotTokenized
    Expected1Export
    ExpectedStructs
    Expected1Field
    ExportStructCantBeReferenced
    UnknownType
    AStructCantReferenceItself
    DidntFoundAStructToExport
    ExpectedNameForExport
    UnexpectedTokenAfterStruct
    AStructMustHaveAtleast1Field
    AnotherStructWithSameNameExists
    StructShouldStartWithCurlyBrace
    ExpectedNameForStruct
    ExpectedFieldAfterAnotherField
    AnotherFieldWithSameNameExists
    ExpectedNameForField
    BracketNotClosed
    ExpectedMapDefinition
    CurlyBraceNotClosed
    NoTypeSpecifiedForMap
    ExpectedAValidType
    InvalidStructNaming
    UnexpectedTokenAfterField
    CyclicReference
)
