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

	Config  map[string]string
	Runtime wazero.Runtime
	Module  api.Module
}

func (m *Mod) IsApiVersionValid() error {
	if m.Header.ApiVersion < API_VERSION {
		return fmt.Errorf("API version of engine and mod do not match (%d != %d)", API_VERSION, m.Header.ApiVersion)
	}
	return nil
}

func (m *Mod) Init(folder string) *Mod {
	m.Folder = folder
	return m
}
func (m *Mod) Close(ctx context.Context) {
	if m.Module != nil {
		m.Module.Close(ctx)
	}
	if m.Runtime != nil {
		m.Runtime.Close(ctx)
	}
}
func (m *Mod) LoadHeader() error {
	file_mod_json, err := os.ReadFile(path.Join(m.Folder, "mod.json"))
	if err != nil {
		return fmt.Errorf("failed to read mod.json, %w", err)
	}

	var head ModHeader
	err = json.Unmarshal(file_mod_json, &head)
	if err != nil {
		return fmt.Errorf("failed to parse mod.json, %w", err)
	}
	m.Header = head

	vers, err := semver.NewVersion(head.Version)
	if err != nil {
		return fmt.Errorf("failed to parse mod version, %w", err)
	}
	m.Version = *vers

	m.Dependencies = make(map[string]semver.Constraints)
	for k, v := range head.Dependencies {
		vers, err := semver.NewConstraint(v)
		if err != nil {
			return fmt.Errorf("failed to parse mod dependency (%s: %s), %w", k, v, err)
		}
		m.Dependencies[k] = *vers
	}

	return nil
}
func (m *Mod) LoadGuestModule(ctx context.Context) error {
	file_main_wasm, err := os.ReadFile(path.Join(m.Folder, "main.wasm"))
	if err != nil {
		return fmt.Errorf("failed to load (%s), %w", m.String(), err)
	}

	r := wazero.NewRuntime(ctx)
	m.Runtime = r

	// Load api
	err = load_api(ctx, r, m)
	if err != nil {
		return fmt.Errorf("failed to load (%s), %w", m.String(), err)
	}

	// Instantiate module
	config := wazero.
		NewModuleConfig().
		WithName(m.Header.Name).
		WithEnv("MOD_NAME", m.Header.Name).
		WithEnv("MOD_VERSION", m.Header.Version).
		WithStartFunctions("_initialize")

	for k, v := range m.Header.Environment {
		config = config.WithEnv(fmt.Sprintf("MOD_ENV_%s", k), v)
	}

	mod, err := r.InstantiateWithConfig(ctx, file_main_wasm, config)
	if err != nil {
		return fmt.Errorf("failed to instantiate with config (%s), %w", m.String(), err)
	}

	m.Module = mod
	return nil
}
func (m *Mod) String() string {
	return fmt.Sprintf("%s %s", m.Header.Name, m.Version.String())
}
func (m *Mod) Equal(other *Mod) bool {
	if (m == nil && other == nil) || (m == nil) || (other == nil) {
		return false
	}
	return m.Header.Name == other.Header.Name && m.Version.Equal(&other.Version)
}
