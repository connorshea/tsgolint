package utils

import (
	"strconv"
	"testing"
)

// BenchmarkSetAdd benchmarks adding elements to a Set.
func BenchmarkSetAdd(b *testing.B) {
	for _, size := range []int{10, 100, 1000} {
		b.Run("size_"+strconv.Itoa(size), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				s := NewSetWithSizeHint[int](size)
				for i := range size {
					s.Add(i)
				}
			}
		})
	}
}

// BenchmarkSetHas benchmarks looking up elements in a Set.
func BenchmarkSetHas(b *testing.B) {
	for _, size := range []int{10, 100, 1000} {
		s := NewSetWithSizeHint[int](size)
		for i := range size {
			s.Add(i)
		}

		b.Run("size_"+strconv.Itoa(size)+"_hit", func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				for i := range size {
					_ = s.Has(i)
				}
			}
		})

		b.Run("size_"+strconv.Itoa(size)+"_miss", func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				for i := size; i < size*2; i++ {
					_ = s.Has(i)
				}
			}
		})
	}
}

// BenchmarkSetDelete benchmarks deleting elements from a Set.
func BenchmarkSetDelete(b *testing.B) {
	for _, size := range []int{10, 100, 1000} {
		b.Run("size_"+strconv.Itoa(size), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				s := NewSetWithSizeHint[int](size)
				for i := range size {
					s.Add(i)
				}
				for i := range size {
					s.Delete(i)
				}
			}
		})
	}
}

// BenchmarkNewSetFromItems benchmarks creating a Set from items.
func BenchmarkNewSetFromItems(b *testing.B) {
	for _, size := range []int{10, 100, 1000} {
		items := make([]int, size)
		for i := range items {
			items[i] = i
		}

		b.Run("size_"+strconv.Itoa(size), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_ = NewSetFromItems(items...)
			}
		})
	}
}

// BenchmarkSetStringKeys benchmarks Set with string keys to measure hashing overhead.
func BenchmarkSetStringKeys(b *testing.B) {
	for _, size := range []int{10, 100, 1000} {
		keys := make([]string, size)
		for i := range keys {
			keys[i] = "key_" + strconv.Itoa(i) + "_some_longer_string_to_measure_hash"
		}

		b.Run("add_"+strconv.Itoa(size), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				s := NewSetWithSizeHint[string](size)
				for _, k := range keys {
					s.Add(k)
				}
			}
		})

		s := NewSetWithSizeHint[string](size)
		for _, k := range keys {
			s.Add(k)
		}

		b.Run("has_"+strconv.Itoa(size), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				for _, k := range keys {
					_ = s.Has(k)
				}
			}
		})
	}
}
