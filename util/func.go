package util

import "time"

func ToPtr[T uint64 | uint32 | time.Time](v T) *T {
	return &v
}

func FromPtr[T uint64 | uint32 | time.Time](v *T) T {
	return *v
}
