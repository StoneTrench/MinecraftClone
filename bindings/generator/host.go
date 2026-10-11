package main

import (
	"fmt"
	"strings"
)

func DeclHostFunc(b *strings.Builder, name string, fn TypeFunc) {
	bind_name := HyphenCaseToPascalCase(name)
	config := &ConfigHostWriter{
		Golang:        GOLANG,
		NameApiModule: "mod",
		NameFnAlloc:   "alloc",
	}
	const fn_name = "fn"
	state := config.Golang
	uses_malloc := false

	wasm_params := []Param{
		{"ctx", TypeIdent{nil, "context.Context"}},
		{config.NameApiModule, TypeIdent{nil, "api.Module"}},
	}

	params := fn.Input
	rets := fn.Output

	for _, p := range rets {
		uses_malloc = uses_malloc || IsUseMalloc(p)
	}

	wasm_params = append(wasm_params, FullyLowerParams(params)...)
	for _, p := range FullyLowerParams(rets) {
		p.Type = TypeP_WasmPointer
		wasm_params = append(wasm_params, p)
	}

	Writes(b, state.StringFnHeader(bind_name, TypeFunc{
		Input: []Param{
			{fn_name, TypeFunc{
				Input:  append([]Param{{"m", TypeIdent{nil, "api.Module"}}}, params...),
				Output: rets,
			}},
		},
		Output: []Param{
			{"w", TypeIdent{nil, "HostBinding"}},
		},
	}))
	Writes(b, state.StringBlockBegin())
	{
		Writef(b, "return HostBinding{\"%s\",", name)
		Writes(b, state.StringType(TypeFunc{Input: wasm_params}))
		Writes(b, state.StringBlockBegin())
		{
			if uses_malloc {
				Writef(b, "_malloc:=%s.ExportedFunction(\"malloc\");", config.NameApiModule)
				Writef(b, "%s:=func(size uint32)uint32{res,err:=_malloc.Call(ctx,uint64(size));if err!=nil{panic(err)};return uint32(res[0])};", config.NameFnAlloc)
			}
			for _, p := range params {
				WriteGoHostDecode(config, b, p)
			}
			if len(rets) > 0 {
				for i, p := range rets {
					Writes(b, p.Name)
					if i < len(rets)-1 {
						Writes(b, ",")
					}
				}
				Writes(b, ":=")
			}
			Writef(b, "%s(%s,", fn_name, config.NameApiModule)
			for i, p := range params {
				Writes(b, p.Name)
				if i < len(params)-1 {
					Writes(b, ",")
				}
			}
			Writes(b, ");")
			if len(rets) > 0 {
				for _, p := range rets {
					WriteGoHostEncode(config, b, p)
				}
			}
		}
		Writes(b, "}}\n")
	}
	Writes(b, state.StringBlockEnd())
}
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
