package main

import (
	"fmt"
	"strings"
)

type GoLangBuilder struct{}

func (g *GoLangBuilder) WriteGuestEncode(b *strings.Builder, p Param) {
	pr, _ := LowerParam(p)

	if len(pr) == 0 {
		return
	}

	name := p.Name
	switch e := p.Type.(type) {
	case TypePrimitive:
		switch e {
		case TypeP_string:
			Writes(b, g.StringDeclare(pr[0]))
			Writes(b, g.StringDeclare(pr[1]))
			Writes(b, g.StringAssign(pr[0], "uintptr(unsafe.Pointer(unsafe.StringData(%s)))", name))
			Writes(b, g.StringAssign(pr[1], "len(%s)", name))
			return
		case TypeP_bool:
			Writes(b, g.StringDeclare(pr[0]))
			Writef(b, "if %s{", name)
			Writes(b, g.StringAssign(pr[0], "1"))
			Writef(b, "}else{")
			Writes(b, g.StringAssign(pr[0], "0"))
			Writef(b, "};")
			return
		}

		Writes(b, g.StringDeclare(pr[0]))
		Writes(b, g.StringAssign(pr[0], "%s", name))
		return
	case TypeArray:
		Writes(b, g.StringDeclare(pr[0]))
		Writef(b, "if len(%s)>0{", name)
		Writes(b, g.StringAssign(pr[0], "uintptr(unsafe.Pointer(&%s[0]))", name))
		Writef(b, "};")
		return
	case TypeSlice:
		Writes(b, g.StringDeclare(pr[0]))
		Writes(b, g.StringDeclare(pr[1]))
		Writef(b, "if len(%s)>0{", name)
		Writes(b, g.StringAssign(pr[0], "uintptr(unsafe.Pointer(&%s[0]))", name))
		Writef(b, "};")
		Writes(b, g.StringAssign(pr[1], "len(%s)", name))
		return
	case TypeFunc:
		Writes(b, g.StringDeclare(pr[0]))
		Writes(b, g.StringAssign(pr[0], "uintptr(unsafe.Pointer(&%s))", name))
		return
	case TypePointer:
		Writes(b, g.StringDeclare(pr[0]))
		Writes(b, g.StringAssign(pr[0], "uintptr(unsafe.Pointer(%s))", name))
		return
	case TypeStruct:
		for i, f := range e.Fields {
			Writes(b, g.StringDeclare(pr[i]))
			Writes(b, g.StringAssign(pr[i], "%s.%s", name, f.Name))
			g.WriteGuestEncode(b, pr[i])
		}
		return
	case TypeIdent:
		p.Type = e.Type
		g.WriteGuestEncode(b, p)
		return
	default:
		panic(fmt.Sprintf("invalid type type: %v", p.Type))
	}
	// panic(fmt.Sprintf("unreachable code: %v", p.Type))
}

func (g *GoLangBuilder) WriteGuestDecode(b *strings.Builder, p Param) {
	pr, _ := LowerParam(p)

	if len(pr) == 0 {
		return
	}

	Writes(b, g.StringDeclare(p))

	name := p.Name
	switch e := p.Type.(type) {
	case TypePrimitive:
		switch e {
		case TypeP_string:
			Writes(b, g.StringAssign(p, "unsafe.String((*byte)(unsafe.Pointer(uintptr(%s))),int(%s))", pr[0].Name, pr[1].Name))
			return
		case TypeP_bool:
			Writes(b, g.StringAssign(p, "%s!=0", pr[0].Name))
			return
		}

		Writes(b, g.StringAssign(p, "%s", pr[0].Name))
		return
	case TypeArray:
		Writef(b, "%s=*(*%s)(unsafe.Pointer(uintptr(%s)));", name, g.StringType(p.Type), pr[0].Name)
		return
	case TypeSlice:
		Writes(b, g.StringAssign(p, "unsafe.Slice((*%s)(unsafe.Pointer(uintptr(%s)),int(%s)))", g.StringType(e.Type), pr[0].Name, pr[1].Name))
		return
	case TypeFunc:
		Writef(b, "%s=*(*%s)(unsafe.Pointer(uintptr(%s)));", name, g.StringType(p.Type), pr[0].Name)
		return
	case TypePointer:
		Writef(b, "%s=*(*%s)(unsafe.Pointer(uintptr(%s)));", name, g.StringType(p.Type), pr[0].Name)
		return
	case TypeStruct:
		for i, f := range e.Fields {
			g.WriteGuestDecode(b, pr[i])
			Writef(b, "%s.%s=%s;", name, f.Name, pr[i].Name)
		}
		return
	case TypeIdent:
		p.Type = e.Type
		g.WriteGuestDecode(b, p)
		return
	default:
		panic(fmt.Sprintf("invalid type type: %v", p.Type))
	}
	// panic(fmt.Sprintf("unreachable code: %v", p.Type))
}

func (g *GoLangBuilder) StringType(t IType) string {
	switch e := t.(type) {
	case TypePrimitive:
		switch e {
		case TypeP_void:
			return "struct{}"
		case TypeP_u8:
			return "uint8"
		case TypeP_u16:
			return "uint16"
		case TypeP_u32:
			return "uint32"
		case TypeP_u64:
			return "uint64"
		case TypeP_i8:
			return "int8"
		case TypeP_i16:
			return "int16"
		case TypeP_i32:
			return "int32"
		case TypeP_i64:
			return "int64"
		case TypeP_f32:
			return "float32"
		case TypeP_f64:
			return "float64"
		case TypeP_string:
			return "string"
		case TypeP_bool:
			return "bool"
		case TypeP_error:
			return "error"
		}
	case TypeArray:
		return fmt.Sprintf("[%d]%s", e.Count, g.StringType(e.Type))
	case TypeSlice:
		return fmt.Sprintf("[]%s", g.StringType(e.Type))
	case TypeFunc:
		input := make([]string, len(e.Input))
		output := make([]string, len(e.Output))

		for i, val := range e.Input {
			input[i] = fmt.Sprintf("%s %s", val.Name, g.StringType(val.Type))
		}
		for i, val := range e.Output {
			output[i] = fmt.Sprintf("%s %s", val.Name, g.StringType(val.Type))
		}

		return fmt.Sprintf("func(%s)(%s)", strings.Join(input, ","), strings.Join(output, ","))
	case TypePointer:
		return fmt.Sprintf("*%s", g.StringType(e.Type))
	case TypeStruct:
		fields := make([]string, len(e.Fields))
		for i, val := range e.Fields {
			fields[i] = fmt.Sprintf("%s %s", val.Name, g.StringType(val.Type))
		}
		return fmt.Sprintf("struct{%s}", strings.Join(fields, ";"))
	case TypeIdent:
		return e.Name
	}
	panic(fmt.Sprintf("invalid type type: %v", t))
}

func (g *GoLangBuilder) StringDeclare(p Param) string {
	return fmt.Sprintf("var %s %s;", p.Name, g.StringType(p.Type))
}

func (g *GoLangBuilder) StringAssign(dst Param, format string, args ...any) string {
	return fmt.Sprintf("%s=(%s)(%s);", dst.Name, g.StringType(dst.Type), fmt.Sprintf(format, args...))
}

func (g *GoLangBuilder) StringCall(name string, args []Param, rets []Param) string {
	str := strings.Builder{}
	for i, e := range args {
		str.WriteString(e.Name)
		if i < len(args)-1 {
			str.WriteString(",")
		}
	}
	if len(rets) > 0 {
		str.WriteString(",")
		for i, e := range rets {
			Writef(&str, "(%s)(uintptr(unsafe.Pointer(&%s)))", g.StringType(e.Type), e.Name)
			if i < len(rets)-1 {
				str.WriteString(",")
			}
		}
	}
	return fmt.Sprintf("%s(%s);", name, str.String())
}

func (g *GoLangBuilder) StringReturn(p []Param) string {
	args := strings.Builder{}
	for i, e := range p {
		args.WriteString(e.Name)
		if i < len(p)-1 {
			args.WriteString(",")
		}
	}
	return fmt.Sprintf("return %s;", args.String())
}

func (g *GoLangBuilder) StringFnHeader(name string, fn TypeFunc) string {
	params := strings.Builder{}
	for i, e := range fn.Input {
		Writef(&params, "%s %s", e.Name, g.StringType(e.Type))
		if i < len(fn.Input)-1 {
			params.WriteString(",")
		}
	}
	rets := strings.Builder{}
	if len(fn.Output) > 0 {
		rets.WriteString("(")
		for i, e := range fn.Output {
			Writef(&rets, "%s", g.StringType(e.Type))
			if i < len(fn.Output)-1 {
				rets.WriteString(",")
			}
		}
		rets.WriteString(")")
	}
	return fmt.Sprintf("func %s(%s)%s", name, params.String(), rets.String())
}

func (g *GoLangBuilder) StringWasmImport(module string, name string, bind_name string, fn TypeFunc) string {
	params := strings.Builder{}
	for i, e := range fn.Input {
		Writef(&params, "%s %s", e.Name, g.StringType(e.Type))
		if i < len(fn.Input)-1 {
			params.WriteString(",")
		}
	}
	rets := strings.Builder{}
	if len(fn.Output) > 0 {
		rets.WriteString("(")
		for i, e := range fn.Output {
			Writef(&rets, "%s", g.StringType(e.Type))
			if i < len(fn.Output)-1 {
				rets.WriteString(",")
			}
		}
		rets.WriteString(")")
	}
	return fmt.Sprintf("//go:wasmimport %s %s\nfunc %s(%s)%s;\n", module, name, bind_name, params.String(), rets.String())
}

func (g *GoLangBuilder) StringResetMalloc() string {
	return "uResetMalloc();"
}

func (g *GoLangBuilder) StringBlockBegin() string {
	return "{\n"
}

func (g *GoLangBuilder) StringBlockEnd() string {
	return "}\n"
}
