package utils

import (
	"strconv"
	"testing"
)

// BenchmarkFilter benchmarks Filter with various sizes.
func BenchmarkFilter(b *testing.B) {
	isEven := func(n int) bool { return n%2 == 0 }

	for _, size := range []int{10, 100, 1000} {
		slice := make([]int, size)
		for i := range slice {
			slice[i] = i
		}

		b.Run("size_"+strconv.Itoa(size)+"_half_match", func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_ = Filter(slice, isEven)
			}
		})

		b.Run("size_"+strconv.Itoa(size)+"_all_match", func(b *testing.B) {
			b.ReportAllocs()
			alwaysTrue := func(int) bool { return true }
			for b.Loop() {
				_ = Filter(slice, alwaysTrue)
			}
		})

		b.Run("size_"+strconv.Itoa(size)+"_none_match", func(b *testing.B) {
			b.ReportAllocs()
			alwaysFalse := func(int) bool { return false }
			for b.Loop() {
				_ = Filter(slice, alwaysFalse)
			}
		})
	}
}

// BenchmarkFilterIndex benchmarks FilterIndex.
func BenchmarkFilterIndex(b *testing.B) {
	slice := make([]int, 100)
	for i := range slice {
		slice[i] = i
	}

	b.Run("half_match", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			_ = FilterIndex(slice, func(_ int, i int, _ []int) bool { return i%2 == 0 })
		}
	})

	b.Run("all_match", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			_ = FilterIndex(slice, func(_ int, _ int, _ []int) bool { return true })
		}
	})
}

// BenchmarkMap benchmarks Map with various sizes.
func BenchmarkMap(b *testing.B) {
	for _, size := range []int{10, 100, 1000} {
		slice := make([]int, size)
		for i := range slice {
			slice[i] = i
		}

		b.Run("int_to_int_"+strconv.Itoa(size), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_ = Map(slice, func(n int) int { return n * 2 })
			}
		})

		b.Run("int_to_string_"+strconv.Itoa(size), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_ = Map(slice, func(n int) string { return strconv.Itoa(n) })
			}
		})
	}
}

// BenchmarkSome benchmarks Some with early exit and full scan.
func BenchmarkSome(b *testing.B) {
	slice := make([]int, 1000)
	for i := range slice {
		slice[i] = i
	}

	b.Run("early_exit", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			_ = Some(slice, func(n int) bool { return n == 5 })
		}
	})

	b.Run("full_scan_no_match", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			_ = Some(slice, func(n int) bool { return n < 0 })
		}
	})
}

// BenchmarkEvery benchmarks Every with early exit and full scan.
func BenchmarkEvery(b *testing.B) {
	slice := make([]int, 1000)
	for i := range slice {
		slice[i] = i
	}

	b.Run("early_exit", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			_ = Every(slice, func(n int) bool { return n > 500 })
		}
	})

	b.Run("full_scan_all_match", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			_ = Every(slice, func(n int) bool { return n >= 0 })
		}
	})
}

// BenchmarkFlatten benchmarks Flatten with various inner sizes.
func BenchmarkFlatten(b *testing.B) {
	b.Run("10x10", func(b *testing.B) {
		data := make([][]int, 10)
		for i := range data {
			data[i] = make([]int, 10)
			for j := range data[i] {
				data[i][j] = i*10 + j
			}
		}
		b.ReportAllocs()
		for b.Loop() {
			_ = Flatten(data)
		}
	})

	b.Run("100x10", func(b *testing.B) {
		data := make([][]int, 100)
		for i := range data {
			data[i] = make([]int, 10)
			for j := range data[i] {
				data[i][j] = i*10 + j
			}
		}
		b.ReportAllocs()
		for b.Loop() {
			_ = Flatten(data)
		}
	})

	b.Run("10x100", func(b *testing.B) {
		data := make([][]int, 10)
		for i := range data {
			data[i] = make([]int, 100)
			for j := range data[i] {
				data[i][j] = i*100 + j
			}
		}
		b.ReportAllocs()
		for b.Loop() {
			_ = Flatten(data)
		}
	})
}

// BenchmarkIsStringWhiteSpace benchmarks whitespace checking on various inputs.
func BenchmarkIsStringWhiteSpace(b *testing.B) {
	b.Run("empty", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			_ = IsStringWhiteSpace("")
		}
	})

	b.Run("whitespace_short", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			_ = IsStringWhiteSpace("   \t\n  ")
		}
	})

	b.Run("non_whitespace_early", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			_ = IsStringWhiteSpace("hello world")
		}
	})

	b.Run("mixed_long", func(b *testing.B) {
		s := "                                                            x"
		b.ReportAllocs()
		for b.Loop() {
			_ = IsStringWhiteSpace(s)
		}
	})
}
