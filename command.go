package command

import (
	"math/rand"
	"strconv"

	gloo "github.com/gloo-foo/framework"
	"github.com/gloo-foo/framework/patterns"
)

// Shuf returns a command that randomly shuffles input lines.
//   - ShufCount(n) (-n): cap the output at n lines.
//   - ShufSeed(s) (--seed): deterministic output from a fixed seed.
//   - ShufRange(lo, hi) (-i): shuffle the integers lo..hi instead of stdin.
//   - ShufEcho(args...) (-e): shuffle the given arguments instead of stdin.
func Shuf(opts ...any) gloo.Command[[]byte, []byte] {
	f := gloo.NewParameters[gloo.File, flags](opts...).Flags
	return patterns.Accumulate(func(lines [][]byte) ([][]byte, error) {
		shuffled := permute(resolveSource(f)(f.seed), inputLines(f, lines))
		return capCount(f.count, shuffled), nil
	})
}

// inputLines selects the lines to shuffle: an integer range (-i) and echo args
// (-e) each override stdin, with echo taking precedence to match GNU shuf.
func inputLines(f flags, stdin [][]byte) [][]byte {
	switch {
	case f.echo != nil:
		return echoLines(f.echo.args)
	case f.inputRange != nil:
		return rangeLines(*f.inputRange)
	default:
		return stdin
	}
}

// echoLines turns -e arguments into input lines.
func echoLines(args []string) [][]byte {
	lines := make([][]byte, len(args))
	for i, a := range args {
		lines[i] = []byte(a)
	}
	return lines
}

// rangeLines renders the inclusive integer range lo..hi as input lines.
func rangeLines(r shufRangeFlag) [][]byte {
	lines := make([][]byte, 0, r.hi-r.lo+1)
	for i := r.lo; i <= r.hi; i++ {
		lines = append(lines, []byte(strconv.Itoa(i)))
	}
	return lines
}

// resolveSource resolves the random-source factory, defaulting when none was
// injected.
func resolveSource(f flags) shufflerFor {
	if f.source != nil {
		return f.source
	}
	return defaultShuffler
}

// defaultShuffler is the production random source: seeded and deterministic when
// a seed is given, process-default randomness otherwise. shuf is a
// non-cryptographic line shuffler, so the weak RNG (gosec G404) is the correct,
// GNU-compatible behavior — the G404 exclusion in .golangci.yaml is scoped here.
func defaultShuffler(seed *int64) shuffle {
	if seed != nil {
		return rand.New(rand.NewSource(*seed)).Shuffle
	}
	return rand.Shuffle
}

// permute returns a shuffled copy of lines using the given source. The input
// slice is never mutated.
func permute(s shuffle, lines [][]byte) [][]byte {
	out := make([][]byte, len(lines))
	copy(out, lines)
	s(len(out), func(i, j int) {
		out[i], out[j] = out[j], out[i]
	})
	return out
}

// capCount truncates lines to at most n when a positive count caps the output.
func capCount(n ShufCount, lines [][]byte) [][]byte {
	if n > 0 && int(n) < len(lines) {
		return lines[:n]
	}
	return lines
}
