// package interval implements floating-point interval arithmetic.
// TODO: Handle signed zero
// TODO: Implement outward rounding
package interval

import (
	"fmt"
	"math"
)

// An Interval is a subset of the extended real numbers.
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

// New returns a pointer to an Interval with endpoints x and y.
// Ends describes whether the endpoints are open or closed.
// New panics if the interval is empty or if it contains a closed endpoint of infinite value.
func New(x, y float64, ends Ends) *Interval {
	in := &Interval{x, y, ends}
	if in.IsEmpty() {
		panic(fmt.Sprintf("New: %v is empty", in))
	}
	if in.a == neginf && in.LeftIsClosed() {
		panic(fmt.Sprintf("New: %v is closed at -Inf", in))
	}
	if in.b == inf && in.RightIsClosed() {
		panic(fmt.Sprintf("New: %v is closed at +Inf", in))
	}
	return in
}

// NewUnit is shorthand for New(x, x, Closed).
func NewUnit(x float64) *Interval { return New(x, x, Closed) }

// Left returns in's left endpoint.
func (in *Interval) Left() float64 { return in.a }

// Right returns in's right endpoint.
func (in *Interval) Right() float64 { return in.b }

// Ends returns in's Ends.
func (in *Interval) Ends() Ends { return in.ends }

// IsEmpty reports whether in is an empty interval.
func (in *Interval) IsEmpty() bool { return in.a > in.b || in.a == in.b && in.ends != Closed }

// IsMixed reports whether in contains at least one positive and one negative real number.
func (in *Interval) IsMixed() bool { return in.a < 0 && 0 < in.b }

// IsUnit reports whether in represents a single real value.
func (in *Interval) IsUnit() bool { return in.a == in.b && in.ends == Closed }

// IsZero reports whether in is a closed interval containing only 0.
func (in *Interval) IsZero() bool { return in.IsUnit() && in.a == 0 }

// Contains reports whether in contains x.
// TODO: Infinities
func (in *Interval) Contains(x float64) bool {
	return (in.a < x || in.a == x && in.LeftIsClosed()) && (in.b > x || in.b == x && in.RightIsClosed())
}

// LeftIsClosed reports whether in contains its left endpoint,
// that is, whether its Ends is Closed or LeftClosed.
func (in *Interval) LeftIsClosed() bool { return in.ends&leftEndMask != 0 }

// RightIsClosed reports whether in contains its right endpoint,
// that is, whether its Ends is Closed or RightClosed.
func (in *Interval) RightIsClosed() bool { return in.ends&rightEndMask != 0 }

// empty returns the empty interval (0, 0).
func empty() *Interval { return &Interval{} }

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
