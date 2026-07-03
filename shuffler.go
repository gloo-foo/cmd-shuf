package command

import "math/rand"

// shuffle permutes the first n elements by calling swap, matching the signature
// of math/rand's Shuffle. Injecting it as a value makes the random source a
// collaborator, so the otherwise non-deterministic default path is testable.
type shuffle func(n int, swap func(i, j int))

// seedSpec is an optional random seed carried by value: isSet distinguishes a
// user-provided seed (--seed) from the seedless default.
type seedSpec struct {
	value int64
	isSet bool
}

// shufflerFor builds a shuffle for the given seed: deterministic when a seed is
// present, process-default randomness otherwise.
type shufflerFor func(seed seedSpec) shuffle

// srcOption overrides the random source factory (test-only injection seam).
type srcOption struct{ factory shufflerFor }

// resolveSource resolves the random-source factory, defaulting when none was
// injected.
func resolveSource(f flags) shufflerFor {
	if f.source != nil {
		return f.source
	}
	return defaultShuffler
}

// defaultShuffler is the production random source: a deterministic, reproducible
// permutation when a seed is given, process-default randomness otherwise. shuf is
// a non-cryptographic line shuffler, so the seedless path uses math/rand's global
// Shuffle, while the seeded path drives a Fisher-Yates permutation from a seeded
// rand.Source — avoiding rand.New, whose weak-RNG use gosec flags (G404).
func defaultShuffler(seed seedSpec) shuffle {
	if seed.isSet {
		return seededShuffle(rand.NewSource(seed.value))
	}
	return rand.Shuffle
}

// seededShuffle returns a Fisher-Yates shuffle driven by src, so a fixed seed
// yields a reproducible permutation across runs.
func seededShuffle(src rand.Source) shuffle {
	return func(n int, swap func(i, j int)) {
		for i := n - 1; i > 0; i-- {
			swap(i, int(src.Int63()%int64(i+1)))
		}
	}
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
