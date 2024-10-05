package cmd

import (
	"github.com/fxyoge/hledger-merge/internal/merger"
	"github.com/urfave/cli/v2"
)

func Execute(c *cli.Context) error {
	inputs := c.StringSlice("input")
	output := c.String("output")

	return merger.MergeFiles(inputs, output)
}
