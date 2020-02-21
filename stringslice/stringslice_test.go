package stringslice

import (
	"testing"
)

/*
Diff
*/
func Test1(t *testing.T) {
	a := []string{"a", "b", "c", "d", "e", "f"}
	b := []string{"b", "d", "f", "g"}
	diff := Diff(a, b)
	for _, e := range a {
		if Contains(b, e) {
			if Contains(diff, e) {
				t.Fail()
			}
		} else {
			if !Contains(diff, e) {
				t.Fail()
			}
		}
	}
}

func Test2(t *testing.T) {
	a := []string{"a", "b", "c"}
	b := []string{"c", "b", "a"}
	if !Equal(b, Reverse(a)) {
		t.Fail()
	}
}

func Test3(t *testing.T) {
	a := []string{"a", "b", "c"}
	b := Clone(a)
	if !Equal(a, b) {
		t.Fail()
	}
}
