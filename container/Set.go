package container

type Set map[any]struct{} // Define a Set type

// Add an element to the set
func (s Set) Add(value any) {
	s[value] = struct{}{} // struct{}{} takes up no memory
}

// Remove an element from the set
func (s Set) Remove(value any) {
	delete(s, value)
}

// Check if the set contains an element
func (s Set) Has(value any) bool {
	_, exists := s[value]
	return exists
}

// Get all elements in the set
func (s Set) Values() []any {
	keys := make([]any, 0, len(s))
	for key := range s {
		keys = append(keys, key)
	}
	return keys
}
