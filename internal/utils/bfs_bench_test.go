package utils

import (
	"strconv"
	"testing"
)

// BenchmarkBreadthFirstSearchLinear benchmarks BFS on a linear graph (chain).
func BenchmarkBreadthFirstSearchLinear(b *testing.B) {
	for _, size := range []int{10, 100, 1000} {
		// Build adjacency: 0->1->2->...->size-1
		adj := make(map[int][]int, size)
		for i := range size - 1 {
			adj[i] = []int{i + 1}
		}

		b.Run("size_"+strconv.Itoa(size)+"_find_last", func(b *testing.B) {
			target := size - 1
			b.ReportAllocs()
			for b.Loop() {
				_ = BreadthFirstSearch(
					0,
					func(n int) []int { return adj[n] },
					func(n int) (bool, bool) {
						return n == target, n == target
					},
					BreadthFirstSearchOptions[int]{},
				)
			}
		})

		b.Run("size_"+strconv.Itoa(size)+"_find_first", func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_ = BreadthFirstSearch(
					0,
					func(n int) []int { return adj[n] },
					func(n int) (bool, bool) {
						return n == 0, true
					},
					BreadthFirstSearchOptions[int]{},
				)
			}
		})
	}
}

// BenchmarkBreadthFirstSearchTree benchmarks BFS on a binary tree graph.
func BenchmarkBreadthFirstSearchTree(b *testing.B) {
	for _, depth := range []int{5, 8, 10} {
		// Build binary tree: node i has children 2*i+1, 2*i+2
		size := (1 << depth) - 1 // 2^depth - 1 nodes
		adj := make(map[int][]int, size)
		for i := range size {
			left := 2*i + 1
			right := 2*i + 2
			var children []int
			if left < size {
				children = append(children, left)
			}
			if right < size {
				children = append(children, right)
			}
			if len(children) > 0 {
				adj[i] = children
			}
		}

		b.Run("depth_"+strconv.Itoa(depth)+"_find_leaf", func(b *testing.B) {
			target := size - 1
			b.ReportAllocs()
			for b.Loop() {
				_ = BreadthFirstSearch(
					0,
					func(n int) []int { return adj[n] },
					func(n int) (bool, bool) {
						return n == target, true
					},
					BreadthFirstSearchOptions[int]{},
				)
			}
		})

		b.Run("depth_"+strconv.Itoa(depth)+"_no_match", func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_ = BreadthFirstSearch(
					0,
					func(n int) []int { return adj[n] },
					func(n int) (bool, bool) {
						return false, false
					},
					BreadthFirstSearchOptions[int]{},
				)
			}
		})
	}
}

// BenchmarkBreadthFirstSearchWithCycles benchmarks BFS on a graph with cycles.
func BenchmarkBreadthFirstSearchWithCycles(b *testing.B) {
	// Build a dense graph where each node connects to the next 3 nodes (mod size)
	size := 100
	adj := make(map[int][]int, size)
	for i := range size {
		adj[i] = []int{(i + 1) % size, (i + 2) % size, (i + 3) % size}
	}

	b.Run("find_opposite", func(b *testing.B) {
		target := size / 2
		b.ReportAllocs()
		for b.Loop() {
			_ = BreadthFirstSearch(
				0,
				func(n int) []int { return adj[n] },
				func(n int) (bool, bool) {
					return n == target, true
				},
				BreadthFirstSearchOptions[int]{},
			)
		}
	})

	b.Run("exhaustive_no_match", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			_ = BreadthFirstSearch(
				0,
				func(n int) []int { return adj[n] },
				func(n int) (bool, bool) {
					return false, false
				},
				BreadthFirstSearchOptions[int]{},
			)
		}
	})
}
