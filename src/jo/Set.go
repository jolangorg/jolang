package jo

type Set[T comparable] map[T]struct{}

func NewSet[T comparable](elems []T) Set[T] {
	s := make(Set[T])
	s.AddMany(elems)
	return s
}

func (s Set[T]) Add(e T) {
	s[e] = struct{}{}
}

func (s Set[T]) AddMany(elems []T) {
	for _, e := range elems {
		s[e] = struct{}{}
	}
}

func (s Set[T]) Contains(e T) bool {
	_, ok := s[e]
	return ok
}

func (s Set[T]) Length(e T) int {
	return len(s)
}
