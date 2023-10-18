package cli

import (
	"os"
	"testing"

	"gotest.tools/v3/assert"
)

func TestVersion(t *testing.T) {
	expected := "version 1, commit unknown, built at today\n"
	icli := geticliForTesting(os.DirFS("../.."))

	reseticliBuffer(&icli)
	icli.Run([]string{"initium", "version"})
	output := icliOutput(icli)
	assert.Assert(t, output == expected, "Expected %s, got %s", expected, output)
}
