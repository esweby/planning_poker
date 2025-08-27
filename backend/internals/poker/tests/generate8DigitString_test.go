package poker

import (
	"backend/internals/poker"
	"testing"
	"unicode"
)

func TestGenerate8DigitStringLength(t *testing.T) {
	for i := 0; i < 100; i++ {
		s := poker.Generate8DigitString()
		if len(s) != 8 {
			t.Errorf("expected length 8, got %d (%s)", len(s), s)
		}
	}
}

func TestGenerate8DigitStringDigits(t *testing.T) {
	for i := 0; i < 100; i++ {
		s :=poker.Generate8DigitString()
		for _, r := range s {
			if !unicode.IsDigit(r) {
				t.Errorf("expected digit, got %q in %s", r, s)
			}
		}
	}
}

func TestGenerate8DigitStringRandomness(t *testing.T) {
	seen := make(map[string]bool)
	for i := 0; i < 1000; i++ {
		s :=poker.Generate8DigitString()
		if seen[s] {
			t.Logf("duplicate detected: %s", s)
		}
		seen[s] = true
	}
	if len(seen) < 900 {
		t.Errorf("too many duplicates, randomness may be broken (unique=%d)", len(seen))
	}
}
