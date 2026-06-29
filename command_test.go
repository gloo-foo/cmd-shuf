package command_test

import (
	"testing"

	"github.com/gloo-foo/testable"
	"github.com/gloo-foo/testable/assertion"

	command "github.com/gloo-foo/cmd-shuf"
)

// These tests pin the exact permutation produced under a fixed seed. A fixed
// seed makes the shuffle deterministic, so the assertions verify the precise
// shuffle behavior — every input line is present, exactly once, in the order
// the seeded source dictates — not merely the output length or set membership.

// ==============================================================================
// Basic Shuffle (permute stdin lines)
// ==============================================================================

func TestShuf_PermutesAllInputLines(t *testing.T) {
	lines, err := testable.TestLines(command.Shuf(command.ShufSeed(42)), "a\nb\nc\nd\ne\n")
	assertion.NoError(t, err)
	assertion.Lines(t, lines, []string{"c", "b", "e", "d", "a"})
}

func TestShuf_SingleLine(t *testing.T) {
	lines, err := testable.TestLines(command.Shuf(command.ShufSeed(1)), "only\n")
	assertion.NoError(t, err)
	assertion.Lines(t, lines, []string{"only"})
}

// ==============================================================================
// Count Flag (-n): cap the output at n lines
// ==============================================================================

func TestShuf_CountCapsOutput(t *testing.T) {
	lines, err := testable.TestLines(command.Shuf(command.ShufCount(3), command.ShufSeed(42)), "a\nb\nc\nd\ne\n")
	assertion.NoError(t, err)
	// The first 3 lines of the seed-42 permutation [c b e d a].
	assertion.Lines(t, lines, []string{"c", "b", "e"})
}

func TestShuf_CountExceedingInputKeepsAll(t *testing.T) {
	lines, err := testable.TestLines(command.Shuf(command.ShufCount(100), command.ShufSeed(42)), "a\nb\nc\n")
	assertion.NoError(t, err)
	assertion.Count(t, lines, 3)
}

func TestShuf_CountOne(t *testing.T) {
	lines, err := testable.TestLines(command.Shuf(command.ShufCount(1), command.ShufSeed(42)), "a\nb\nc\nd\ne\n")
	assertion.NoError(t, err)
	assertion.Lines(t, lines, []string{"c"})
}

// ==============================================================================
// Seed Flag (--seed): deterministic, reproducible permutations
// ==============================================================================

func TestShuf_SameSeedSameOutput(t *testing.T) {
	input := "alpha\nbeta\ngamma\ndelta\nepsilon\n"

	lines1, err := testable.TestLines(command.Shuf(command.ShufSeed(99)), input)
	assertion.NoError(t, err)

	lines2, err := testable.TestLines(command.Shuf(command.ShufSeed(99)), input)
	assertion.NoError(t, err)

	assertion.Lines(t, lines1, lines2)
}

func TestShuf_DifferentSeedsDifferentOutput(t *testing.T) {
	input := "a\nb\nc\nd\ne\nf\ng\nh\ni\nj\n"

	lines1, err := testable.TestLines(command.Shuf(command.ShufSeed(1)), input)
	assertion.NoError(t, err)

	lines2, err := testable.TestLines(command.Shuf(command.ShufSeed(2)), input)
	assertion.NoError(t, err)

	same := true
	for i := range lines1 {
		if lines1[i] != lines2[i] {
			same = false
			break
		}
	}
	assertion.False(t, same, "different seeds should produce different output")
}

// ==============================================================================
// Empty Input
// ==============================================================================

func TestShuf_EmptyInput(t *testing.T) {
	lines, err := testable.TestLines(command.Shuf(), "")
	assertion.NoError(t, err)
	assertion.Empty(t, lines)
}

// ==============================================================================
// Range Flag (-i): shuffle the inclusive integer range, ignoring stdin
// ==============================================================================

func TestShuf_RangePermutesIntegers(t *testing.T) {
	lines, err := testable.TestLines(command.Shuf(command.ShufRange(1, 5), command.ShufSeed(42)), "")
	assertion.NoError(t, err)
	assertion.Lines(t, lines, []string{"3", "2", "5", "4", "1"})
}

func TestShuf_RangeIgnoresStdin(t *testing.T) {
	lines, err := testable.TestLines(command.Shuf(command.ShufRange(10, 12), command.ShufSeed(1)), "ignored\ninput\n")
	assertion.NoError(t, err)
	for _, l := range lines {
		assertion.True(t, l == "10" || l == "11" || l == "12", "range output must come from 10..12, not stdin")
	}
	assertion.Count(t, lines, 3)
}

func TestShuf_RangeWithCountCaps(t *testing.T) {
	lines, err := testable.TestLines(
		command.Shuf(command.ShufRange(1, 10), command.ShufCount(3), command.ShufSeed(42)),
		"",
	)
	assertion.NoError(t, err)
	assertion.Lines(t, lines, []string{"5", "3", "9"})
}

// ==============================================================================
// Echo Flag (-e): treat args as input lines, ignoring stdin
// ==============================================================================

func TestShuf_EchoPermutesArgs(t *testing.T) {
	lines, err := testable.TestLines(command.Shuf(command.ShufEcho("alpha", "beta", "gamma"), command.ShufSeed(42)), "")
	assertion.NoError(t, err)
	assertion.Lines(t, lines, []string{"alpha", "gamma", "beta"})
}

func TestShuf_EchoIgnoresStdin(t *testing.T) {
	lines, err := testable.TestLines(command.Shuf(command.ShufEcho("x", "y"), command.ShufSeed(1)), "ignored\n")
	assertion.NoError(t, err)
	for _, l := range lines {
		assertion.True(t, l == "x" || l == "y", "echo output must come from args, not stdin")
	}
	assertion.Count(t, lines, 2)
}

func TestShuf_EchoWithCountCaps(t *testing.T) {
	lines, err := testable.TestLines(
		command.Shuf(command.ShufEcho("a", "b", "c", "d"), command.ShufCount(2), command.ShufSeed(42)),
		"",
	)
	assertion.NoError(t, err)
	assertion.Lines(t, lines, []string{"c", "a"})
}

// echo takes precedence over -i, matching GNU shuf.
func TestShuf_EchoOverridesRange(t *testing.T) {
	lines, err := testable.TestLines(
		command.Shuf(command.ShufEcho("alpha", "beta", "gamma"), command.ShufRange(1, 100), command.ShufSeed(42)),
		"",
	)
	assertion.NoError(t, err)
	assertion.Lines(t, lines, []string{"alpha", "gamma", "beta"})
}
