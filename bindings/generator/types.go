package main

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

type Param struct {
	Name string
	Type IType
}

type IType interface{}
type TypePrimitive uint32
type TypeStruct struct {
	Fields []Param
}
type TypeArray struct {
	Count int
	Type  IType
}
type TypeSlice struct {
	Type IType
}
type TypePointer struct {
	Type IType
}
type TypeFunc struct {
	Input  []Param
	Output []Param
}
type TypeIdent struct {
	Type IType
	Name string
}

const (
	TypeP_void TypePrimitive = iota
	TypeP_u8
	TypeP_u16
	TypeP_u32
	TypeP_u64

	TypeP_i8
	TypeP_i16
	TypeP_i32
	TypeP_i64

	TypeP_f32
	TypeP_f64

	TypeP_string
	TypeP_bool
	TypeP_error

	TypeP_WasmPointer = TypeP_u32
)

func LowerParam(p Param) ([]Param, bool) {
	name := p.Name
	switch e := p.Type.(type) {
	case TypePrimitive:
		p.Name = fmt.Sprintf("%s_w", name)
		switch e {
		case TypeP_string:
			return []Param{
				{fmt.Sprintf("%s_str", name), TypeP_WasmPointer},
				{fmt.Sprintf("%s_len", name), TypeP_u32},
			}, true
		case TypeP_bool:
			return []Param{{fmt.Sprintf("%s_b", name), TypeP_u32}}, true
		case TypeP_u8, TypeP_u16, TypeP_u32:
			p.Type = TypeP_u32
			return []Param{p}, true
		case TypeP_i8, TypeP_i16, TypeP_i32:
			p.Type = TypeP_i32
			return []Param{p}, true
		case TypeP_u64:
			p.Type = TypeP_u64
			return []Param{p}, true
		case TypeP_i64:
			p.Type = TypeP_i64
			return []Param{p}, true
		case TypeP_error:
			return []Param{}, true
		}

		return []Param{p}, true
	case TypeArray:
		return []Param{
			{fmt.Sprintf("%s_ptr", name), TypeP_WasmPointer},
		}, true
	case TypeSlice:
		return []Param{
			{fmt.Sprintf("%s_ptr", name), TypeP_WasmPointer},
			{fmt.Sprintf("%s_len", name), TypeP_u32},
		}, true
	case TypeFunc:
		return []Param{
			{fmt.Sprintf("%s_fn", name), TypeP_WasmPointer},
		}, true
	case TypePointer:
		return []Param{
			{fmt.Sprintf("%s_ptr", name), TypeP_WasmPointer},
		}, true
	case TypeStruct:
		res := []Param{}
		for _, f := range e.Fields {
			f.Name = fmt.Sprintf("%s_%s_f", name, f.Name)
			res = append(res, f)
		}
		return res, false
	case TypeIdent:
		p.Type = e.Type
		return LowerParam(p)
	default:
		panic(fmt.Sprintf("invalid type type: %v", p.Type))
	}
}
func FullyLowerParams(p []Param) []Param {
	to_be_lowered := p
	result := []Param{}

	for len(to_be_lowered) > 0 {
		next_to_be_lowered := []Param{}
		for _, p := range to_be_lowered {
			r, is_low := LowerParam(p)
			if is_low {
				result = append(result, r...)
			} else {
				next_to_be_lowered = append(next_to_be_lowered, r...)
			}
		}

		to_be_lowered = next_to_be_lowered
	}

	return result
}
func ToUpperFirst(s string) string {
	if s == "" {
		return s
	}

	r, size := utf8.DecodeRuneInString(s)
	return string(unicode.ToUpper(r)) + s[size:]
}
func HyphenCaseToPascalCase(s string) string {
	res := strings.Builder{}

	for _, s := range strings.Split(s, "-") {
		res.WriteString(ToUpperFirst(strings.ToLower(s)))
	}

	return res.String()
}

func Writef(b *strings.Builder, format string, args ...any) {
	b.WriteString(fmt.Sprintf(format, args...))
}
func Writes(b *strings.Builder, msg string) {
	b.WriteString(msg)
}

type ILangBuilder interface {
	WriteGuestEncode(b *strings.Builder, p Param)
	WriteGuestDecode(b *strings.Builder, p Param)

	StringType(t IType) string

	StringDeclare(p Param) string
	StringAssign(dst Param, format string, args ...any) string
	StringCall(name string, args []Param, rets []Param) string
	StringResetMalloc() string
	StringReturn(p []Param) string

	StringFnHeader(name string, fn TypeFunc) string
	StringWasmImport(module string, name string, bind string, fn TypeFunc) string

	StringBlockBegin() string
	StringBlockEnd() string
}

type ConfigHostWriter struct {
	Golang        *GoLangBuilder
	NameApiModule string
	NameFnAlloc   string
}

var Type_AllocFn = TypeFunc{Input: []Param{{"size", TypeP_u32}}, Output: []Param{{"offset", TypeP_WasmPointer}}}

func WriteGoHostEncode(config *ConfigHostWriter, b *strings.Builder, p Param) {
	g := config.Golang
	m := config.NameApiModule
	alloc := config.NameFnAlloc
	pr, _ := LowerParam(p)

	name := p.Name
	switch e := p.Type.(type) {
	case TypePrimitive:
		switch e {
		case TypeP_string:
			Writef(b, "uWriteString(%s,%s,%s,%s,%s);", m, alloc, name, pr[0].Name, pr[1].Name)
			return
		case TypeP_bool:
			Writef(b, "if %s{", name)
			Writef(b, "uWriteAnyFixed[%s](%s,%s,%s)", g.StringType(pr[0].Type), m, "1", pr[0].Name)
			Writes(b, "}else{")
			Writef(b, "uWriteAnyFixed[%s](%s,%s,%s)", g.StringType(pr[0].Type), m, "0", pr[0].Name)
			Writes(b, "};")
			return
		}

		Writef(b, "uWriteAnyFixed[%s](%s,%s,%s);", g.StringType(pr[0].Type), m, name, pr[0].Name)
		return
	case TypeArray:
		Writef(b, "uWritePlex(%s,%s,%s,%s);", m, alloc, name, pr[0].Name)
		return
	case TypeSlice:
		Writef(b, "uWriteSlice(%s,%s,%s,%s,%s);", m, alloc, name, pr[0].Name, pr[1].Name)
		return
	case TypeFunc:
		panic("can't pass function through wasm abi")
	case TypePointer:
		Writef(b, "uWritePointer(%s,%s,%s,%s);", m, alloc, name, pr[0].Name)
		return
	case TypeStruct:
		for i, f := range e.Fields {
			Writes(b, g.StringAssign(pr[i], "%s.%s", name, f.Name))
			WriteGoHostEncode(config, b, pr[i])
		}
		return
	case TypeIdent:
		p.Type = e.Type
		WriteGoHostEncode(config, b, p)
		return
	default:
		panic(fmt.Sprintf("invalid type type: %v", p.Type))
	}
	// panic(fmt.Sprintf("unreachable code: %v", p.Type))
}
func WriteGoHostDecode(config *ConfigHostWriter, b *strings.Builder, p Param) {
	g := config.Golang
	m := config.NameApiModule
	pr, _ := LowerParam(p)

	name := p.Name
	switch e := p.Type.(type) {
	case TypePrimitive:
		switch e {
		case TypeP_string:
			Writes(b, g.StringDeclare(p))
			Writes(b, g.StringAssign(p, "uReadString(%s, %s, %s)", m, pr[0].Name, pr[1].Name))
			return
		case TypeP_bool:
			Writes(b, g.StringDeclare(p))
			Writes(b, g.StringAssign(p, "%s != 0", pr[0].Name))
			return
		}

		Writes(b, g.StringDeclare(p))
		Writes(b, g.StringAssign(p, "%s", pr[0].Name))
		return
	case TypeArray:
		Writes(b, g.StringDeclare(p))
		Writes(b, g.StringAssign(p, "uReadPlex[%s](%s, %s)", g.StringType(p.Type), m, pr[0].Name))
		return
	case TypeSlice:
		Writes(b, g.StringDeclare(p))
		Writes(b, g.StringAssign(p, "uReadSlice[%s](%s, %s, %s)", g.StringType(p.Type), m, pr[0].Name, pr[1].Name))
		return
	case TypeFunc:
		panic("can't pass function through wasm abi")
	case TypePointer:
		Writes(b, g.StringDeclare(p))
		Writes(b, g.StringAssign(p, "uReadPointer[%s](%s, %s)", g.StringType(p.Type), m, pr[0].Name))
		return
	case TypeStruct:
		Writes(b, g.StringDeclare(p))
		for i, f := range e.Fields {
			WriteGoHostDecode(config, b, pr[i])
			Writef(b, "%s.%s=%s;", name, f.Name, pr[i].Name)
		}
		return
	case TypeIdent:
		p.Type = e.Type
		WriteGoHostDecode(config, b, p)
		return
	default:
		panic(fmt.Sprintf("invalid type type: %v", p.Type))
	}
	// panic(fmt.Sprintf("unreachable code: %v", p.Type))
}
func IsUseMalloc(p Param) bool {
	switch e := p.Type.(type) {
	case TypePrimitive:
		switch e {
		case TypeP_string:
			return true
		case TypeP_bool:
			return false
		}

		return false
	case TypeArray:
		return true
	case TypeSlice:
		return true
	case TypeFunc:
		panic("can't pass function through wasm abi")
	case TypePointer:
		return true
	case TypeStruct:
		return true
	case TypeIdent:
		p.Type = e.Type
		return IsUseMalloc(p)
	default:
		panic(fmt.Sprintf("invalid type type: %v", p.Type))
	}
	// panic(fmt.Sprintf("unreachable code: %v", p.Type))
}
