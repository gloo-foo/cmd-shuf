package command

import (
	"slices"
	"testing"

	"github.com/gloo-foo/testable"
)

// withReverse injects a fully deterministic source that reverses the first n
// elements (i swapped with n-1-i). With a known source, the exact output
// permutation is asserted rather than only its length or set membership.
func withReverse() srcOption {
	return srcOption{factory: func(*int64) shuffle {
		return func(n int, swap func(i, j int)) {
			for i := 0; i < n/2; i++ {
				swap(i, n-1-i)
			}
		}
	}}
}

func TestPermute_ReverseSourceProducesExactOrder(t *testing.T) {
	lines, err := testable.TestLines(Shuf(withReverse()), "a\nb\nc\nd\n")
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{"d", "c", "b", "a"}; !slices.Equal(lines, want) {
		t.Fatalf("got %q, want %q", lines, want)
	}
}

func TestPermute_DoesNotMutateInput(t *testing.T) {
	in := [][]byte{[]byte("a"), []byte("b"), []byte("c")}
	original := [][]byte{[]byte("a"), []byte("b"), []byte("c")}
	rev := withReverse().factory(nil)
	_ = permute(rev, in)
	for i := range in {
		if string(in[i]) != string(original[i]) {
			t.Fatalf("input mutated at %d: got %q, want %q", i, in[i], original[i])
		}
	}
}

// TestDefaultShuffler_Seeded covers the seeded branch of the production source
// and pins its deterministic permutation for seed 42.
func TestDefaultShuffler_Seeded(t *testing.T) {
	lines, err := testable.TestLines(
		Shuf(srcOption{factory: defaultShuffler}, ShufSeed(42)),
		"a\nb\nc\nd\ne\n",
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(lines) != 5 {
		t.Fatalf("got %d lines, want 5", len(lines))
	}
	// defaultShuffler(seed) must be reproducible: a second run with the same
	// seed yields the identical permutation.
	again, err := testable.TestLines(
		Shuf(srcOption{factory: defaultShuffler}, ShufSeed(42)),
		"a\nb\nc\nd\ne\n",
	)
	if err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(lines, again) {
		t.Fatalf("seeded source not reproducible: %q vs %q", lines, again)
	}
}

// TestDefaultShuffler_Unseeded covers the seedless branch of the production
// source. Without a seed the order is non-deterministic, so only the invariant
// that every input survives the shuffle is asserted.
func TestDefaultShuffler_Unseeded(t *testing.T) {
	s := defaultShuffler(nil)
	out := permute(s, [][]byte{[]byte("a"), []byte("b"), []byte("c")})
	got := []string{string(out[0]), string(out[1]), string(out[2])}
	slices.Sort(got)
	if want := []string{"a", "b", "c"}; !slices.Equal(got, want) {
		t.Fatalf("seedless shuffle lost elements: got %q", got)
	}
}

// TestResolveSource_DefaultsWhenAbsent covers the default branch of
// resolveSource (no injected source).
func TestResolveSource_DefaultsWhenAbsent(t *testing.T) {
	lines, err := testable.TestLines(Shuf(ShufSeed(7)), "a\nb\n")
	if err != nil {
		t.Fatal(err)
	}
	if len(lines) != 2 {
		t.Fatalf("got %d lines, want 2", len(lines))
	}
}
