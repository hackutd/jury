package ranking_test

import (
	"server/ranking"
	"testing"
)

func TestUpdate(t *testing.T) {
	a, b, c, d, e, f := ranking.Update(20.0, 12.2, 4.2, 1.0, 3.11, 0.65)
	Assert(t, a, 20.29342378562617)
	Assert(t, b, 12.144888310192417)
	Assert(t, c, 4.28143039999674)
	Assert(t, d, 0.9529174440716865)
	Assert(t, e, 3.057070240002119)
	Assert(t, f, 0.6301076201202875)
}

func Assert(t *testing.T, actual any, expected any) {
	if actual != expected {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}
