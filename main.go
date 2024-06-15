package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "hltools",
		Usage: "Various useful cli tools to use with hledger",
		Commands: []*cli.Command{
			{
				Name:    "depreciation",
				Aliases: []string{"depr"},
				Usage:   "calculate depreciations and return postings for each period",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "file",
						Usage:    "path of toml file with depreciation details",
						Aliases:  []string{"f"},
						Required: true,
					},
					&cli.BoolFlag{
						Name:  "thismonth",
						Usage: "Only return posting for this month",
					},
				},
				Action: Depreciations,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
