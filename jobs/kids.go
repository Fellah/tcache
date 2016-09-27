package jobs

import (
	"github.com/fellah/tcache/data"
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

func processKidsValue(tour *data.Tour) {
	var kids int

	if tour.Kid1Age != nil {
		kids++
	} else {
		kidsAge := -1
		tour.Kid1Age = &kidsAge
	}

	if tour.Kid2Age != nil {
		kids++
	} else {
		kidsAge := -1
		tour.Kid2Age = &kidsAge
	}

	if tour.Kid3Age != nil {
		kids++
	} else {
		kidsAge := -1
		tour.Kid3Age = &kidsAge
	}

	if kids != tour.Kids {
		switch tour.Kids {
		case 0:
			*tour.Kid1Age, *tour.Kid2Age, *tour.Kid3Age = -1, -1, -1
		case 1:
			*tour.Kid2Age, *tour.Kid3Age = -1, -1
		case 2:
			*tour.Kid3Age = -1
		}
	}

	kidsSlice := make(KidsSlice, 3)

	kidsSlice[0] = *tour.Kid1Age
	kidsSlice[1] = *tour.Kid2Age
	kidsSlice[2] = *tour.Kid3Age

	kidsSlice.Sort()

	tour.Kid1Age = &kidsSlice[0]
	tour.Kid2Age = &kidsSlice[1]
	tour.Kid3Age = &kidsSlice[2]
}

func processKidAgeValue(kidAge int) (age int) {
	if kidAge >= 0 && kidAge <= 1 {
		// Variable 'age' equal zero by default.
	} else if kidAge >= 2 && kidAge <= 6 {
		age = 2
	} else if kidAge >= 7 && kidAge <= 8 {
		age = 7
	} else if kidAge >= 9 && kidAge <= 12 {
		age = 9
	} else if kidAge >= 13 {
		age = 13
	}

	return age
}

func isKidsValid(tour *data.Tour) bool {
	kids := 0

	if tour.Kid1Age != nil && *tour.Kid1Age >= 0 {
		kids++
	}

	if tour.Kid2Age != nil && *tour.Kid2Age >= 0 {
		kids++
	}

	if tour.Kid3Age != nil && *tour.Kid3Age >= 0 {
		kids++
	}

	if tour.Kids == kids {
		return true
	}

	return false
}
