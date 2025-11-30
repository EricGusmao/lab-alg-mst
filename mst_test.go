package main

import (
	"math/rand/v2" 
	"slices"
	"testing"
)

// Unit Test: Kept small, checks basic correctness.
func TestKruskalSmall(t *testing.T) {
	graph := Graph{
		V: 4,
		Edges: []Edge{
			{0, 1, 10},
			{0, 2, 6},
			{0, 3, 5},
			{1, 3, 15},
			{2, 3, 4},
		},
	}

	expected := 19
	_, cost := KruskalMST(graph)

	if cost != expected {
		t.Errorf("Incorrect result: expected %d, got %d", expected, cost)
	}
}

func TestKruskalMST_Scenarios(t *testing.T) {
	tests := []struct {
		name           string
		graph          Graph
		expectedWeight int
		expectedEdges  int // How many edges should be in the final solution
	}{
		{
			name: "Basic - Triangle",
			// 0 --(1)-- 1
			// |       /
			//(3)    (2)
			// |   /
			// 2
			graph: Graph{
				V: 3,
				Edges: []Edge{
					{0, 1, 1},
					{1, 2, 2},
					{0, 2, 3},
				},
			},
			expectedWeight: 3, // 1 + 2
			expectedEdges:  2, // V-1
		},
		{
			name: "Parallel Edges (Multigraph)",
			// Should choose the edge with the lowest weight between the same nodes
			graph: Graph{
				V: 2,
				Edges: []Edge{
					{0, 1, 100},
					{0, 1, 10}, // This one should be chosen
					{0, 1, 50},
				},
			},
			expectedWeight: 10,
			expectedEdges:  1,
		},
		{
			name: "Disconnected Graph (Forest)",
			// Two separate islands: 0-1 and 2-3.
			// The algorithm should return the sum of MSTs of each component.
			graph: Graph{
				V: 4,
				Edges: []Edge{
					{0, 1, 5},
					{2, 3, 10},
				},
			},
			expectedWeight: 15,
			expectedEdges:  2, // Cannot reach V-1 (3)
		},
		{
			name: "Linear Graph (Linked List)",
			// 0-1-2-3-4
			graph: Graph{
				V: 5,
				Edges: []Edge{
					{0, 1, 1},
					{1, 2, 2},
					{2, 3, 3},
					{3, 4, 4},
				},
			},
			expectedWeight: 10,
			expectedEdges:  4,
		},
		{
			name: "Lonely Node (V=1)",
			// A single node without edges has 0 weight and 0 edges in MST
			graph: Graph{
				V:     1,
				Edges: []Edge{},
			},
			expectedWeight: 0,
			expectedEdges:  0,
		},
		{
			name: "Cycle with Equal Weights",
			// 0-1(10), 1-2(10), 2-0(10). Should remove any one edge.
			graph: Graph{
				V: 3,
				Edges: []Edge{
					{0, 1, 10},
					{1, 2, 10},
					{2, 0, 10},
				},
			},
			expectedWeight: 20, // 10 + 10
			expectedEdges:  2,
		},
		{
			name: "Complex Graph (Wikipedia Example)",
			// A slightly denser graph to ensure robustness
			graph: Graph{
				V: 6,
				Edges: []Edge{
					{0, 1, 4}, {0, 2, 4}, // 0 connects to 1 and 2
					{1, 2, 2},            // Cycle 0-1-2
					{1, 0, 4},            // Duplicate inverted edge (robustness test)
					{2, 3, 3}, {2, 5, 2}, {2, 4, 4},
					{3, 5, 3}, {4, 5, 3},
				},
			},
			// Expected MST:
			// (1-2):2, (2-5):2, (2-3):3 or (3-5):3, (4-5):3, (0-1):4 or (0-2):4
			// Sum: 2+2+3+3+4 = 14
			expectedWeight: 14,
			expectedEdges:  5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inputEdges := slices.Clone(tt.graph.Edges)
			graphInput := Graph{V: tt.graph.V, Edges: inputEdges}

			result, totalWeight := KruskalMST(graphInput)

			if totalWeight != tt.expectedWeight {
				t.Errorf("Incorrect weight. Expected: %d, Got: %d", tt.expectedWeight, totalWeight)
			}

			if len(result) != tt.expectedEdges {
				t.Errorf("Incorrect number of edges. Expected: %d, Got: %d", tt.expectedEdges, len(result))
			}

			// Extra check: Check that there are no cycles in the result
			if hasCycle(tt.graph.V, result) {
				t.Errorf("The resulting MST contains a cycle!")
			}
		})
	}
}

func FuzzKruskalMST(f *testing.F) {
	// 1. Seed Corpus: Add known cases for the fuzzer to start
	// The fuzzer uses []byte, so let's create a simple format:
	// Byte 0: Number of Nodes
	// Bytes 1..N: Triplets (Source, Dest, Weight)
	f.Add([]byte{4, 0, 1, 10, 0, 2, 6, 0, 3, 5, 1, 3, 15, 2, 3, 4}) // Unit test case

	f.Fuzz(func(t *testing.T, data []byte) {
		// We need at least 1 byte to define the number of nodes
		if len(data) < 1 {
			return
		}

		// Limit V between 2 and 255 so the test doesn't get too heavy
		numNodes := max(int(data[0]), 2)

		// The rest of the bytes are the edges. Read in chunks of 3.
		// Format: [source, dest, weight]
		rawEdges := data[1:]
		var edges []Edge

		for i := 0; i+2 < len(rawEdges); i += 3 {
			u := int(rawEdges[i]) % numNodes   // Ensures node exists
			v := int(rawEdges[i+1]) % numNodes // Ensures node exists
			w := int(rawEdges[i+2])            // Weight (0-255)

			// Avoid self-loops (optional, but Kruskal should handle it)
			if u != v {
				edges = append(edges, Edge{Source: u, Dest: v, Weight: w})
			}
		}

		graphInput := Graph{V: numNodes, Edges: edges}

		mst, totalWeight := KruskalMST(graphInput)


		// Invariant A: The number of edges can never exceed V-1
		if len(mst) > numNodes-1 {
			t.Errorf("MST has too many edges! V=%d, Edges=%d", numNodes, len(mst))
		}

		// Invariant B: The returned total weight must be the exact sum of the edges
		calculatedSum := 0
		for _, a := range mst {
			calculatedSum += a.Weight
		}
		if calculatedSum != totalWeight {
			t.Errorf("Weight divergence! Returned: %d, Calculated: %d", totalWeight, calculatedSum)
		}

		// Invariant C: The resulting MST CANNOT have cycles
		// (Reusing the hasCycle function we created earlier)
		if hasCycle(numNodes, mst) {
			// If the fuzzer fails here, it will save the input to a file for debugging
			t.Errorf("The generated MST contains a cycle! Input: %v", data)
		}
	})
}

// Simple helper to validate if the generated MST has cycles (which would be a fatal error)
func hasCycle(V int, edges []Edge) bool {
	parent := make([]int, V)
	for i := range V {
		parent[i] = i
	}

	// Reusing the find logic from the main code
	// If we need to unite two nodes that already have the same parent, a cycle exists.
	for _, a := range edges {
		rootX := findIterative(parent, a.Source)
		rootY := findIterative(parent, a.Dest)

		if rootX == rootY {
			return true // Cycle detected
		}

		// Simple union for testing
		parent[rootX] = rootY
	}
	return false
}

// Helper: Updated for rand/v2
func generateRandomGraph(numNodes, numEdges int) Graph {
	graph := Graph{V: numNodes, Edges: make([]Edge, numEdges)}
	for i := range numEdges {
		graph.Edges[i] = Edge{
			Source: rand.IntN(numNodes), // rand.IntN is the new v2 API
			Dest:   rand.IntN(numNodes),
			Weight: rand.IntN(100) + 1,
		}
	}
	return graph
}

func BenchmarkKruskalLarge(b *testing.B) {
	numNodes := 1000
	numEdges := 5000
	baseGraph := generateRandomGraph(numNodes, numEdges)

	edgesBuffer := make([]Edge, len(baseGraph.Edges))
	
	for b.Loop() {
		b.StopTimer()

		copy(edgesBuffer, baseGraph.Edges)
		testGraph := Graph{V: numNodes, Edges: edgesBuffer}

		b.StartTimer()

		KruskalMST(testGraph)
	}
}