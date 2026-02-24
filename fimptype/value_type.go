package fimptype

type ValueTypeT string

const (
	VTypeString     ValueTypeT = "string"
	VTypeInt        ValueTypeT = "int"
	VTypeFloat      ValueTypeT = "float"
	VTypeBool       ValueTypeT = "bool"
	VTypeStrMap     ValueTypeT = "str_map"
	VTypeIntMap     ValueTypeT = "int_map"
	VTypeFloatMap   ValueTypeT = "float_map"
	VTypeBoolMap    ValueTypeT = "bool_map"
	VTypeStrArray   ValueTypeT = "str_array"
	VTypeIntArray   ValueTypeT = "int_array"
	VTypeFloatArray ValueTypeT = "float_array"
	VTypeBoolArray  ValueTypeT = "bool_array"
	VTypeObject     ValueTypeT = "object"
	VTypeBase64     ValueTypeT = "base64"
	VTypeBinary     ValueTypeT = "bin"
	VTypeNull       ValueTypeT = "null"
)

func (vt ValueTypeT) Str() string {
	return string(vt)
}
