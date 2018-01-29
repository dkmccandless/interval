package interval

func (e Ends) flip() Ends { return e&leftEndMask<<1 + e&rightEndMask>>1 }

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
