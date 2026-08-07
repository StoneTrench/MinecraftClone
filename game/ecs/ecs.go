package ecs

import (
	"fmt"
	"reflect"
	"sync"
)

type EntityId uint64
type empty struct{}

type ECS struct {
	mu         sync.RWMutex
	next_id    EntityId
	entities   map[EntityId]empty
	components map[reflect.Type]any
}

func CreateSystem() *ECS {
	return &ECS{
		entities:   make(map[EntityId]empty),
		components: make(map[reflect.Type]any),
	}
}

func (s *ECS) CreateEntity() EntityId {
	s.mu.Lock()
	defer s.mu.Unlock()
	id := s.next_id
	s.next_id++
	s.entities[id] = empty{}
	return id
}

func (s *ECS) DestroyEntity(id EntityId) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.entities, id)
	for _, store := range s.components {
		reflect.ValueOf(store).MethodByName("Delete").Call([]reflect.Value{reflect.ValueOf(id)})
	}
}

func AddComponent[T any](s *ECS, id EntityId, c T) {
	s.mu.Lock()
	defer s.mu.Unlock()
	t := reflect.TypeOf(c)
	store, ok := s.components[t]
	if !ok {
		store = createComponentStore[T]()
		s.components[t] = store
	}
	store.(*componentStore[T]).Data[id] = c
}

func GetComponent[T any](s *ECS, id EntityId) (*T, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	t := reflect.TypeFor[T]()
	store, ok := s.components[t]
	if !ok {
		return nil, fmt.Errorf("could not get component from entity %v, no store found with type %v", id, t)
	}
	res, ok := store.(*componentStore[T]).Data[id]
	if !ok {
		return nil, fmt.Errorf("could not get component from entity %v, no component found with type %v", id, t)
	}
	return &res, nil
}

func RemoveComponent[T any](s *ECS, id EntityId) {
	s.mu.Lock()
	defer s.mu.Unlock()
	t := reflect.TypeFor[T]()
	store, ok := s.components[t]
	if !ok {
		return
	}
	delete(store.(*componentStore[T]).Data, id)
}

func Query2[T1 any, T2 any](s *ECS) (entities []EntityId, comp1 []T1, comp2 []T2) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	t1 := reflect.TypeFor[T1]()
	t2 := reflect.TypeFor[T2]()

	store1, ok1 := s.components[t1]
	store2, ok2 := s.components[t2]
	if !ok1 || !ok2 {
		return
	}

	map1 := store1.(*componentStore[T1]).Data
	map2 := store2.(*componentStore[T2]).Data

	// Loop for the least amount of iterations
	if len(map1) < len(map2) {
		for id, c1 := range map1 {
			if c2, ok := map2[id]; ok {
				entities = append(entities, id)
				comp1 = append(comp1, c1)
				comp2 = append(comp2, c2)
			}
		}
	} else {
		for id, c2 := range map2 {
			if c1, ok := map1[id]; ok {
				entities = append(entities, id)
				comp1 = append(comp1, c1)
				comp2 = append(comp2, c2)
			}
		}
	}
	return
}

// Has to tightly pack components.
//
// Has to be able to quickly iterate those components.
//
// Has to allow queries of components based on EntityId.
//
// Has to be able to map EntityId to a component.
type componentStore[T any] struct {
	Data map[EntityId]T // TODO: This is a shitty implementation, fix it at some point please. [https://github.com/SanderMertens/ecs-faq]
}

func createComponentStore[T any]() *componentStore[T] {
	return &componentStore[T]{
		Data: make(map[EntityId]T),
	}
}

func (c *componentStore[T]) Delete(id EntityId) {
	delete(c.Data, id)
}
