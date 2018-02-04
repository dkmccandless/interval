/*
Package interval implements floating-point interval arithmetic and set operations.

Define intervals using New and NewSingle:

	interval.New(0, 1, Closed)               // the unit interval
	interval.New(0, math.Inf(1), LeftClosed) // the non-negative reals
	interval.NewSingle(-3)                   // the number -3

The arithmetic functions in this package are based on Hickey, Ju, and van Emden,
"Interval Arithmetic: from Principles to Implementation", in particular
their definition of the "functional division" operation. Addition, subtraction,
and multiplication of non-empty intervals yield non-empty interval results.
Division is undefined when the denominator is [0, 0] and results in a union of
disjoint unbounded intervals when the denominator contains zero but the numerator
does not. Otherwise, when both arguments are non-empty intervals, division yields
a single non-empty interval result (which may be unbounded).
Set operations on intervals are defined when their intersection is non-empty.
Operations on empty intervals are semantically undefined and yield an empty
interval result.

Default hardware rounding of floating-point operations involving
interval endpoints may lead to imprecise and potentially incorrect
representation of the values contained in the interval.
*/
package interval
