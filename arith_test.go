package interval

import "testing"

func TestNeg(t *testing.T) {
	for _, test := range []struct{ in, want *Interval }{
		{&Interval{0, 0, Closed}, &Interval{0, 0, Closed}},
		{&Interval{0, 1, LeftClosed}, &Interval{-1, 0, RightClosed}},
		{&Interval{0, 1, RightClosed}, &Interval{-1, 0, LeftClosed}},
		{&Interval{0, 1, Open}, &Interval{-1, 0, Open}},
		{&Interval{2, 4, Closed}, &Interval{-4, -2, Closed}},
		{&Interval{2, 4, LeftClosed}, &Interval{-4, -2, RightClosed}},
		{&Interval{2, 4, RightClosed}, &Interval{-4, -2, LeftClosed}},
		{&Interval{2, 4, Open}, &Interval{-4, -2, Open}},
		{&Interval{neginf, 3, RightClosed}, &Interval{-3, inf, LeftClosed}},
		{&Interval{neginf, inf, Open}, &Interval{neginf, inf, Open}},
	} {
		if got := test.in.Neg(); !Equal(got, test.want) {
			t.Errorf("%v.Neg(): got %v, want %v", test.in, got, test.want)
		}
		if got := test.want.Neg(); !Equal(got, test.in) {
			t.Errorf("%v.Neg(): got %v, want %v", test.want, got, test.in)
		}
	}
}

var (
	ine  = empty()
	inz  = &Interval{0, 0, Closed}
	inp0 = &Interval{0, 0.5, Closed}
	inp1 = &Interval{1, 2, Closed}
	inm  = &Interval{-2, 4, Closed}
	inn0 = &Interval{-0.25, 0, Closed}
	inn1 = &Interval{-8, -4, Closed}
)

var arithTests = []struct{ x, y, add, sub *Interval }{
	{ine, inp1, ine, ine},
	{inp1, ine, ine, ine},
	{inp1, inz, inp1, inp1},
	{inz, inp1, inp1, inp1.Neg()},
	{inp0, inp1, &Interval{1, 2.5, Closed}, &Interval{-2, -0.5, Closed}},
	{inp1, inp0, &Interval{1, 2.5, Closed}, &Interval{0.5, 2, Closed}},
	{inp1, inm, &Interval{-1, 6, Closed}, &Interval{-3, 4, Closed}},
	{inp1, inn0, &Interval{0.75, 2, Closed}, &Interval{1, 2.25, Closed}},
	{inp1, inn1, &Interval{-7, -2, Closed}, &Interval{5, 10, Closed}},
	{inm, inp1, &Interval{-1, 6, Closed}, &Interval{-4, 3, Closed}},
	{inm, inp0, &Interval{-2, 4.5, Closed}, &Interval{-2.5, 4, Closed}},
	{inm, inm, &Interval{-4, 8, Closed}, &Interval{-6, 6, Closed}},
	{inm, inn0, &Interval{-2.25, 4, Closed}, &Interval{-2, 4.25, Closed}},
	{inm, inn1, &Interval{-10, 0, Closed}, &Interval{2, 12, Closed}},
	{inn1, inp1, &Interval{-7, -2, Closed}, &Interval{-10, -5, Closed}},
	{inn1, inp0, &Interval{-8, -3.5, Closed}, &Interval{-8.5, -4, Closed}},
	{inn1, inm, &Interval{-10, 0, Closed}, &Interval{-12, -2, Closed}},
	{inn1, inn0, &Interval{-8.25, -4, Closed}, &Interval{-8, -3.75, Closed}},
	{inn0, inn1, &Interval{-8.25, -4, Closed}, &Interval{3.75, 8, Closed}},
}

func TestAdd(t *testing.T) {
	for _, test := range arithTests {
		if got := Add(test.x, test.y); !Equal(got, test.add) {
			t.Errorf("Add(%v, %v): got %v, want %v", test.x, test.y, got, test.add)
		}
	}
}

func TestSub(t *testing.T) {
	for _, test := range arithTests {
		if got := Sub(test.x, test.y); !Equal(got, test.sub) {
			t.Errorf("Sub(%v, %v): got %v, want %v", test.x, test.y, got, test.sub)
		}
	}
}
