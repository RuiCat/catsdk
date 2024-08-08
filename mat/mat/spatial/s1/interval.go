// Copyright 2014 Google Inc. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package s1

import (
	"math"
	"strconv"
)

// An Interval represents a closed interval on a unit circle (also known
// as a 1-dimensional sphere). It is capable of representing the empty
// interval (containing no points), the full interval (containing all
// points), and zero-length intervals (containing a single point).
//
// Points are represented by the angle they make with the positive x-axis in
// the range [-π, π]. An interval is represented by its lower and upper
// bounds (both inclusive, since the interval is closed). The lower bound may
// be greater than the upper bound, in which case the interval is "inverted"
// (i.e. it passes through the point (-1, 0)).
//
// The point (-1, 0) has two valid representations, π and -π. The
// normalized representation of this point is π, so that endpoints
// of normal intervals are in the range (-π, π]. We normalize the latter to
// the former in IntervalFromEndpoints. However, we take advantage of the point
// -π to construct two special intervals:
//
//	The full interval is [-π, π]
//	The empty interval is [π, -π].
//
// Treat the exported fields as read-only.
type Interval struct {
	Min, Max float64
}

// IntervalFromEndpoints constructs a new interval from endpoints.
// Both arguments must be in the range [-π,π]. This function allows inverted intervals
// to be created.
func IntervalFromEndpoints(lo, hi float64) Interval {
	i := Interval{lo, hi}
	if lo == -math.Pi && hi != math.Pi {
		i.Min = math.Pi
	}
	if hi == -math.Pi && lo != math.Pi {
		i.Max = math.Pi
	}
	return i
}

// IntervalFromPointPair returns the minimal interval containing the two given points.
// Both arguments must be in [-π,π].
func IntervalFromPointPair(a, b float64) Interval {
	if a == -math.Pi {
		a = math.Pi
	}
	if b == -math.Pi {
		b = math.Pi
	}
	if positiveDistance(a, b) <= math.Pi {
		return Interval{a, b}
	}
	return Interval{b, a}
}

// EmptyInterval returns an empty interval.
func EmptyInterval() Interval { return Interval{math.Pi, -math.Pi} }

// FullInterval returns a full interval.
func FullInterval() Interval { return Interval{-math.Pi, math.Pi} }

// IsValid reports whether the interval is valid.
func (i Interval) IsValid() bool {
	return (math.Abs(i.Min) <= math.Pi && math.Abs(i.Max) <= math.Pi &&
		!(i.Min == -math.Pi && i.Max != math.Pi) &&
		!(i.Max == -math.Pi && i.Min != math.Pi))
}

// IsFull reports whether the interval is full.
func (i Interval) IsFull() bool { return i.Min == -math.Pi && i.Max == math.Pi }

// IsEmpty reports whether the interval is empty.
func (i Interval) IsEmpty() bool { return i.Min == math.Pi && i.Max == -math.Pi }

// IsInverted reports whether the interval is inverted; that is, whether Lo > Hi.
func (i Interval) IsInverted() bool { return i.Min > i.Max }

// Invert returns the interval with endpoints swapped.
func (i Interval) Invert() Interval {
	return Interval{i.Max, i.Min}
}

// Center returns the midpoint of the interval.
// It is undefined for full and empty intervals.
func (i Interval) Center() float64 {
	c := 0.5 * (i.Min + i.Max)
	if !i.IsInverted() {
		return c
	}
	if c <= 0 {
		return c + math.Pi
	}
	return c - math.Pi
}

// Length returns the length of the interval.
// The length of an empty interval is negative.
func (i Interval) Length() float64 {
	l := i.Max - i.Min
	if l >= 0 {
		return l
	}
	l += 2 * math.Pi
	if l > 0 {
		return l
	}
	return -1
}

// Assumes p ∈ (-π,π].
func (i Interval) fastContains(p float64) bool {
	if i.IsInverted() {
		return (p >= i.Min || p <= i.Max) && !i.IsEmpty()
	}
	return p >= i.Min && p <= i.Max
}

// Contains returns true iff the interval contains p.
// Assumes p ∈ [-π,π].
func (i Interval) Contains(p float64) bool {
	if p == -math.Pi {
		p = math.Pi
	}
	return i.fastContains(p)
}

// ContainsInterval returns true iff the interval contains oi.
func (i Interval) ContainsInterval(oi Interval) bool {
	if i.IsInverted() {
		if oi.IsInverted() {
			return oi.Min >= i.Min && oi.Max <= i.Max
		}
		return (oi.Min >= i.Min || oi.Max <= i.Max) && !i.IsEmpty()
	}
	if oi.IsInverted() {
		return i.IsFull() || oi.IsEmpty()
	}
	return oi.Min >= i.Min && oi.Max <= i.Max
}

// InteriorContains returns true iff the interior of the interval contains p.
// Assumes p ∈ [-π,π].
func (i Interval) InteriorContains(p float64) bool {
	if p == -math.Pi {
		p = math.Pi
	}
	if i.IsInverted() {
		return p > i.Min || p < i.Max
	}
	return (p > i.Min && p < i.Max) || i.IsFull()
}

// InteriorContainsInterval returns true iff the interior of the interval contains oi.
func (i Interval) InteriorContainsInterval(oi Interval) bool {
	if i.IsInverted() {
		if oi.IsInverted() {
			return (oi.Min > i.Min && oi.Max < i.Max) || oi.IsEmpty()
		}
		return oi.Min > i.Min || oi.Max < i.Max
	}
	if oi.IsInverted() {
		return i.IsFull() || oi.IsEmpty()
	}
	return (oi.Min > i.Min && oi.Max < i.Max) || i.IsFull()
}

// Intersects returns true iff the interval contains any points in common with oi.
func (i Interval) Intersects(oi Interval) bool {
	if i.IsEmpty() || oi.IsEmpty() {
		return false
	}
	if i.IsInverted() {
		return oi.IsInverted() || oi.Min <= i.Max || oi.Max >= i.Min
	}
	if oi.IsInverted() {
		return oi.Min <= i.Max || oi.Max >= i.Min
	}
	return oi.Min <= i.Max && oi.Max >= i.Min
}

// InteriorIntersects returns true iff the interior of the interval contains any points in common with oi, including the latter's boundary.
func (i Interval) InteriorIntersects(oi Interval) bool {
	if i.IsEmpty() || oi.IsEmpty() || i.Min == i.Max {
		return false
	}
	if i.IsInverted() {
		return oi.IsInverted() || oi.Min < i.Max || oi.Max > i.Min
	}
	if oi.IsInverted() {
		return oi.Min < i.Max || oi.Max > i.Min
	}
	return (oi.Min < i.Max && oi.Max > i.Min) || i.IsFull()
}

// Compute distance from a to b in [0,2π], in a numerically stable way.
func positiveDistance(a, b float64) float64 {
	d := b - a
	if d >= 0 {
		return d
	}
	return (b + math.Pi) - (a - math.Pi)
}

// Union returns the smallest interval that contains both the interval and oi.
func (i Interval) Union(oi Interval) Interval {
	if oi.IsEmpty() {
		return i
	}
	if i.fastContains(oi.Min) {
		if i.fastContains(oi.Max) {
			// Either oi ⊂ i, or i ∪ oi is the full interval.
			if i.ContainsInterval(oi) {
				return i
			}
			return FullInterval()
		}
		return Interval{i.Min, oi.Max}
	}
	if i.fastContains(oi.Max) {
		return Interval{oi.Min, i.Max}
	}

	// Neither endpoint of oi is in i. Either i ⊂ oi, or i and oi are disjoint.
	if i.IsEmpty() || oi.fastContains(i.Min) {
		return oi
	}

	// This is the only hard case where we need to find the closest pair of endpoints.
	if positiveDistance(oi.Max, i.Min) < positiveDistance(i.Max, oi.Min) {
		return Interval{oi.Min, i.Max}
	}
	return Interval{i.Min, oi.Max}
}

// Intersection returns the smallest interval that contains the intersection of the interval and oi.
func (i Interval) Intersection(oi Interval) Interval {
	if oi.IsEmpty() {
		return EmptyInterval()
	}
	if i.fastContains(oi.Min) {
		if i.fastContains(oi.Max) {
			// Either oi ⊂ i, or i and oi intersect twice. Neither are empty.
			// In the first case we want to return i (which is shorter than oi).
			// In the second case one of them is inverted, and the smallest interval
			// that covers the two disjoint pieces is the shorter of i and oi.
			// We thus want to pick the shorter of i and oi in both cases.
			if oi.Length() < i.Length() {
				return oi
			}
			return i
		}
		return Interval{oi.Min, i.Max}
	}
	if i.fastContains(oi.Max) {
		return Interval{i.Min, oi.Max}
	}

	// Neither endpoint of oi is in i. Either i ⊂ oi, or i and oi are disjoint.
	if oi.fastContains(i.Min) {
		return i
	}
	return EmptyInterval()
}

// AddPoint returns the interval expanded by the minimum amount necessary such
// that it contains the given point "p" (an angle in the range [-π, π]).
func (i Interval) AddPoint(p float64) Interval {
	if math.Abs(p) > math.Pi {
		return i
	}
	if p == -math.Pi {
		p = math.Pi
	}
	if i.fastContains(p) {
		return i
	}
	if i.IsEmpty() {
		return Interval{p, p}
	}
	if positiveDistance(p, i.Min) < positiveDistance(i.Max, p) {
		return Interval{p, i.Max}
	}
	return Interval{i.Min, p}
}

// Define the maximum rounding error for arithmetic operations. Depending on the
// platform the mantissa precision may be different than others, so we choose to
// use specific values to be consistent across all.
// The values come from the C++ implementation.
var (
	// epsilon is a small number that represents a reasonable level of noise between two
	// values that can be considered to be equal.
	epsilon = 1e-15
	// dblEpsilon is a smaller number for values that require more precision.
	dblEpsilon = 2.220446049e-16
)

// Expanded returns an interval that has been expanded on each side by margin.
// If margin is negative, then the function shrinks the interval on
// each side by margin instead. The resulting interval may be empty or
// full. Any expansion (positive or negative) of a full interval remains
// full, and any expansion of an empty interval remains empty.
func (i Interval) Expanded(margin float64) Interval {
	if margin >= 0 {
		if i.IsEmpty() {
			return i
		}
		// Check whether this interval will be full after expansion, allowing
		// for a rounding error when computing each endpoint.
		if i.Length()+2*margin+2*dblEpsilon >= 2*math.Pi {
			return FullInterval()
		}
	} else {
		if i.IsFull() {
			return i
		}
		// Check whether this interval will be empty after expansion, allowing
		// for a rounding error when computing each endpoint.
		if i.Length()+2*margin-2*dblEpsilon <= 0 {
			return EmptyInterval()
		}
	}
	result := IntervalFromEndpoints(
		math.Remainder(i.Min-margin, 2*math.Pi),
		math.Remainder(i.Max+margin, 2*math.Pi),
	)
	if result.Min <= -math.Pi {
		result.Min = math.Pi
	}
	return result
}

// ApproxEqual reports whether this interval can be transformed into the given
// interval by moving each endpoint by at most ε, without the
// endpoints crossing (which would invert the interval). Empty and full
// intervals are considered to start at an arbitrary point on the unit circle,
// so any interval with (length <= 2*ε) matches the empty interval, and
// any interval with (length >= 2*π - 2*ε) matches the full interval.
func (i Interval) ApproxEqual(other Interval) bool {
	// Full and empty intervals require special cases because the endpoints
	// are considered to be positioned arbitrarily.
	if i.IsEmpty() {
		return other.Length() <= 2*epsilon
	}
	if other.IsEmpty() {
		return i.Length() <= 2*epsilon
	}
	if i.IsFull() {
		return other.Length() >= 2*(math.Pi-epsilon)
	}
	if other.IsFull() {
		return i.Length() >= 2*(math.Pi-epsilon)
	}

	// The purpose of the last test below is to verify that moving the endpoints
	// does not invert the interval, e.g. [-1e20, 1e20] vs. [1e20, -1e20].
	return (math.Abs(math.Remainder(other.Min-i.Min, 2*math.Pi)) <= epsilon &&
		math.Abs(math.Remainder(other.Max-i.Max, 2*math.Pi)) <= epsilon &&
		math.Abs(i.Length()-other.Length()) <= 2*epsilon)

}

func (i Interval) String() string {
	// like "[%.7f, %.7f]"
	return "[" + strconv.FormatFloat(i.Min, 'f', 7, 64) + ", " + strconv.FormatFloat(i.Max, 'f', 7, 64) + "]"
}

// Complement returns the complement of the interior of the interval. An interval and
// its complement have the same boundary but do not share any interior
// values. The complement operator is not a bijection, since the complement
// of a singleton interval (containing a single value) is the same as the
// complement of an empty interval.
func (i Interval) Complement() Interval {
	if i.Min == i.Max {
		// Singleton. The interval just contains a single point.
		return FullInterval()
	}
	// Handles empty and full.
	return Interval{i.Max, i.Min}
}

// ComplementCenter returns the midpoint of the complement of the interval. For full and empty
// intervals, the result is arbitrary. For a singleton interval (containing a
// single point), the result is its antipodal point on S1.
func (i Interval) ComplementCenter() float64 {
	if i.Min != i.Max {
		return i.Complement().Center()
	}
	// Singleton. The interval just contains a single point.
	if i.Max <= 0 {
		return i.Max + math.Pi
	}
	return i.Max - math.Pi
}

// DirectedHausdorffDistance returns the Hausdorff distance to the given interval.
// For two intervals i and y, this distance is defined by
//
//	h(i, y) = max_{p in i} min_{q in y} d(p, q),
//
// where d(.,.) is measured along S1.
func (i Interval) DirectedHausdorffDistance(y Interval) Angle {
	if y.ContainsInterval(i) {
		return 0 // This includes the case i is empty.
	}
	if y.IsEmpty() {
		return Angle(math.Pi) // maximum possible distance on s1.
	}
	yComplementCenter := y.ComplementCenter()
	if i.Contains(yComplementCenter) {
		return Angle(positiveDistance(y.Max, yComplementCenter))
	}

	// The Hausdorff distance is realized by either two i.Max endpoints or two
	// i.Min endpoints, whichever is farther apart.
	hiHi := 0.0
	if IntervalFromEndpoints(y.Max, yComplementCenter).Contains(i.Max) {
		hiHi = positiveDistance(y.Max, i.Max)
	}

	loLo := 0.0
	if IntervalFromEndpoints(yComplementCenter, y.Min).Contains(i.Min) {
		loLo = positiveDistance(i.Min, y.Min)
	}

	return Angle(math.Max(hiHi, loLo))
}

// Project returns the closest point in the interval to the given point p.
// The interval must be non-empty.
func (i Interval) Project(p float64) float64 {
	if p == -math.Pi {
		p = math.Pi
	}
	if i.fastContains(p) {
		return p
	}
	// Compute distance from p to each endpoint.
	dlo := positiveDistance(p, i.Min)
	dhi := positiveDistance(i.Max, p)
	if dlo < dhi {
		return i.Min
	}
	return i.Max
}
