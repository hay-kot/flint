// Package set implements a very basic set data structure. It is immutable and
// provides some helper methods specifically for working with the flint
package set

type ImmutableSet[T key] struct {
	s *Set[T]
}

func NewImmutable[T key](v ...T) *ImmutableSet[T] {
	return &ImmutableSet[T]{New(v...)}
}

func (s *ImmutableSet[T]) Contains(v T) bool {
	return s.s.Contains(v)
}
func (s *ImmutableSet[T]) ContainsAll(v ...T) bool {
	return s.s.ContainsAll(v...)
}

func (s *ImmutableSet[T]) Len() int {
	return s.s.Len()
}

func (s *ImmutableSet[T]) Intersection(o *Set[T]) *ImmutableSet[T] {
	return &ImmutableSet[T]{
		s: s.s.Intersection(o),
	}
}

func (s *ImmutableSet[T]) Missing(o *Set[T]) *ImmutableSet[T] {
	return &ImmutableSet[T]{
		s: s.s.Missing(o),
	}
}

func (s *ImmutableSet[T]) Slice() []T {
	return s.s.Slice()
}
