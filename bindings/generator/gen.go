package main

import (
	"fmt"
	"strings"
)

func DeclFunc(b *strings.Builder, state ILangBuilder, module, name string, fn TypeFunc) {
	bind_name := HyphenCaseToPascalCase(name)
	api_name := fmt.Sprintf("%s_%s", module, bind_name)

	params := fn.Input
	rets := fn.Output

	uses_malloc := false
	wasm_params := []Param{}
	wasm_input_vars := []Param{}
	wasm_output_vars := []Param{}

	for _, p := range rets {
		uses_malloc = uses_malloc || IsUseMalloc(p)
	}

	for _, p := range FullyLowerParams(params) {
		wasm_input_vars = append(wasm_input_vars, p)
		wasm_params = append(wasm_params, p)
	}

	for _, e := range FullyLowerParams(rets) {
		e.Type = TypeP_WasmPointer
		wasm_output_vars = append(wasm_output_vars, e)
		wasm_params = append(wasm_params, e)
	}

	Writes(
		b,
		state.StringWasmImport(
			module, name,
			api_name,
			TypeFunc{
				Input: wasm_params,
			},
		),
	)
	Writes(b, state.StringFnHeader(bind_name, fn))
	Writes(b, state.StringBlockBegin())
	{
		for _, p := range params {
			state.WriteGuestEncode(b, p)
		}
		if len(rets) > 0 {
			for _, p := range wasm_output_vars {
				Writes(b, state.StringDeclare(p))
			}
		}
		Writes(b, state.StringCall(api_name, wasm_input_vars, wasm_output_vars))
		if len(rets) > 0 {
			for _, p := range rets {
				state.WriteGuestDecode(b, p)
			}
		}
		if uses_malloc {
			Writes(b, state.StringResetMalloc())
		}
		if len(rets) > 0 {
			Writes(b, state.StringReturn(rets))
		}
	}
	Writes(b, state.StringBlockEnd())
}

func DeclType(b *strings.Builder, state ILangBuilder, name string, t IType) TypeIdent {
	Writef(b, "type %s %s\n", name, state.StringType(t))
	return TypeIdent{t, name}
}

func DeclHostFunc(b *strings.Builder, name string, fn TypeFunc) {
	bind_name := HyphenCaseToPascalCase(name)
	config := &ConfigHostWriter{
		Golang:        &GoLangBuilder{},
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
