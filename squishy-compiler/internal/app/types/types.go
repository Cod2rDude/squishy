package types

// Public Constants
const (
	UnknownToken = iota
	KeywordToken
	TypeToken
	StructNameToken
	IntToken
	OperatorToken
	InvalidToken
	FieldNameToken
	ExportNameToken
	CommentToken
)

// Public Structs
type Token struct {
	Position     int
	Is           int
	Value        string
	RealPosition string
}

type Type struct {
	Name                        string
	IsArray                     bool
	ArraySize                   int
	IsMap                       bool
	IsShortMap                  bool
	IsReferenceToAnotherStruct  bool
}

type Field struct {
	Name string
	Type *Type
}

type Struct struct {
	Reference             string
	Name                  string
	Fields                []*Field
	OtherStructReferences map[string][]int
	EverReferenced        bool
	ReferencedBy          map[string]int
}

type Scheme struct {
	Structs map[string]*Struct
	Exports string
}
