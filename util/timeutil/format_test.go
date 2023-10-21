package timeutil

import (
	"testing"
	"time"
)

func TestIsToday(t *testing.T) {
	today := time.Now()
	want := true
	if IsToday(today) != want {
		t.Errorf("IsToday(%v) = false, want true", today)
	}
}
func TestIsNotToday(t *testing.T) {
	yesterday := time.Now().AddDate(0, 0, -1)
	want := false
	if IsToday(yesterday) != want {
		t.Errorf("IsToday(%v) = true, want false", yesterday)
	}
}
