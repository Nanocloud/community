package utils

import "testing"

func TestRandomness(t *testing.T) {
	a := RandomString(10)
	b := RandomString(10)

	if a == b {
		t.Errorf("RandomString shouldn't return the same values when called in a row")
	}
}

func TestLength(t *testing.T) {
	for i := 1; i < 50; i++ {
		l := len(RandomString(i))
		if l != i {
			t.Errorf("RandomString(%d) should return a %d characters long string. Got a %d characters long one.", i, i, l)
		}
	}
}
