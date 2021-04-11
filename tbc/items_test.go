package tbc

import "testing"

func TestColorIntersection(t *testing.T) {
	chryo := GemLookup["Rune Covered Chrysoprase"]

	if !chryo.Color.Intersects(GemColorBlue) {
		t.Fatalf("Chryo intersects blue...")
	}
	if !chryo.Color.Intersects(GemColorGreen) {
		t.Fatalf("Chryo intersects blue...")
	}
}
