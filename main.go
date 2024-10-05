package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/fxyoge/hledger-merge/cmd"
)

func main() {
	app := &cli.App{
		Name:  "hledger-merge",
		Usage: "Merge hledger-compatible files",
		Flags: []cli.Flag{
			&cli.StringSliceFlag{
				Name:     "input",
				Aliases:  []string{"i"},
				Usage:    "Input hledger files",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "output",
				Aliases:  []string{"o"},
				Usage:    "Output file",
				Required: true,
			},
		},
		Action: cmd.Execute,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
