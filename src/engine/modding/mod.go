package modding

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/Masterminds/semver/v3"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
)

type ModHeader struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Author      string `json:"author"`
	Description string `json:"description"`

	Homepage string `json:"homepage"`

	ApiVersion        int               `json:"api_version"`
	Dependencies      map[string]string `json:"dependencies"`
	Incompatibilities map[string]string `json:"incompatibilities"`

	Environment map[string]string `json:"environment"`
}
type Mod struct {
	Folder  string
	Version semver.Version

	Dependencies map[string]semver.Constraints

	Header ModHeader

	Config       map[string]string
	WazeroModule api.Module

	ctx     context.Context
	Runtime wazero.Runtime
}

func (m *Mod) GetName() string {
	return m.Header.Name
}
func (m *Mod) GetVersion() string {
	return m.Version.String()
}
func (m *Mod) GetAuthor() string {
	return m.Header.Author
}
func (m *Mod) GetDescription() string {
	return m.Header.Description
}
func (m *Mod) GetHomepage() string {
	return m.Header.Homepage
}
func (m *Mod) IsApiVersionValid(api_version int) error {
	if m.Header.ApiVersion != api_version {
		return fmt.Errorf("API version of engine and mod do not match (%d != %d)", api_version, m.Header.ApiVersion)
	}
	return nil
}

func (m *Mod) Init(folder string) *Mod {
	m.Folder = folder
	return m
}
func (m *Mod) Close() {
	if m.WazeroModule != nil {
		m.WazeroModule.Close(m.ctx)
	}
}
func (m *Mod) LoadHeader() error {
	file_mod_json, err := os.ReadFile(path.Join(m.Folder, "mod.json"))
	if err != nil {
		return fmt.Errorf("failed to read mod.json: %w", err)
	}

	var head ModHeader
	err = json.Unmarshal(file_mod_json, &head)
	if err != nil {
		return fmt.Errorf("failed to parse mod.json: %w", err)
	}
	m.Header = head

	vers, err := semver.NewVersion(head.Version)
	if err != nil {
		return fmt.Errorf("failed to parse mod version: %w", err)
	}
	m.Version = *vers

	m.Dependencies = make(map[string]semver.Constraints)
	for k, v := range head.Dependencies {
		vers, err := semver.NewConstraint(v)
		if err != nil {
			return fmt.Errorf("failed to parse mod dependency (%s: %s): %w", k, v, err)
		}
		m.Dependencies[k] = *vers
	}

	return nil
}
func (m *Mod) LoadWasm() error {
	m.ctx = context.Background()
	m.Runtime = wazero.NewRuntime(m.ctx)

	file_path := path.Join(m.Folder, "main.wasm")
	bin, err := os.ReadFile(file_path)
	if err != nil {
		return fmt.Errorf("failed to open wasm file (%s): %w", m.String(), err)
	}

	comp_mod, err := m.Runtime.CompileModule(m.ctx, bin)
	if err != nil {
		return fmt.Errorf("failed to compile module (%s): %w", m.String(), err)
	}

	cnf := wazero.NewModuleConfig().
		WithStartFunctions("_initialize"). // this has to be called here manually for some reason
		WithName(m.GetName())

	err = load_api(m.Runtime, m)
	if err != nil {
		return fmt.Errorf("failed load api (%s): %w", m.String(), err)
	}

	inst_mod, err := m.Runtime.InstantiateModule(m.ctx, comp_mod, cnf)
	if err != nil {
		return fmt.Errorf("failed to instantiate module (%s): %w", m.String(), err)
	}

	m.WazeroModule = inst_mod
	return nil
}
func (m *Mod) String() string {
	return fmt.Sprintf("%s@%s", m.GetName(), m.GetVersion())
}
func (m *Mod) Equal(other *Mod) bool {
	if (m == nil) || (other == nil) {
		return false
	}
	return m.String() == other.String()
}
