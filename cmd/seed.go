package main

import (
	"github.com/urfave/cli/v2"

	"go-web-app/db/seeder"
)

var (
	cmdsSeeder = []*cli.Command{
		{
			Name:    "seed",
			Aliases: []string{},
			Usage:   "Do Seeders",
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name: FLAG_MIGRATION_SEED,
				},
			},
			Action: func(c *cli.Context) error {
				var err error

				err = InitDatabase(c, getFuncLoggerByCli(c))
				if err != nil {
					return err
				}

				err = seeder.GenerateSeederAutogenFile()
				if err != nil {
					return err
				}

				if c.Args().Len() > 0 {
					err = seeder.SeedBySeederName(c.Args().First())
					if err != nil {
						return err
					}
				}
				err = seeder.SeedAll()
				if err != nil {
					return err
				}

				return nil
			},
		},
	}
)
