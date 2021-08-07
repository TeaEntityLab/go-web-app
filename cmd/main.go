package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "wsgiCli"
	app.HelpName = "./cli.sh"
	app.EnableBashCompletion = true
	app.Commands = concatCopyPreAllocate([][]*cli.Command{{
		/*
			{
				Name:    "template",
				Aliases: []string{"t"},
				Usage:   "options for task templates",
				Subcommands: []*cli.Command{
					{
						Name:  "add",
						Usage: "add a new template",
						Action: func(c *cli.Context) error {
							fmt.Println("new task template: ", c.Args().First())
							return nil
						},
					},
					{
						Name:  "remove",
						Usage: "remove an existing template",
						Action: func(c *cli.Context) error {
							fmt.Println("removed task template: ", c.Args().First())
							return nil
						},
					},
				},
			},
		*/

	}, cmdsMake, cmdsSeeder, cmdsMigration})
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
