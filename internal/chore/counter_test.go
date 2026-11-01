package chore

import (
	"testing"

	chModel "donetick.com/core/internal/chore/model"
)

func ip(v int) *int { return &v }

func TestApplyCounterProgress(t *testing.T) {
	ch := &chModel.Chore{TargetCount: ip(30), RemainingCount: ip(30)}
	if a, r, err := applyCounterProgress(ch, ip(29)); err != nil || a != 29 || r != 1 {
		t.Fatalf("29 of 30: got %d %d %v", a, r, err)
	}
	ch.RemainingCount = ip(1)
	if a, r, _ := applyCounterProgress(ch, ip(5)); a != 1 || r != 0 {
		t.Fatalf("clamp: got %d %d", a, r)
	}
	for _, bad := range []int{0, -3} {
		if _, _, err := applyCounterProgress(ch, ip(bad)); err == nil {
			t.Fatalf("expected error for %d", bad)
		}
	}
	if a, r, _ := applyCounterProgress(ch, nil); a != 1 || r != 0 {
		t.Fatalf("complete all: got %d %d", a, r)
	}
	// legacy row with nil remaining falls back to target
	ch.RemainingCount = nil
	if _, r, _ := applyCounterProgress(ch, ip(10)); r != 20 {
		t.Fatalf("nil remaining: got %d", r)
	}
}

func TestEditedCounter(t *testing.T) {
	old := &chModel.Chore{TargetCount: ip(30), RemainingCount: ip(4)}
	if tg, r := editedCounter(old, ip(30)); *tg != 30 || *r != 4 {
		t.Fatal("same target should keep progress")
	}
	if tg, r := editedCounter(old, ip(10)); *tg != 10 || *r != 10 {
		t.Fatal("new target should reset")
	}
	if tg, r := editedCounter(old, nil); tg != nil || r != nil {
		t.Fatal("nil target removes counter")
	}
	if tg, _ := editedCounter(nil, ip(0)); tg != nil {
		t.Fatal("0 means no counter")
	}
}
