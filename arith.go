package interval

import (
	"errors"
	"fmt"
)

// ErrDisjointUnion is returned when the result of a call to Div
// is a union of two disjoint intervals.
// This occurs when an interval not containing zero
// is divided by an interval containing zero.
var ErrDisjointUnion = errors.New("union of disjoint intervals")

// ErrDivByZero is returned when a non-empty interval is divided by
// the closed degenerate interval [0, 0].
var ErrDivByZero = errors.New("division by the zero interval")

func (e Ends) flip() Ends { return e&leftEndMask<<1 + e&rightEndMask>>1 }

// Classification functions after Hickey et al.: P contains at least one positive
// number and no negative numbers; P0 contains 0 and P1 does not. Likewise for N.
func (in *Interval) isP0() bool  { return in.a == 0 && in.b > 0 && in.LeftIsClosed() }
func (in *Interval) isP1() bool  { return in.b > 0 && !in.Contains(0) }
func (in *Interval) isPos() bool { return in.isP0() || in.isP1() }
func (in *Interval) isN0() bool  { return in.Neg().isP0() }
func (in *Interval) isN1() bool  { return in.Neg().isP1() }
func (in *Interval) isNeg() bool { return in.Neg().isPos() }

// [0, 0] must be treated as a special case
func (in *Interval) has0() bool { return in.isP0() || in.IsMixed() || in.isN0() }

// Neg returns the additive inverse of x.
func (in *Interval) Neg() *Interval {
	return &Interval{-in.b, -in.a, in.ends.flip()}
}

// Add returns the sum x+y.
// TODO: Handle outward rounding
func Add(x, y *Interval) *Interval {
	if x.IsEmpty() || y.IsEmpty() {
		return empty()
	}
	return &Interval{x.a + y.a, x.b + y.b, x.ends & y.ends}
}

// Sub returns the difference x-y.
// TODO: Handle outward rounding
func Sub(x, y *Interval) *Interval {
	if x.IsEmpty() || y.IsEmpty() {
		return empty()
	}
	return Add(x, y.Neg())
}

// Mul returns the product x*y.
// TODO: Handle outward rounding
func Mul(x, y *Interval) *Interval {
	switch {
	case x.IsEmpty() || y.IsEmpty():
		return empty()
	case x.IsZero() || y.IsZero():
		return zero()
	case x.isNeg():
		return Mul(x.Neg(), y).Neg()
	case y.isNeg():
		return Mul(y.Neg(), x).Neg()
	case x.isPos() && y.isPos():
		e := x.ends & y.ends
		if x.a == 0 && x.LeftIsClosed() || y.a == 0 && y.LeftIsClosed() {
			e |= leftEndMask
		}
		return &Interval{x.a * y.a, x.b * y.b, e}
	case x.isPos() && y.IsMixed():
		if x.RightIsClosed() {
			return &Interval{x.b * y.a, x.b * y.b, y.ends}
		} else {
			return &Interval{x.b * y.a, x.b * y.b, Open}
		}
	case x.IsMixed() && y.IsMixed():
		// Return an interval from min(x.a*y.a, x.a*y.b, x.b*y.a, x.b*y.b)
		// to max(x.a*y.a, x.a*y.b, x.b*y.a, x.b*y.b) with appropriate ends
		return Union(
			&Interval{x.a * y.b, x.b * y.b, x.ends&y.ends.flip()&leftEndMask + x.ends&y.ends&rightEndMask},
			&Interval{x.b * y.a, x.a * y.a, x.ends.flip()&y.ends&leftEndMask + x.ends.flip()&y.ends.flip()&rightEndMask},
		)
	case y.isPos():
		return Mul(y, x)
	default:
		panic(fmt.Sprintf("unhandled case %v*%v", x, y))
	}
}

// Div returns the quotient x/y, defined as the interval containing all values z
// for which there exist values a in x and b in y, with b != 0, such that z = a/b.
//
// If x and y both contain 0 and other values, then a and b might be arbitrarily
// close to zero, and Div returns (-inf, 0], [0, +inf), or (-inf, +inf) according to
// the signs present in x and y, with a nil error.
// If y contains 0 and x does not, then their quotient is a union of two disjoint
// intervals with endpoints -inf, x.a/y.a and x.a/y.b, +inf. In this case,
// Div returns (-inf, +inf) and ErrDisjointUnion to signal the loss of information
// in the return value.
//
// Other special cases are:
//	Div(empty, y) = Div(x, empty) = empty, nil
//	Div(x, [0, 0]) = empty, ErrDivByZero
//	Div([0, 0], y) = [0, 0], nil
func Div(x, y *Interval) (*Interval, error) {
	switch {
	case x.IsEmpty() || y.IsEmpty():
		return empty(), nil
	case y.IsZero():
		return empty(), ErrDivByZero
	case x.IsZero():
		return zero(), nil
	case x.isNeg():
		in, err := Div(x.Neg(), y)
		return in.Neg(), err
	case y.isNeg():
		in, err := Div(x, y.Neg())
		return in.Neg(), err

	// Hickey et al.'s sign convention:
	// When the left endpoint is zero, it is treated as +0 (reciprocal +inf),
	// and a zero-valued right endpoint is treated as -0 (reciprocal -inf).
	case x.has0() && y.IsMixed() || x.IsMixed() && y.has0():
		return &Interval{neginf, inf, Open}, nil
	case x.isP1() && y.IsMixed():
		// The quotient is the union of two disjoint intervals
		// with endpoints -inf, x.a/y.a and x.a/y.b, +inf;
		// return their enclosure.
		return &Interval{neginf, inf, Open}, ErrDisjointUnion
	case y.isP0():
		return &Interval{x.a / y.b, inf, x.ends & y.ends.flip() & leftEndMask}, nil
	// y is P1
	case x.isPos():
		return &Interval{x.a / y.b, x.b / y.a, x.ends & y.ends.flip()}, nil
	case x.IsMixed():
		return &Interval{x.a / y.a, x.b / y.a, x.ends&y.ends&leftEndMask + x.ends&y.ends.flip()&rightEndMask}, nil
	default:
		panic(fmt.Sprintf("unhandled case %v/%v", x, y))
	}
}
