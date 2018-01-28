package interval

import (
	"testing"
)

var boolTests = []struct {
	in                                                Interval
	empty, mixed, unit, zero, leftClosed, rightClosed bool
}{
	{Interval{0, 0, Open}, true, false, false, false, false, false},
	{Interval{0, 0, LeftClosed}, true, false, false, false, true, false},
	{Interval{0, 0, RightClosed}, true, false, false, false, false, true},
	{Interval{0, 0, Closed}, false, false, true, true, true, true},
	{Interval{1, 1, Open}, true, false, false, false, false, false},
	{Interval{1, 2, Open}, false, false, false, false, false, false},
	{Interval{-1, 1, Open}, false, true, false, false, false, false},
	{Interval{1, 1, Closed}, false, false, true, false, true, true},
	{Interval{1, 2, Closed}, false, false, false, false, true, true},
	{Interval{-1, 1, Closed}, false, true, false, false, true, true},
	{Interval{0, inf, LeftClosed}, false, false, false, false, true, false},
	{Interval{neginf, 0, Open}, false, false, false, false, false, false},
	{Interval{neginf, inf, Open}, false, true, false, false, false, false},
}

func TestIsEmpty(t *testing.T) {
	for _, test := range boolTests {
		if got := test.in.IsEmpty(); got != test.empty {
			t.Errorf("IsEmpty(%v): got %v, want %v", test.in, got, test.empty)
		}
	}
}

func TestIsMixed(t *testing.T) {
	for _, test := range boolTests {
		if got := test.in.IsMixed(); got != test.mixed {
			t.Errorf("IsMixed(%v): got %v, want %v", test.in, got, test.mixed)
		}
	}
}

func TestIsUnit(t *testing.T) {
	for _, test := range boolTests {
		if got := test.in.IsUnit(); got != test.unit {
			t.Errorf("IsUnit(%v): got %v, want %v", test.in, got, test.unit)
		}
	}
}

func TestIsZero(t *testing.T) {
	for _, test := range boolTests {
		if got := test.in.IsZero(); got != test.zero {
			t.Errorf("IsZero(%v): got %v, want %v", test.in, got, test.zero)
		}
	}
}

func TestLeftIsClosed(t *testing.T) {
	for _, test := range boolTests {
		if got := test.in.LeftIsClosed(); got != test.leftClosed {
			t.Errorf("LeftIsClosed(%v): got %v, want %v", test.in, got, test.leftClosed)
		}
	}
}

func TestRightIsClosed(t *testing.T) {
	for _, test := range boolTests {
		if got := test.in.RightIsClosed(); got != test.rightClosed {
			t.Errorf("RightIsClosed(%v): got %v, want %v", test.in, got, test.rightClosed)
		}
	}
}

func TestContains(t *testing.T) {
	for _, test := range []struct {
		in   Interval
		x    float64
		want bool
	}{
		{Interval{}, 0, false},
		{Interval{0, 0, Open}, 0, false},
		{Interval{0, 0, LeftClosed}, 0, false},
		{Interval{0, 0, RightClosed}, 0, false},
		{Interval{0, 0, Closed}, 0, true},
		{Interval{1, -1, Closed}, 0, false}, // a>b
		{Interval{2, 4, Open}, 2, false},
		{Interval{2, 4, LeftClosed}, 2, true},
		{Interval{2, 4, RightClosed}, 2, false},
		{Interval{2, 4, Closed}, 2, true},
		{Interval{2, 4, Open}, 2.718, true},
		{Interval{2, 4, LeftClosed}, 2.718, true},
		{Interval{2, 4, RightClosed}, 2.718, true},
		{Interval{2, 4, Closed}, 2.718, true},
		{Interval{2, 4, Open}, 4, false},
		{Interval{2, 4, LeftClosed}, 4, false},
		{Interval{2, 4, RightClosed}, 4, true},
		{Interval{2, 4, Closed}, 4, true},
		{Interval{2, 4, Open}, -3, false},
		{Interval{2, 4, LeftClosed}, -3, false},
		{Interval{2, 4, RightClosed}, -3, false},
		{Interval{2, 4, Closed}, -3, false},
		{Interval{neginf, 0, RightClosed}, -18, true},
		{Interval{neginf, 0, RightClosed}, 0, true},
		{Interval{neginf, 0, RightClosed}, 6, false},
		{Interval{0, inf, LeftClosed}, -18, false},
		{Interval{0, inf, LeftClosed}, 0, true},
		{Interval{0, inf, LeftClosed}, 6, true},
		{Interval{neginf, inf, Open}, -3, true},
		{Interval{neginf, inf, Open}, neginf, false},
		{Interval{neginf, inf, Open}, inf, false},
	} {
		if got := test.in.Contains(test.x); got != test.want {
			t.Errorf("Contains(%v, %v): got %v, want %v", test.in, test.x, got, test.want)
		}
	}
}
