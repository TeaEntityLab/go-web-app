package main

import (
	"strconv"

	"github.com/urfave/cli/v2"

	"go-web-app/db/migration"
	"go-web-app/db/seeder"
)

const (
	FLAG_MIGRATION_SEED = "seed"
)

var (
	cmdsMigration = []*cli.Command{
		{
			Name:    "migrate",
			Aliases: []string{},
			Usage:   "Do Migrations",
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

				err = migration.GenerateMigrationAutogenFile()
				if err != nil {
					return err
				}
				migration.Migrate()

				if c.Bool(FLAG_MIGRATION_SEED) {
					seeder.GenerateSeederAutogenFile()
					err = seeder.SeedAll()
					if err != nil {
						return err
					}
				}

				return nil
			},
		},
		{
			Name:    "migrate:rollback",
			Aliases: []string{},
			Usage:   "Do the rollbacks of the Migrations. migrate:rollback [steps]",
			Action: func(c *cli.Context) error {
				var err error

				err = InitDatabase(c, getFuncLoggerByCli(c))
				if err != nil {
					return err
				}

				steps := 1
				if c.Args().Len() > 0 {
					stepsStr := c.Args().First()
					steps, _ = strconv.Atoi(stepsStr)
					if steps < 1 {
						steps = 1
					}
				}

				err = migration.GenerateMigrationAutogenFile()
				if err != nil {
					return err
				}
				migration.MigrateRollback(steps)

				return nil
			},
		},
		{
			Name:    "migrate:refresh",
			Aliases: []string{},
			Usage:   "Do all the rollbacks and the Migrations. migrate:refresh [--seed]",
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

				err = migration.GenerateMigrationAutogenFile()
				if err != nil {
					return err
				}
				migration.MigrateRollback(99999)
				migration.Migrate()
				//app.Command("migrate").Run(c)

				if c.Bool(FLAG_MIGRATION_SEED) {
					seeder.GenerateSeederAutogenFile()
					err = seeder.SeedAll()
					if err != nil {
						return err
					}
				}

				return nil
			},
		},
	}
)
