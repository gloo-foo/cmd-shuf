package command

import (
	"math/rand"

	yup "github.com/gloo-foo/framework"
)

type command yup.Inputs[yup.File, flags]

func Shuf(parameters ...any) yup.Command {
	return command(yup.Initialize[yup.File, flags](parameters...))
}

func (p command) Executor() yup.CommandExecutor {
	return yup.AccumulateAndProcess(func(lines []string) []string {
		shuffled := make([]string, len(lines))
		copy(shuffled, lines)

		// Shuffle using Fisher-Yates algorithm
		rand.Shuffle(len(shuffled), func(i, j int) {
			shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
		})

		// Apply count limit if specified
		if p.Flags.Count > 0 && int(p.Flags.Count) < len(shuffled) {
			shuffled = shuffled[:p.Flags.Count]
		}

		return shuffled
	}).Executor()
}
