package utils

import (
	"fmt"

	h3 "github.com/uber/h3-go/v4"
)

func CustomCompact2(cells []h3.Cell, res int) []h3.Cell {
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
	fmt.Println(result)
	return result
}
func PolygonToH3Indexes(boundary []h3.LatLng) []string {
    gPoly := h3.GeoPolygon{
        GeoLoop: h3.GeoLoop{},
    }
    for _, c := range boundary {
        gPoly.GeoLoop = append(gPoly.GeoLoop, h3.LatLng{
            Lat: c.Lat,
            Lng: c.Lng,
        })
    }
    return customPolyFillUsingRange(gPoly, 7, 9)
}

func customPolyFillUsingRange(geoPoly h3.GeoPolygon, minResolution, maxResolution int) []string {
    polyH3Map := map[int][]h3.Cell{}
    for i := minResolution; i <= maxResolution; i++ {
        polyH3Map[i], _ = geoPoly.Cells(i)

        fmt.Printf("resolution %d h3 indexes len: %d\n", i, len(polyH3Map[i]))
    }

    compactedMap := map[string]bool{}
    for i := maxResolution; i >= minResolution; i-- {
        for _, idx := range polyH3Map[i] {
            if i < maxResolution {
                childs, _ := idx.ImmediateChildren()
                all := true
                for _, childIdx := range childs {
                    if _, has := compactedMap[childIdx.String()]; !has {
                        all = false
                    }
                }

                if all {
                    compactedMap[idx.String()] = true
                    for _, childIdx := range childs {
                        delete(compactedMap, childIdx.String())
                    }
                }
            } else {
                compactedMap[idx.String()] = true
            }
        }

        if i < maxResolution {
           innerCell := []h3.Cell{}
           for _, idx := range polyH3Map[i+1] {
               rings , _ := idx.GridDisk(1)
               for _, rIdx := range rings {
                   if compactedMap[rIdx.String()] {
                       innerCell = append(innerCell, idx)
                       break
                   }
               }
           }
        
           for _, idx := range innerCell {
               compactedMap[idx.String()] = true
           }
        
           newH3 := []h3.Cell{}
           for _, idx := range polyH3Map[i] {
               if compactedMap[idx.String()] {
                   newH3 = append(newH3, idx)
               }
           }
           polyH3Map[i] = newH3
        }
    }

    fmt.Println("Compacted h3 indexes len: ", len(compactedMap))
    compacted := []string{}
    for k := range compactedMap {
        compacted = append(compacted, k)
    }
    return compacted
}
