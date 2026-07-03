package shuf_test

import (
	"fmt"

	"github.com/gloo-foo/testable"

	command "github.com/gloo-foo/cmd-shuf"
)

func ExampleShuf_basic() {
	// echo "1\n2\n3\n4\n5" | shuf --seed 42 -n 2
	output, _ := testable.Test(
		command.Shuf(command.ShufCount(2), command.ShufSeed(42)),
		"1\n2\n3\n4\n5\n",
	)
	fmt.Print(output)
	// Output:
	// 3
	// 2
}
