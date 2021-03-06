package interval

import (
	"errors"
	"fmt"
	"math"
)

// ErrNaN is returned when New is called with an argument that is NaN.
var ErrNaN = errors.New("argument is NaN")

// ErrEmpty is returned when New is called with arguments
// that would result in an empty interval.
var ErrEmpty = errors.New("empty interval")

// ErrClosedInf is returned when New is called with arguments
// that would create a closed left or right endpoint at -inf or +inf.
var ErrClosedInf = errors.New("closed endpoint of infinite value")

// An Interval is a subset of the real numbers.
// The Interval type's zero value corresponds to the empty interval (0, 0).
type Interval struct {
	a, b float64
	ends Ends
}

// An Ends describes whether an Interval contains its endpoints.
type Ends int

const (
	Open        Ends = iota
	LeftClosed       // right-open
	RightClosed      // left-open
	Closed
)

const (
	leftEndMask  Ends = 1
	rightEndMask Ends = 2
)

var (
	inf    = math.Inf(1)
	neginf = math.Inf(-1)
)

// New returns a pointer to an Interval with endpoints x and y,
// which may be positive or negative infinity.
// Ends describes whether the endpoints are open or closed.
// New returns an empty interval and a non-nil error
// if x or y is NaN or if the interval is empty
// or contains a closed endpoint of infinite value.
func New(x, y float64, ends Ends) (*Interval, error) {
	if math.IsNaN(x) || math.IsNaN(y) {
		return empty(), ErrNaN
	}
	in := &Interval{x, y, ends}
	if in.IsEmpty() {
		return empty(), ErrEmpty
	}
	if in.a == neginf && in.LeftIsClosed() || in.b == inf && in.RightIsClosed() {
		return empty(), ErrClosedInf
	}
	return in, nil
}

// NewSingle is shorthand for New(x, x, Closed).
func NewSingle(x float64) (*Interval, error) { return New(x, x, Closed) }

// Left returns in's left endpoint.
func (in *Interval) Left() float64 { return in.a }

// Right returns in's right endpoint.
func (in *Interval) Right() float64 { return in.b }

// Ends returns in's Ends.
func (in *Interval) Ends() Ends { return in.ends }

// IsEmpty reports whether in is an empty interval.
// An interval with endpoints x and y is empty if x > y
// or if x == y and either endpoint is open.
func (in *Interval) IsEmpty() bool { return in.a > in.b || in.a == in.b && in.ends != Closed }

// IsMixed reports whether in contains at least one positive and one negative real number.
func (in *Interval) IsMixed() bool { return in.a < 0 && 0 < in.b }

// IsSingle reports whether in is a degenerate interval
// representing a single real value.
func (in *Interval) IsSingle() bool { return in.a == in.b && in.ends == Closed }

// IsZero reports whether in is the closed degenerate interval [0, 0].
func (in *Interval) IsZero() bool { return in.IsSingle() && in.a == 0 }

// Contains reports whether in contains x.
func (in *Interval) Contains(x float64) bool {
	return (in.a < x || in.a == x && in.LeftIsClosed()) && (in.b > x || in.b == x && in.RightIsClosed())
}

// Equal reports whether x and y represent the same quantity.
// Two intervals are equal if they are both empty or if they contain the same values.
func Equal(x, y *Interval) bool {
	if x.IsEmpty() {
		return y.IsEmpty()
	}
	return x.a == y.a && x.b == y.b && x.ends == y.ends
}

// LeftIsClosed reports whether in contains its left endpoint,
// that is, whether its Ends is Closed or LeftClosed.
func (in *Interval) LeftIsClosed() bool { return in.ends&leftEndMask != 0 }

// RightIsClosed reports whether in contains its right endpoint,
// that is, whether its Ends is Closed or RightClosed.
func (in *Interval) RightIsClosed() bool { return in.ends&rightEndMask != 0 }

// empty returns the empty interval (0, 0).
func empty() *Interval { return &Interval{} }

// zero returns the closed interval [0, 0].
func zero() *Interval { return &Interval{0, 0, Closed} }

// Intersection returns the intersection of x and y.
func Intersection(x, y *Interval) *Interval {
	if x.IsEmpty() || y.IsEmpty() {
		return empty()
	}
	if Equal(x, y) {
		return &Interval{x.a, x.b, x.ends}
	}
	if !x.Contains(y.a) && !x.Contains(y.b) && !y.Contains(x.a) && !y.Contains(x.b) {
		return empty()
	}
	var e Ends
	switch {
	case x.a < y.a:
		e ^= y.ends & leftEndMask
	case x.a == y.a:
		e ^= x.ends & y.ends & leftEndMask
	case x.a > y.a:
		e ^= x.ends & leftEndMask
	}
	switch {
	case x.b < y.b:
		e ^= x.ends & rightEndMask
	case x.b == y.b:
		e ^= x.ends & y.ends & rightEndMask
	case x.b > y.b:
		e ^= y.ends & rightEndMask
	}
	return &Interval{math.Max(x.a, y.a), math.Min(x.b, y.b), e}
}

// Union returns the union of x and y if their intersection is non-empty,
// or else the empty interval.
func Union(x, y *Interval) *Interval {
	if x.IsEmpty() || y.IsEmpty() {
		return empty()
	}
	if Equal(x, y) {
		return &Interval{x.a, x.b, x.ends}
	}
	if !x.Contains(y.a) && !x.Contains(y.b) && !y.Contains(x.a) && !y.Contains(x.b) {
		return empty()
	}
	var e Ends
	switch {
	case x.a < y.a:
		e ^= x.ends & leftEndMask
	case x.a == y.a:
		e ^= (x.ends | y.ends) & leftEndMask
	case x.a > y.a:
		e ^= y.ends & leftEndMask
	}
	switch {
	case x.b < y.b:
		e ^= y.ends & rightEndMask
	case x.b == y.b:
		e ^= (x.ends | y.ends) & rightEndMask
	case x.b > y.b:
		e ^= x.ends & rightEndMask
	}
	return &Interval{math.Min(x.a, y.a), math.Max(x.b, y.b), e}
}

// String returns a string representation of in.
// Square brackets denote closed endpoints and parentheses denote open endpoints.
func (in *Interval) String() string {
	var l, r string
	if in.LeftIsClosed() {
		l = "["
	} else {
		l = "("
	}
	if in.RightIsClosed() {
		r = "]"
	} else {
		r = ")"
	}
	return fmt.Sprintf("%v%v, %v%v", l, in.a, in.b, r)
}
