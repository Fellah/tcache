package stat

import (
	"fmt"
	"time"
)

const (
	historySize = 24
)

func NewTours() *Tours {
	tours := &Tours{
		total: 0,
		Total:     make(chan uint64),
		skipped: 0,
		Skipped:   make(chan uint64),
		kidsIssue: 0,
		KidsIssue: make(chan uint64),
		End: make(chan bool),
		historyUsed: 0,
		currentHour: "",
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

	idle uint64
	Idle chan uint64

	End chan bool

	currentHour string

	history [historySize]string
	historyUsed int
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
		case idle := <-t.Idle:
			t.idle += idle
		default:
		}

		current_hour := time.Now().Format("2006-01-02 15:00:00")
		if current_hour != t.currentHour {
			t.nextHistory(t.sOutput())
			t.currentHour = current_hour
		}
	}
}

func (t *Tours) nextHistory(s string) {
	if t.historyUsed < historySize {
		t.history[t.historyUsed] = s
		t.historyUsed++
	} else {
		for i := 0; i < (historySize-1); i++ {
			t.history[i] = t.history[i+1]
		}
		t.history[historySize-1] = s
	}
	t.total = 0
	t.skipped = 0
	t.kidsIssue = 0
	t.idle = 0
}

func (t *Tours) Output() {
	println(t.sOutput())
	for i := t.historyUsed - 1; i > 0; i-- {
		println(t.history[i])
	}
}

func (t *Tours) sOutput() string {
	return fmt.Sprintf("HOUR: %s, Tours: %d, Skipped: %d, Send: %d, Kids Issue: %d, Idle times: %d", t.currentHour, t.total, t.skipped, (t.total - t.skipped), t.kidsIssue, t.Idle)
}
