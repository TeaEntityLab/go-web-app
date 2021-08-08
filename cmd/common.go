package main

import (
	"io/ioutil"

	"github.com/caarlos0/env/v6"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"go-web-app/common/util/fileutils"
	"go-web-app/db"
)

var (
	defaultLogger = logrus.New()
)

type TemplateReplaceDef struct {
	templatePath string
	targetPath   string
}

// ReplaceTemplate Replace the template file content variables.
func ReplaceTemplate(templateList []TemplateReplaceDef, replaceFunc func(string) string) error {
	var err error
	for _, item := range templateList {

		fileutils.Copy(item.templatePath, item.targetPath)
		var contentByte []byte
		contentByte, err = ioutil.ReadFile(item.targetPath)
		if err != nil {
			return err
		}

		content := string(contentByte)
		content = replaceFunc(content)
		err = ioutil.WriteFile(item.targetPath, []byte(content), 0)
		if err != nil {
			return err
		}
	}

	return nil
}

// InitDatabase Init the database as possible.
func InitDatabase(c *cli.Context, funcLogger *logrus.Entry) error {
	var envParseErr error
	cfg := db.CommonConfig{}
	envParseErr = env.Parse(&cfg)
	if envParseErr != nil {
		funcLogger.WithError(envParseErr).Fatalf("env.Parse error")
		return envParseErr
	}

	if dbErr := db.InitDefaultDatabase(cfg.DBType, cfg.DBEndpoints); dbErr != nil {
		funcLogger.WithError(dbErr).Fatalf("db.InitDefaultDatabase() error")
		return dbErr
	}

	return nil
}

func getFuncLoggerByCli(c *cli.Context) *logrus.Entry {
	return defaultLogger.WithField("func", c.Command.Name)
}

func concatCopyPreAllocate(slices [][]*cli.Command) []*cli.Command {
	var totalLen int
	for _, s := range slices {
		totalLen += len(s)
	}
	tmp := make([]*cli.Command, totalLen)
	var i int
	for _, s := range slices {
		i += copy(tmp[i:], s)
	}
	return tmp
}

func init() {
	defaultLogger.WithField("ServiceName", "Cli").Debugf("service start")
}
