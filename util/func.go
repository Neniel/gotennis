package util

func ToPtr[T uint64 | uint32](v T) *T {
	return &v
}

func FromPtr[T uint64 | uint32](v *T) T {
	return *v
}
