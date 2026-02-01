package language

// Public Variables
var Keywords = map[string]bool{
    "struct":  true,
    "field":   true,
    "exports": true,
}

var Operators = map[string]bool{"{": true, "}": true, "[": true, "]": true}

var DefaultTypes = map[string]bool{
    // Integers
    "u8":   true,   "i8":   true,
    "u16":  true,   "i16":  true,
    "u32":  true,   "i32":  true,

    // Floats
    "f16":  true,   "f24":  true,
    "f32":  true,   "f64":  true,

    // Vectors
    "vector2":      true,   "vector3":      true,   // f32 * 2, f32 * 3
    "vector2int16": true,   "vector3int16": true,   // i16 * 2, i16 * 3
    "vector2norm":  true,   "vector3norm":  true,   // 2 bytes, 4 bytes

    // CFrames
    "cframe":   true,   // f32 * 3 + f32 * 9
    "cframe_e": true,   // f32 * 3 + i16 * 3
    "cframe_q": true,   // f32 * 3 + f16 * 4

    // Color3
    "color3":       true,   // u8 * 3
    "color3_hdr":   true,   // f16 * 3

    // Other
    "string":   true, // u8 length + data
    "string_l": true, // u16 length + data
    "bool":     true, // Default -> 1 byte, In arrays -> 1 bit
}

var DefaultTypesToRobloxTypes = map[string]string{
    "u8":           "number",
    "i8":           "number",
    "u16":          "number",
    "i16":          "number",
    "u32":          "number",
    "i32":          "number",
    "f16":          "number",
    "f24":          "number",
    "f32":          "number",
    "f64":          "number",
   
    "vector2":      "Vector2",
    "vector3":      "Vector3",
    "vector2int16": "Vector2int16",
    "vector3int16": "Vector3int16",
    "vector2norm":  "Vector2",
    "vector3norm":  "Vector3",

    "cframe":       "CFrame",
    "cframe_e":     "CFrame",
    "cframe_q":     "CFrame",

    "color3":       "Color3",
    "color3_hdr":   "Color3",

    "string":       "string",
    "string_l":     "string",
    "bool":         "boolean",
}
