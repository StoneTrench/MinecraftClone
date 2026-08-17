package ecs

import (
	"sync"
)

// TODO: System scheduling
// TODO: Spatial mapping/querying

type entityid_t uint16
type generation_t uint16 // for entity validation (because ids can get reassigned)

type EntityRef struct {
	id  entityid_t
	gen generation_t
}

type EntityWorld struct {
	mu         sync.Mutex
	OpenSet    []entityid_t
	Generation []generation_t
	NextId     entityid_t

	ComponentSets map[string]*ComponentSet
}

func (s *EntityWorld) CreateEntity() EntityRef {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.OpenSet) == 0 {
		res := EntityRef{
			id:  s.NextId,
			gen: 0,
		}
		s.NextId++
		s.Generation = append(s.Generation, 0)
		return res
	}

	id := s.OpenSet[len(s.OpenSet)-1]
	s.OpenSet = s.OpenSet[:(len(s.OpenSet) - 1)]
	gens := s.Generation[id]

	return EntityRef{
		id:  id,
		gen: gens,
	}
}

func (s *EntityWorld) DestroyEntity(e EntityRef) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.Generation[e.id] != e.gen {
		return // because the entity has already been destroyed
	}
	for _, c := range s.ComponentSets {
		c.RemoveComponent(e)
	}
	s.Generation[e.id] += 1
	s.OpenSet = append(s.OpenSet, e.id)
}

func (s *EntityWorld) CreateComponentSet(id string, size componentsize_t) {
	// TODO: Component set management
}
