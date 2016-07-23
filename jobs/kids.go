package jobs

import (
	"sort"
)

type KidsSlice []int

func (p KidsSlice) Len() int { return len(p) }

func (p KidsSlice) Less(i, j int) bool {
	if p[i] < 0 && p[j] >= 0 {
		return false
	} else if p[i] >= 0 && p[j] < 0 {
		return true
	} else {
		return p[i] < p[j]
	}
}

func (p KidsSlice) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

// Sort is a convenience method.
func (p KidsSlice) Sort() { sort.Sort(p) }
