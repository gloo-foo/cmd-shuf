package command

import gloo "github.com/gloo-foo/framework"

// ShufCount sets the maximum number of output lines (-n flag).
type ShufCount int

// ShufSeed sets the random seed for deterministic output (--seed flag).
type ShufSeed int64

// shufRangeFlag generates integers from lo to hi, ignoring stdin (-i flag).
type shufRangeFlag struct{ lo, hi int }

// ShufRange generates integers from lo to hi (inclusive), shuffled.
func ShufRange(lo, hi int) gloo.Switch[flags] { return shufRangeFlag{lo: lo, hi: hi} }

func (f shufRangeFlag) Configure(flags *flags) { flags.inputRange = &f }

// shufEchoFlag treats the given args as input lines instead of reading stdin (-e flag).
type shufEchoFlag struct{ args []string }

// ShufEcho treats the given arguments as input lines to shuffle.
func ShufEcho(args ...string) gloo.Switch[flags] { return shufEchoFlag{args: args} }

func (f shufEchoFlag) Configure(flags *flags) { flags.echo = &f }

// shuffle permutes the first n elements by calling swap, matching the signature
// of math/rand's Shuffle. Injecting it as a value makes the random source a
// collaborator, so the otherwise non-deterministic default path is testable.
type shuffle func(n int, swap func(i, j int))

// shufflerFor builds a shuffle for the given seed: deterministic when a seed is
// present, process-default randomness otherwise.
type shufflerFor func(seed *int64) shuffle

// srcOption overrides the random source factory (test-only injection seam).
type srcOption struct{ factory shufflerFor }

func (s srcOption) Configure(flags *flags) { flags.source = s.factory }

// flags holds the parsed shuf options. source defaults to defaultShuffler and is
// overridden only by injection in tests.
type flags struct {
	seed       *int64
	inputRange *shufRangeFlag
	echo       *shufEchoFlag
	source     shufflerFor
	count      ShufCount
}

func (c ShufCount) Configure(f *flags) { f.count = c }
func (s ShufSeed) Configure(f *flags)  { v := int64(s); f.seed = &v }
