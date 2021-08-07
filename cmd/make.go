package main

import (
	"fmt"
	"strings"

	"github.com/urfave/cli/v2"

	"go-web-app/common/util/timeutils"
	"go-web-app/db/migration"
	"go-web-app/db/seeder"
)

var (
	cmdsMake = []*cli.Command{
		{
			Name:    "make:model",
			Aliases: []string{"mM"},
			Usage:   "Create Model",
			Action: func(c *cli.Context) error {
				var err error
				if c.Args().Len() < 1 {
					fmt.Println("Model Name needed")
					return nil
				}
				actionName := "CreateModel"
				modelName := c.Args().First()
				datetime := timeutils.GetRFC3339StringForCodeGen(nil)
				datetimeRaw := timeutils.GetRFC3339String(nil)
				templateList := []TemplateReplaceDef{
					{
						templatePath: "asset/golangtemplate/Migration" + actionName + ".go.default",
						targetPath:   migration.GetMigrationFileName(datetime, actionName, modelName),
					},
					{
						templatePath: "asset/golangtemplate/Model.go.default",
						targetPath:   "common/model/" + modelName + ".go",
					},
				}
				err = ReplaceTemplate(templateList, func(content string) string {
					content = strings.Replace(content, "{{datetime}}", datetime, -1)
					content = strings.Replace(content, "{{datetimeRaw}}", datetimeRaw, -1)
					content = strings.Replace(content, "{{actionName}}", actionName, -1)
					content = strings.Replace(content, "{{modelName}}", modelName, -1)

					return content
				})
				if err != nil {
					return err
				}

				err = migration.GenerateMigrationAutogenFile()
				if err != nil {
					return err
				}

				return nil
			},
		},
		{
			Name:    "make:migration",
			Aliases: []string{"mMi"},
			Usage:   "Create Migration",
			Action: func(c *cli.Context) error {
				var err error
				if c.Args().Len() < 1 {
					fmt.Println("Migration Name needed")
					return nil
				}
				actionName := "General"
				migrationName := c.Args().First()
				datetime := timeutils.GetRFC3339StringForCodeGen(nil)
				datetimeRaw := timeutils.GetRFC3339String(nil)
				templateList := []TemplateReplaceDef{
					{
						templatePath: "asset/golangtemplate/Migration" + actionName + ".go.default",
						targetPath:   migration.GetMigrationFileName(datetime, actionName, migrationName),
					},
				}
				err = ReplaceTemplate(templateList, func(content string) string {
					content = strings.Replace(content, "{{datetime}}", datetime, -1)
					content = strings.Replace(content, "{{datetimeRaw}}", datetimeRaw, -1)
					content = strings.Replace(content, "{{actionName}}", actionName, -1)
					content = strings.Replace(content, "{{migrationName}}", migrationName, -1)

					return content
				})
				if err != nil {
					return err
				}

				err = migration.GenerateMigrationAutogenFile()
				if err != nil {
					return err
				}

				return nil
			},
		},
		{
			Name:    "make:seeder",
			Aliases: []string{"mSD"},
			Usage:   "Create Seeder",
			Action: func(c *cli.Context) error {
				var err error
				if c.Args().Len() < 1 {
					fmt.Println("Seeder Name needed")
					return nil
				}
				actionName := "General"
				seederName := c.Args().First()
				datetime := timeutils.GetRFC3339StringForCodeGen(nil)
				datetimeRaw := timeutils.GetRFC3339String(nil)
				templateList := []TemplateReplaceDef{
					{
						templatePath: "asset/golangtemplate/Seeder" + actionName + ".go.default",
						targetPath:   seeder.GetSeederFileName(datetime, actionName, seederName),
					},
				}
				err = ReplaceTemplate(templateList, func(content string) string {
					content = strings.Replace(content, "{{datetime}}", datetime, -1)
					content = strings.Replace(content, "{{datetimeRaw}}", datetimeRaw, -1)
					content = strings.Replace(content, "{{actionName}}", actionName, -1)
					content = strings.Replace(content, "{{seederName}}", seederName, -1)

					return content
				})
				if err != nil {
					return err
				}

				err = seeder.GenerateSeederAutogenFile()
				if err != nil {
					return err
				}

				return nil
			},
		},
	}
)
