package chore

import (
	"errors"

	chModel "donetick.com/core/internal/chore/model"
)

// normalizeTargetCount returns nil for "no counter" (nil or <= 0) and a copy otherwise.
func normalizeTargetCount(target *int) *int {
	if target == nil || *target <= 0 {
		return nil
	}
	v := *target
	return &v
}

// editedCounter computes the new target/remaining pair when a chore is edited.
// Changing the target restarts the counter; keeping it preserves progress.
func editedCounter(old *chModel.Chore, reqTarget *int) (*int, *int) {
	target := normalizeTargetCount(reqTarget)
	if target == nil {
		return nil, nil
	}
	if old != nil && old.TargetCount != nil && *old.TargetCount == *target && old.RemainingCount != nil {
		rem := *old.RemainingCount
		if rem > *target || rem <= 0 {
			rem = *target
		}
		return target, &rem
	}
	rem := *target
	return target, &rem
}

// applyCounterProgress validates a progress submission for a counted chore.
// A nil count means "complete all remaining". Counts above the remaining
// amount are clamped. It returns the applied amount and the new remaining.
func applyCounterProgress(chore *chModel.Chore, count *int) (int, int, error) {
	remaining := *chore.TargetCount
	if chore.RemainingCount != nil && *chore.RemainingCount > 0 {
		remaining = *chore.RemainingCount
	}
	if count == nil {
		return remaining, 0, nil
	}
	if *count <= 0 {
		return 0, remaining, errors.New("count must be a positive number")
	}
	applied := *count
	if applied > remaining {
		applied = remaining
	}
	return applied, remaining - applied, nil
}
