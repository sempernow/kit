// Package bmask provides bitmask operators on its variable-aggregated type (uint64).
package bmask

import (
	"math/bits"
)

// Bitmask value for variably-aggregated (bitmap) type
// whereof each (named) member (max 64 members) is a "flag" having one non-zero bit.
// Bitmask operators handle the aggregated Bitmap values over its range of flags:
// math.Pow(2, 0) to math.Pow(2, 63).
type Bitmask uint64 //... CLIPS @ math.Pow(2, 63); does NOT cycle.

// ----------------------------------------------------------------------------
// Bitmask operators

// FlagIndex returns bit position (0-63).
func FlagIndex(flag Bitmask) int {
	return bits.Len64(uint64(flag)) - 1
} //... flag = math.Pow(2, float64(FlagIndex(flag)))

// HasFlag tests if flag is in receiver's Bitmask
func (r Bitmask) HasFlag(flag Bitmask) bool { return r&flag != 0 }

// AddFlag adds flag to receiver's Bitmask
func (r *Bitmask) AddFlag(flag Bitmask) { *r |= flag }

// ClearFlag removes flag from receiver's Bitmask
func (r *Bitmask) ClearFlag(flag Bitmask) { *r &= ^flag }

// ToggleFlag toggles flag at receiver's Bitmask
func (r *Bitmask) ToggleFlag(flag Bitmask) { *r ^= flag }
