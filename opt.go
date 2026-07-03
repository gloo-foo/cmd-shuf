package command

// ShufCount sets the maximum number of output lines (-n flag).
type ShufCount int

// ShufSeed sets the random seed for deterministic output (--seed flag).
type ShufSeed int64

// ShufBound is an inclusive endpoint of the -i integer range.
type ShufBound int

// ShufRangeSpec is the -i option value: shuffle the integers of an inclusive
// range instead of stdin. Build it with ShufRange.
type ShufRangeSpec struct{ lo, hi int }

// ShufRange generates integers from lo to hi (inclusive), shuffled (-i flag).
func ShufRange(lo, hi ShufBound) ShufRangeSpec { return ShufRangeSpec{lo: int(lo), hi: int(hi)} }

// ShufLine is one input line supplied as a command argument (-e flag).
type ShufLine string

// ShufEchoArgs is the -e option value: the argument lines to shuffle instead of
// stdin. Build it with ShufEcho.
type ShufEchoArgs struct{ args []ShufLine }

// ShufEcho treats the given arguments as input lines to shuffle (-e flag).
func ShufEcho(args ...ShufLine) ShufEchoArgs { return ShufEchoArgs{args: args} }

// flags holds the parsed shuf options. source defaults to defaultShuffler and is
// overridden only by injection in tests.
type flags struct {
	inputRange *ShufRangeSpec
	echo       *ShufEchoArgs
	source     shufflerFor
	seed       seedSpec
	count      ShufCount
}

// with folds one option value into the flags, returning the updated copy.
// Values of unrecognized types are ignored.
func (f flags) with(o any) flags {
	switch v := o.(type) {
	case ShufCount:
		f.count = v
	case ShufSeed:
		f.seed = seedSpec{value: int64(v), isSet: true}
	case ShufRangeSpec:
		f.inputRange = &v
	case ShufEchoArgs:
		f.echo = &v
	case srcOption:
		f.source = v.factory
	}
	return f
}
