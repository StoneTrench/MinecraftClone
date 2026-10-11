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
