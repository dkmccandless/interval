package interval

func (e Ends) flip() Ends { return e&leftEndMask<<1 + e&rightEndMask>>1 }

func (in *Interval) isP0() bool  { return in.a == 0 && in.b > 0 }
func (in *Interval) isP1() bool  { return in.a > 0 }
func (in *Interval) isPos() bool { return in.isP1() || in.isP0() }
func (in *Interval) isN0() bool  { return in.a < 0 && in.b == 0 }
func (in *Interval) isN1() bool  { return in.b < 0 }
func (in *Interval) isNeg() bool { return in.isN1() || in.isN0() }

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
		return Union(
			&Interval{x.a * y.b, x.b * y.b, x.ends&y.ends.flip()&leftEndMask + x.ends&y.ends&rightEndMask},
			&Interval{x.b * y.a, x.a * y.a, x.ends.flip()&y.ends&leftEndMask + x.ends.flip()&y.ends.flip()&rightEndMask},
		)
	case x.IsMixed() && y.isPos():
		return Mul(y, x)
	case x.isNeg():
		return Mul(x.Neg(), y).Neg()
	case y.isNeg():
		return Mul(y.Neg(), x).Neg()
	default:
		panic("unhandled case")
	}
}
