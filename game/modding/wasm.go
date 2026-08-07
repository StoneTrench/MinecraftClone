package modding

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/Masterminds/semver/v3"
	"github.com/StoneTrench/go-mc-clone/game/log"
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
	err := m.IsApiVersionValid()
	if err != nil {
		log.Warn(err)
	}

	log.Infof("Loading wasm module for (%s)", m.String())
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

func query_dependencies(mod_list []Mod, m Mod) ([]Mod, []error) {
	dep_count := len(m.Dependencies)
	res := make([]Mod, dep_count)

	err_list := []error{}

	dep_index := 0
	for k, v := range m.Dependencies {
		var found = false

		for _, dep := range mod_list {
			if dep.Header.Name == k {
				if !v.Check(&dep.Version) {
					err := fmt.Errorf(
						"wrong version of mod, (%s %s) expects (%s %s), got (%s %s) instead",
						m.Header.Name,
						m.Version.String(),
						k,
						v.String(),
						dep.Header.Name,
						dep.Version.String(),
					)

					err_list = append(err_list, err)
					continue
				}

				res[dep_index] = dep
				dep_index++
				found = true
			}
		}

		if !found {
			err := fmt.Errorf(
				"missing dependency, (%s %s) requires dependency (%s %s)",
				m.Header.Name,
				m.Version.String(),
				k,
				v.String(),
			)
			err_list = append(err_list, err)
			continue
		}
	}

	return res, err_list
}
func sort_by_dependency(mod_list []Mod) ([]Mod, error) {
	sorted := []Mod{}
	visited := make(map[string]struct{})

	err_list := []error{}

	for _, v := range mod_list {
		err := sort_by_dependency_view(mod_list, v, visited, &sorted)
		if err != nil {
			err_list = append(err_list, err...)
		}
	}

	return sorted, errors.Join(err_list...)
}
func sort_by_dependency_view(mod_list []Mod, mod Mod, visited map[string]struct{}, sorted *[]Mod) []error {
	if _, exists := visited[mod.Header.Name]; exists {
		not_contains := true
		for _, m := range *sorted {
			if mod.Equal(&m) {
				not_contains = false
				break
			}
		}
		if not_contains {
			return []error{fmt.Errorf("cyclic dependency found in relation to mod (%s)", mod.String())}
		}
	}

	err_list := []error{}

	deps, err := query_dependencies(mod_list, mod)
	if err != nil {
		err_list = append(err_list, err...)
	}

	for _, e := range deps {
		err := sort_by_dependency_view(mod_list, e, visited, sorted)
		if err != nil {
			err_list = append(err_list, err...)
		}
	}

	*sorted = append((*sorted), mod)

	return err_list
}

func LoadMods(ctx context.Context) ([]Mod, error) {
	entries, err := os.ReadDir("./mods/")
	if err != nil {
		return nil, fmt.Errorf("failed to load mods, %w", err)
	}

	// Load mod headers
	mods := []Mod{}

	for _, e := range entries {
		if e.IsDir() {
			mod := Mod{}
			mod.Init(path.Join("./mods", e.Name()))

			err := mod.LoadHeader()
			if err != nil {
				log.Error(fmt.Errorf("failed to load a mod header, %w", err))
			} else {
				mods = append(mods, mod)
			}
		}
	}

	// Sort mods based on dependency
	// Skip mods with missing dependencies and log errors
	mods, err = sort_by_dependency(mods)
	if err != nil {
		return nil, err
	}

	// Load guest modules
	for _, m := range mods {
		err := m.LoadGuestModule(ctx)
		if err != nil {
			return nil, err
		}
	}

	// Finish

	return mods, nil
}

func Init() error {
	log.Infof("Init wasm (api %d)", API_VERSION)
	ctx := context.Background()

	_, err := LoadMods(ctx)
	if err != nil {
		return err
	}
	// for _, m := range mods {
	// 	defer m.Close(ctx)
	// }

	return nil
}
