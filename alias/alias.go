// Package alias provides unprefixed names for shuf command flags.
// This allows users to import and use shorter names:
//
//	import "github.com/gloo-foo/cmd-shuf/alias"
//	shuf.Shuf(alias.Count(3), alias.Seed(42))
package alias

import (
	gloo "github.com/gloo-foo/framework"

	command "github.com/gloo-foo/cmd-shuf"
)

// Shuf is the command constructor; it forwards to command.Shuf.
func Shuf(opts ...any) gloo.Command[[]byte, []byte] { return command.Shuf(opts...) }

// Count sets the maximum number of output lines (-n flag).
type Count = command.ShufCount

// Seed sets the random seed for deterministic output (--seed flag).
type Seed = command.ShufSeed

// Range generates integers from lo to hi (inclusive), shuffled (-i flag).
func Range(lo, hi command.ShufBound) command.ShufRangeSpec { return command.ShufRange(lo, hi) }

// Echo treats the given arguments as input lines to shuffle (-e flag).
func Echo(args ...command.ShufLine) command.ShufEchoArgs { return command.ShufEcho(args...) }
