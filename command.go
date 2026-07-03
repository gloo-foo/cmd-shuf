package command

import (
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
	var f flags
	for _, o := range opts {
		f = f.with(o)
	}
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
func echoLines(args []ShufLine) [][]byte {
	lines := make([][]byte, len(args))
	for i, a := range args {
		lines[i] = []byte(a)
	}
	return lines
}

// rangeLines renders the inclusive integer range lo..hi as input lines.
func rangeLines(r ShufRangeSpec) [][]byte {
	lines := make([][]byte, 0, r.hi-r.lo+1)
	for i := r.lo; i <= r.hi; i++ {
		lines = append(lines, []byte(strconv.Itoa(i)))
	}
	return lines
}

// capCount truncates lines to at most n when a positive count caps the output.
func capCount(n ShufCount, lines [][]byte) [][]byte {
	if n > 0 && int(n) < len(lines) {
		return lines[:n]
	}
	return lines
}
