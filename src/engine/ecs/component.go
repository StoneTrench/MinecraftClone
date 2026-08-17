package ecs

import (
	"fmt"
	"sync"
)

type componentid_t uint16 // for individual components in a component set, not component types
type componentsize_t uint16

type ComponentSet struct {
	mu                sync.RWMutex
	ComponentSize     componentsize_t
	EntityToComponent []componentid_t // 0 is a sentinel value, meaning there is no component attached
	Components        []byte
}

func (s *ComponentSet) TryResize(e EntityRef) {
	s.mu.Lock()
	defer s.mu.Unlock()
	n := int(e.id) - len(s.EntityToComponent) + 1
	if n > 0 {
		s.EntityToComponent = append(s.EntityToComponent, make([]componentid_t, n)...)
	}
}

func (s *ComponentSet) SetComponent(e EntityRef, c []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(c) != int(s.ComponentSize) {
		return fmt.Errorf("invalid component size, expected %d got %d bytes", s.ComponentSize, len(c))
	}
	component_id := s.EntityToComponent[e.id]
	if component_id == 0 {
		s.Components = append(s.Components, c...)
		s.EntityToComponent[e.id] = componentid_t(len(s.Components) / int(s.ComponentSize))
	} else {
		component_ptr := uint32(component_id-1) * uint32(s.ComponentSize)
		copy(s.Components[component_ptr:(component_ptr+uint32(s.ComponentSize))], c)
	}
	return nil
}

func (s *ComponentSet) RemoveComponent(e EntityRef) {
	s.mu.Lock()
	defer s.mu.Unlock()
	remove_c_id := s.EntityToComponent[e.id]                             // get the component id
	last_c_id := componentid_t(len(s.Components) / int(s.ComponentSize)) // get the last component's id

	if last_c_id == 0 || remove_c_id == 0 {
		return // no components
	}

	remove_c_ptr := uint32(remove_c_id-1) * uint32(s.ComponentSize)
	last_c_ptr := uint32(last_c_id-1) * uint32(s.ComponentSize)

	s.EntityToComponent[e.id] = 0 // set the component id to the sentinel value

	// move the last component to the target component
	if last_c_id != remove_c_id {
		copy(s.Components[remove_c_ptr:(remove_c_ptr+uint32(s.ComponentSize))], s.Components[last_c_ptr:(last_c_ptr+uint32(s.ComponentSize))])
		for c, c_id := range s.EntityToComponent {
			if c_id == last_c_id {
				s.EntityToComponent[c] = remove_c_id // update the last component's reference index to the removed value, because it was moved there
				break
			}
		}
	}
	s.Components = s.Components[:last_c_ptr] // shrink the slice by 1
}
