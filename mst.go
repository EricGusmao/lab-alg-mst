package main

import (
	"cmp"
	"slices"
)

// Edge represents a connection between two nodes with a weight
type Edge struct {
	Source int
	Dest   int
	Weight int
}

// Graph basic structure
type Graph struct {
	V     int    // Number of Vertices
	Edges []Edge // List of Edges
}

// findIterative implements Find with path compression (iterative)
func findIterative(parent []int, i int) int {
	// Find the root
	root := i
	for parent[root] != root {
		root = parent[root]
	}

	// Path compression: point all nodes in the path directly to the root
	curr := i
	for curr != root {
		next := parent[curr]
		parent[curr] = root
		curr = next
	}

	return root
}

// KruskalMST executes the optimized Kruskal's algorithm
// Returns the MST edges and the total weight
func KruskalMST(graph Graph) ([]Edge, int) {
	// The maximum size of an MST is always V-1.
	treeSize := graph.V - 1
	// Handle edge case where V=0 or V=1, prevent negative capacity
	if treeSize < 0 {
		treeSize = 0
	}
	result := make([]Edge, 0, treeSize)
	totalWeight := 0

	// In-Place Sort
	slices.SortFunc(graph.Edges, func(a, b Edge) int {
		return cmp.Compare(a.Weight, b.Weight)
	})

	parent := make([]int, graph.V)
	rank := make([]int, graph.V)
	for i := range graph.V {
		parent[i] = i
	}

	edgesCount := 0

	// Iterate over sorted edges
	for _, edge := range graph.Edges {
		// If we have already formed the tree (V-1 edges), stop immediately.
		if edgesCount >= treeSize {
			break
		}

		rootX := findIterative(parent, edge.Source)
		rootY := findIterative(parent, edge.Dest)

		if rootX != rootY {
			result = append(result, edge)
			totalWeight += edge.Weight
			edgesCount++

			// Union by Rank (inlined logic for performance)
			if rank[rootX] < rank[rootY] {
				parent[rootX] = rootY
			} else if rank[rootX] > rank[rootY] {
				parent[rootY] = rootX
			} else {
				parent[rootY] = rootX
				rank[rootX]++
			}
		}
	}

	return result, totalWeight
}