package stat

import "fmt"

func NewTours() *Tours {
	tours := &Tours{
		Total:     make(chan uint64),
		Skipped:   make(chan uint64),
		KidsIssue: make(chan uint64),
	}

	go tours.collect()

	return tours
}

type Tours struct {
	total uint64
	Total chan uint64

	skipped uint64
	Skipped chan uint64

	kidsIssue uint64
	KidsIssue chan uint64
}

func (t *Tours) collect() {
	for {
		select {
		case total := <-t.Total:
			t.total += total
		case skipped := <-t.Skipped:
			t.skipped += skipped
		case kidsIssue := <-t.KidsIssue:
			t.kidsIssue += kidsIssue
		default:
		}
	}
}

func (t *Tours) Output() {
	fmt.Printf("Tours: %d, Skipped: %d, Send: %d, Kids Issue: %d\n", t.total, t.skipped, (t.total - t.skipped), t.kidsIssue)
}
