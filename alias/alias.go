// Package alias provides unprefixed type aliases for shuf command flags.
// This allows users to import and use shorter names:
//
//	import "github.com/gloo-foo/cmd-shuf/alias"
//	shuf.Shuf(alias.Count(3), alias.Seed(42))
package alias

import command "github.com/gloo-foo/cmd-shuf"

// Shuf is the command constructor.
var Shuf = command.Shuf

// Count sets the maximum number of output lines (-n flag).
type Count = command.ShufCount

// Seed sets the random seed for deterministic output (--seed flag).
type Seed = command.ShufSeed

// Range generates integers from lo to hi (inclusive), shuffled (-i flag).
var Range = command.ShufRange

// Echo treats the given arguments as input lines to shuffle (-e flag).
var Echo = command.ShufEcho
