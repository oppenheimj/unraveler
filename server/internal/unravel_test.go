package internal

import (
	"reflect"
	"testing"
)

func TestCalculateBlocks(t *testing.T) {
	tables := []struct {
		numThreads     int
		numNodes       int
		expectedBlocks []int
	}{
		{8, 200, []int{0, 13, 27, 42, 59, 78, 100, 130, 200}},
		{4, 1000, []int{0, 134, 293, 500, 1000}},
		{8, 9, []int{0, 20}},
	}

	for _, table := range tables {
		actualBlocks := calculateBlocks(table.numThreads, table.numNodes)
		if !reflect.DeepEqual(actualBlocks, table.expectedBlocks) {
			t.Errorf("Indices for t=%d n=%d incorrect, got: %d, want: %d.", table.numThreads, table.numNodes, actualBlocks, table.expectedBlocks)
		}
	}
}
