package modding

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path"
)

func load_mod_headers(entries []os.DirEntry) ([]Mod, error) {
	mods := []Mod{}
	err_list := []error{}
	for _, e := range entries {
		if e.IsDir() {
			mod := Mod{}
			mod.Init(path.Join("./mods", e.Name()))

			err := mod.LoadHeader()
			if err != nil {
				err_list = append(err_list, fmt.Errorf("failed to load a mod header, %w", err))
			} else {
				mods = append(mods, mod)
			}
		}
	}
	return mods, errors.Join(err_list...)
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

func LoadMods(ctx context.Context) ([]Mod, error) {
	entries, err := os.ReadDir("./mods/")
	if err != nil {
		return nil, fmt.Errorf("failed to load mods, %w", err)
	}

	// Load mod headers
	mods, err := load_mod_headers(entries)
	if err != nil {
		return nil, err
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
