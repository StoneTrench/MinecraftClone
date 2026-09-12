package main

import (
	"os"
	"strings"
)

// TODO: Finish binding generator, implement c bindings

const MODULE_NAME = "engine_bindings"

func hDeclType(guest, host *strings.Builder, state ILangBuilder, name string, t IType) TypeIdent {
	DeclType(guest, state, name, t)
	return DeclType(host, GOLANG, name, t)
}

func hDeclFunc(guest, host *strings.Builder, state ILangBuilder, name string, t TypeFunc) {
	DeclFunc(guest, state, MODULE_NAME, name, t)
	DeclHostFunc(host, name, t)
}

func main() {
	bguest := &strings.Builder{}
	bhost := &strings.Builder{}
	state := &GoLangBuilder{}

	bguest.WriteString("package api_bindings\n")
	bguest.WriteString("import (\n")
	bguest.WriteString("\"unsafe\"")
	bguest.WriteString("\n")
	bguest.WriteString(")\n")
	bhost.WriteString("package host\n")
	bhost.WriteString("import \"context\"\n")
	bhost.WriteString("import \"github.com/tetratelabs/wazero/api\"\n")

	Color_t := hDeclType(bguest, bhost, state, "Color", TypeStruct{Fields: []Param{
		{"R", TypeP_u8},
		{"G", TypeP_u8},
		{"B", TypeP_u8},
		{"A", TypeP_u8},
	}})
	BlockType_t := hDeclType(bguest, bhost, state, "BlockType", TypeStruct{Fields: []Param{
		{"Name", TypeP_string},
		{"Color", Color_t},
		{"IsSolid", TypeP_bool},
	}})

	fn_slog := TypeFunc{[]Param{{"msg", TypeP_string}}, nil}

	hDeclFunc(bguest, bhost, state, "log-info", fn_slog)
	hDeclFunc(bguest, bhost, state, "log-warn", fn_slog)
	hDeclFunc(bguest, bhost, state, "log-error", fn_slog)

	hDeclFunc(bguest, bhost, state, "register-block", TypeFunc{
		[]Param{{"block", BlockType_t}},
		[]Param{{"namespaceid", TypeP_string}},
	})
	hDeclFunc(bguest, bhost, state, "append-lang", TypeFunc{
		[]Param{{"name", TypeP_string}},
		nil,
	})

	hDeclFunc(bguest, bhost, state, "config-set-int", TypeFunc{Input: []Param{{"name", TypeP_string}, {"value", TypeP_i64}}})
	hDeclFunc(bguest, bhost, state, "config-set-float", TypeFunc{Input: []Param{{"name", TypeP_string}, {"value", TypeP_f64}}})
	hDeclFunc(bguest, bhost, state, "config-set-string", TypeFunc{Input: []Param{{"name", TypeP_string}, {"value", TypeP_string}}})

	hDeclFunc(bguest, bhost, state, "config-get-int", TypeFunc{
		[]Param{{"name", TypeP_string}},
		[]Param{{"value", TypeP_i64}},
	})
	hDeclFunc(bguest, bhost, state, "config-get-float", TypeFunc{
		[]Param{{"name", TypeP_string}},
		[]Param{{"value", TypeP_f64}},
	})
	hDeclFunc(bguest, bhost, state, "config-get-string", TypeFunc{
		[]Param{{"name", TypeP_string}},
		[]Param{{"value", TypeP_string}},
	})

	file, err := os.Create("bindings/go/api.go")
	if err != nil {
		panic(err)
	}
	_, err = file.WriteString(strings.Replace(bguest.String(), ";", ";\n", -1))
	if err != nil {
		panic(err)
	}
	file.Close()

	file, err = os.Create("bindings/host/api.go")
	if err != nil {
		panic(err)
	}
	_, err = file.WriteString(strings.Replace(bhost.String(), ";", ";\n", -1))
	if err != nil {
		panic(err)
	}
	file.Close()
}
