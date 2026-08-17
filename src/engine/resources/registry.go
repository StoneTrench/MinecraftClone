package resources

import (
	"fmt"
	"iter"
	"log/slog"
	"sync"

	"golang.org/x/exp/constraints"
)

/*

There's two phases of the registry's lifespan:
+ Registering resources
+ Querying for resources, updating ids, saving id mappings
The numerical ids of resources, can be assigned after the registration phase is over.

Q: What happens when a new mod is added?
A: Id mappings that already exist, get correctly assigned, then the registry handles the new resources, by just assigning new, unique Ids.

Q: What happens when a mod is removed?
A: In the save file we store the mapping of resources, and a list of mods used. The game will issue a warning of missing mods.

Q: What happens when a mod updates names?
A: Don't.

Q: What happens when mods are loaded in a different order?
A: Order doesn't matter, because resources can be loaded in any order, as their numerical ids are handled after registration.

Q: How specifically does id mapping work?
A:

*/

type Registry[T any, I constraints.Integer] struct {
	mu sync.RWMutex

	name string

	resources map[string]T

	is_frozen   bool
	is_id_valid bool // when you leave a game, ids get invalidated, when you join a game, the ids get loaded

	version uint32 // registry_version is incremented every time the the registry gets invalidated, and once when it is frozen

	id2ident_arr []string
	ident2id_map map[string]I
}

func CreateRegistry[T any, I constraints.Integer](capacity int, name string) *Registry[T, I] {
	return &Registry[T, I]{
		name:         name,
		resources:    make(map[string]T),
		is_frozen:    false,
		is_id_valid:  false,
		version:      0,
		id2ident_arr: make([]string, 0, capacity),
		ident2id_map: make(map[string]I, capacity),
	}
}

func (r *Registry[T, I]) GetVersion() uint32 {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.version
}
func (r *Registry[T, I]) errorf(format string, args ...any) error {
	return fmt.Errorf("in registry (name:\"%s\",frozen:%v,valid:%v,version:%d): %w", r.name, r.is_frozen, r.is_id_valid, r.version, fmt.Errorf(format, args...))
}
func (r *Registry[T, I]) err_is_frozen() error {
	if r.is_frozen {
		return fmt.Errorf("registry is frozen: registration phase already ended")
	}
	return nil
}
func (r *Registry[T, I]) err_not_is_frozen() error {
	if !r.is_frozen {
		return fmt.Errorf("registry is not frozen: registration phase needs to end before registry can be queried")
	}
	return nil
}
func (r *Registry[T, I]) err_is_id_valid() error {
	if r.is_frozen {
		return fmt.Errorf("registry ids have not yet been invalidated: this means that the assigned numerical ids of resources still match their namespace identifiers, invalidate first")
	}
	return nil
}
func (r *Registry[T, I]) err_not_is_id_valid() error {
	if !r.is_id_valid {
		return fmt.Errorf("registry ids have been invalidated: this means that the assigned numerical ids of resources will no longer match their namespace identifiers, waiting for new id map")
	}
	return nil
}

func (r *Registry[T, I]) Register(mod string, name string, resource T) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if err := r.err_is_frozen(); err != nil {
		return "", r.errorf("failed to register resource: %w", err)
	}
	ident := fmt.Sprintf("%s:%s", mod, name)
	if _, exists := r.resources[ident]; exists {
		return "", r.errorf("failed to register resource: resource identifier already used")
	}
	r.resources[ident] = resource
	return ident, nil
}
func (r *Registry[T, I]) RegisterPanic(mod string, name string, resource T) string {
	n, err := r.Register(mod, name, resource)
	if err != nil {
		panic(err)
	}
	return n
}
func (r *Registry[T, I]) Freeze() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.is_frozen = true
	r.version++
}

func (r *Registry[T, I]) InvalidateIds() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if err := r.err_not_is_frozen(); err != nil {
		return r.errorf("failed to invalidate id map: %w", err)
	}
	r.is_id_valid = false
	r.version++
	return nil
}
func (r *Registry[T, I]) ApplyIdMap(m []string) (new_map []string, err error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if err = r.err_not_is_frozen(); err != nil {
		return nil, r.errorf("failed to handle id map: %w", err)
	}
	m_len := len(m)
	r.id2ident_arr = m
	r.ident2id_map = make(map[string]I)

	missing := make([]string, 0)
	loaded_len := 0

	for id, ident := range m {
		r.ident2id_map[ident] = I(id)

		_, exists := r.resources[ident]
		if !exists {
			missing = append(missing, ident)
		} else {
			loaded_len++
		}
	}

	extra_len := len(r.resources) - loaded_len

	// log.Assert(extra_len >= 0, "registry: extra_len must always be >= 0")

	if extra_len > 0 {
		for ident := range r.resources {
			_, exists := r.ident2id_map[ident]

			if !exists {
				r.ident2id_map[ident] = I(m_len)
				r.id2ident_arr = append(r.id2ident_arr, ident)
				m_len++
			}
		}
	}

	if len(missing) > 0 {
		slog.Warn(fmt.Sprintf("missing registry entries when loading id map: %v", missing))
	}

	// log.Assert(m_len == len(r.id2ident_arr), "registry: m_len and id2ident_arr should be equal")
	// log.Assert(len(r.id2ident_arr) == len(r.ident2id_map), "registry: id2ident_arr and ident2id_map should be equal")

	r.is_id_valid = true

	return r.id2ident_arr, nil
}

func (r *Registry[T, I]) GetByIdentifier(ident string) (T, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if err := r.err_not_is_frozen(); err != nil {
		var NA T
		return NA, r.errorf("failed to get resource by identifier: %w", err)
	}
	res, exists := r.resources[ident]
	if !exists {
		return res, r.errorf("resource with identifier (%s) not registered", ident)
	}
	return res, nil
}
func (r *Registry[T, I]) GetId(ident string) (I, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var NA I
	if err := r.err_not_is_frozen(); err != nil {
		return NA, r.errorf("failed to get resource id: %w", err)
	}
	if err := r.err_not_is_id_valid(); err != nil {
		return NA, r.errorf("failed to get resource by id: %w", err)
	}
	res, exists := r.ident2id_map[ident]
	if !exists {
		return res, r.errorf("resource with identifier (%s) not registered", ident)
	}
	return res, nil
}
func (r *Registry[T, I]) GetIdentifier(id I) (string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var NA string
	if err := r.err_not_is_frozen(); err != nil {
		return NA, r.errorf("failed to get resource identifier: %w", err)
	}
	if err := r.err_not_is_id_valid(); err != nil {
		return NA, r.errorf("failed to get resource by id: %w", err)
	}
	if 0 > id || id >= I(len(r.id2ident_arr)) {
		return NA, r.errorf("resource with id (%d) not registered", id)
	}
	return r.id2ident_arr[id], nil
}
func (r *Registry[T, I]) GetById(id I) (T, error) {
	var NA T
	ident, err := r.GetIdentifier(id)
	if err != nil {
		return NA, err
	}
	return r.GetByIdentifier(ident)
}

func (r *Registry[T, I]) Iterate() iter.Seq2[string, T] {
	return func(yield func(string, T) bool) {
		r.mu.RLock()
		defer r.mu.RUnlock()
		for k, v := range r.resources {
			if !yield(k, v) {
				return
			}
		}
	}
}
