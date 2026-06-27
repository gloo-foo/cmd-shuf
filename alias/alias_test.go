package alias_test

import (
	"slices"
	"testing"

	shuf "github.com/gloo-foo/cmd-shuf/alias"
	"github.com/gloo-foo/testable"
)

// The alias package re-exports the constructor and flag types under unprefixed
// names. A mis-wired re-export (Count bound to Seed, Echo bound to Range, Shuf
// bound to the wrong function) compiles cleanly, so only behavior can prove the
// wiring. Each test exercises one re-export under a fixed seed and asserts the
// exact GNU shuf permutation it must produce.

func assertLines(t *testing.T, got, want []string) {
	t.Helper()
	if !slices.Equal(got, want) {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestAlias_ShufPermutesInput(t *testing.T) {
	// Shuf re-exports the constructor; seed 42 over a..e is [c d e a b].
	lines, err := testable.TestLines(shuf.Shuf(shuf.Seed(42)), "a\nb\nc\nd\ne\n")
	if err != nil {
		t.Fatal(err)
	}
	assertLines(t, lines, []string{"c", "d", "e", "a", "b"})
}

func TestAlias_SeedIsDeterministic(t *testing.T) {
	// Seed must bind to ShufSeed: same seed, identical permutation twice.
	first, err := testable.TestLines(shuf.Shuf(shuf.Seed(7)), "a\nb\nc\nd\n")
	if err != nil {
		t.Fatal(err)
	}
	second, err := testable.TestLines(shuf.Shuf(shuf.Seed(7)), "a\nb\nc\nd\n")
	if err != nil {
		t.Fatal(err)
	}
	assertLines(t, first, second)
}

func TestAlias_CountCapsOutput(t *testing.T) {
	// Count must bind to ShufCount (-n): the first 3 of [c d e a b].
	lines, err := testable.TestLines(shuf.Shuf(shuf.Count(3), shuf.Seed(42)), "a\nb\nc\nd\ne\n")
	if err != nil {
		t.Fatal(err)
	}
	assertLines(t, lines, []string{"c", "d", "e"})
}

func TestAlias_RangeShufflesIntegers(t *testing.T) {
	// Range must bind to ShufRange (-i): seed 42 over 1..5 is [3 4 5 1 2].
	lines, err := testable.TestLines(shuf.Shuf(shuf.Range(1, 5), shuf.Seed(42)), "")
	if err != nil {
		t.Fatal(err)
	}
	assertLines(t, lines, []string{"3", "4", "5", "1", "2"})
}

func TestAlias_EchoShufflesArgs(t *testing.T) {
	// Echo must bind to ShufEcho (-e): seed 42 over the three args.
	lines, err := testable.TestLines(shuf.Shuf(shuf.Echo("alpha", "beta", "gamma"), shuf.Seed(42)), "")
	if err != nil {
		t.Fatal(err)
	}
	assertLines(t, lines, []string{"gamma", "alpha", "beta"})
}
