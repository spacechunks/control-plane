package ptr

// Deref dereferences a pointer or returns the types zero value, if it's null.
func Deref[T any](p *T) T {
	if p == nil {
		var zero T
		return zero
	}
	return *p
}
