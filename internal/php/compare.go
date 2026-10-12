package php

import (
	"math"
	"strings"
)

func scalarOrder(a, b value) (int, bool) {
	// null/string is an exception to the general bool/null truth comparison.
	if a.kind == nullKind && b.kind == stringKind {
		return strings.Compare("", b.s), true
	}
	if a.kind == stringKind && b.kind == nullKind {
		return strings.Compare(a.s, ""), true
	}
	if a.kind == boolKind || b.kind == boolKind || a.kind == nullKind || b.kind == nullKind {
		// zend_compare preserves the NaN-as-true result without a conversion
		// warning, including PHP 8.5. Do not use evaluator.toBool here.
		return integerOrder(boolean(a.truthy()).i, boolean(b.truthy()).i), true
	}
	if a.kind == stringKind && b.kind == stringKind {
		return stringOrder(a, b), true
	}
	if a.kind == floatKind && math.IsNaN(a.f) || b.kind == floatKind && math.IsNaN(b.f) {
		return 0, false
	}
	if a.kind == stringKind {
		order, _ := scalarOrder(b, a)
		return -order, true
	}
	if b.kind == stringKind {
		if n, _, ok := comparisonNumber(b); ok {
			return numberOrder(a, n), true
		}
		// PHP 8 compares a number with a non-numeric string lexically,
		// converting the number to text using the default precision.
		return strings.Compare(a.text(), b.s), true
	}
	return numberOrder(a, b), true
}

func integerOrder(a, b int64) int {
	if a < b {
		return -1
	}
	if a > b {
		return 1
	}
	return 0
}

func numberOrder(a, b value) int {
	if a.kind == intKind && b.kind == intKind {
		return integerOrder(a.i, b.i)
	}
	if a.asFloat() < b.asFloat() {
		return -1
	}
	if a.asFloat() > b.asFloat() {
		return 1
	}
	return 0
}

// comparisonNumber requires the entire string to be numeric, apart from ASCII
// whitespace. Unlike arithmetic, a numeric prefix alone does not suffice and
// never raises a warning. overflow records a decimal integer outside int64.
func comparisonNumber(v value) (n value, overflow int, ok bool) {
	n, trailing, valid := v.numeric()
	if !valid || trailing {
		return value{}, 0, false
	}
	if n.kind == floatKind && !strings.ContainsAny(v.s, ".eE") {
		overflow = 1
		if n.f < 0 {
			overflow = -1
		}
	}
	return n, overflow, true
}

func stringOrder(a, b value) int {
	x, overflowX, numericX := comparisonNumber(a)
	y, overflowY, numericY := comparisonNumber(b)
	if numericX && numericY {
		// Zend's smart string comparison falls back to original bytes when
		// overflowing integer strings have lost their distinguishing bits,
		// or when both numeric strings overflow to the same infinity.
		lostIntegerPrecision := overflowX != 0 && overflowX == overflowY && x.f-y.f == 0
		sameInfinity := x.kind == floatKind && y.kind == floatKind && x.f == y.f && math.IsInf(x.f, 0)
		if !lostIntegerPrecision && !sameInfinity {
			if x.kind == intKind && overflowY != 0 {
				return -overflowY
			}
			if y.kind == intKind && overflowX != 0 {
				return overflowX
			}
			return numberOrder(x, y)
		}
	}
	return strings.Compare(a.s, b.s)
}
