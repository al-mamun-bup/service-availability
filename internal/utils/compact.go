package utils

import (
	h3 "github.com/uber/h3-go/v4"
)

// CustomCompact simulates the behavior of H3's CompactCells using custom logic implementation
func CustomCompact(cells []h3.Cell, res int) []h3.Cell {
	current := make(map[h3.Cell]struct{})
	for _, c := range cells {
		current[c] = struct{}{}
	}

	for {
		parentToChildren := make(map[h3.Cell][]h3.Cell)

		// Group cells by their parent
		for cell := range current {
			parent, err := cell.Parent(res - 1)
			if err != nil {
				continue 
			}
			parentToChildren[parent] = append(parentToChildren[parent], cell)
		}

		updated := make(map[h3.Cell]struct{})
		changed := false

		// Attempt to compact cells by checking if a full set of children exists
		for parent, children := range parentToChildren {
			allChildren, err := parent.Children(res)
			if err != nil {
				// Keep children as-is if we can't get all children
				for _, child := range children {
					updated[child] = struct{}{}
				}
				continue
			}
			if len(children) == len(allChildren) {
				// If all children are present, use parent
				updated[parent] = struct{}{}
				changed = true
			} else {
				// Otherwise keep the existing children
				for _, child := range children {
					updated[child] = struct{}{}
				}
			}
		}

		if !changed {
			break
		}

		current = updated
		res-- // Move to a coarser resolution
	}

	// Convert map to slice
	result := make([]h3.Cell, 0, len(current))
	for c := range current {
		result = append(result, c)
	}
	return result
}
