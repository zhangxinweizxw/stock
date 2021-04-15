package structures

import (
    "strconv"
	"strings"
	"sync"
)

type Set struct {
	m map[int64]bool
	sync.RWMutex
}

func New() *Set {
	return &Set{
		m: map[int64]bool{},
	}
}

func (s *Set) Add(item int64) {
	s.Lock()
	defer s.Unlock()
	s.m[item] = true
}

func (s *Set) Remove(item int64) {
	s.Lock()
	s.Unlock()
	delete(s.m, item)
}

func (s *Set) Has(item int64) bool {
	s.RLock()
	defer s.RUnlock()
	_, ok := s.m[item]
	return ok
}

func (s *Set) Len() int {
	return len(s.List())
}

func (s *Set) Clear() {
	s.Lock()
	defer s.Unlock()
	s.m = map[int64]bool{}
}

func (s *Set) IsEmpty() bool {
	if s.Len() == 0 {
		return true
	}
	return false
}

func (s *Set) String() string {
	s.RLock()
	defer s.RUnlock()
	list := []string{}
	for item := range s.m {
		list = append(list, strconv.FormatInt(item, 10))
	}
	return strings.Join(list, ",")
}

func (s *Set) List() []int64 {
	s.RLock()
	defer s.RUnlock()
	list := []int64{}
	for item := range s.m {
		list = append(list, item)
	}
	return list
}
