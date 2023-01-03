package main

import (
	"fmt"
	"github.com/cocov-ci/coverage-reporter/command"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

func envs(base ...string) []string {
	return base
}

func main() {
	app := cli.NewApp()
	app.Name = "coverage-reporter"
	app.Usage = "Emits coverage reports to a Cocov instance"
	app.Version = "0.1"
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:     "token",
			Usage:    "Defines the token identifying this repository",
			Required: true,
			EnvVars:  envs("COCOV_REPOSITORY_TOKEN"),
		},
		&cli.BoolFlag{
			Name:    "debug",
			Usage:   "Enables debug logging",
			EnvVars: envs("COCOV_REPORTER_DEBUG"),
		},
	}
	app.Authors = []*cli.Author{
		{Name: "Victor \"Vito\" Gama", Email: "hey@vito.io"},
	}
	app.Before = func(ctx *cli.Context) error {
		var config zap.Config
		if ctx.Bool("debug") {
			config = zap.NewDevelopmentConfig()
			if os.Getenv("COCOV_REPORTER_DEV") == "true" {
				config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
			}
			config.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
		} else {
			config = zap.NewProductionConfig()
			config.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
		}
		logger, err := config.Build()
		if err != nil {
			return err
		}

		zap.ReplaceGlobals(logger)
		return nil
	}
	app.Copyright = "Copyright (c) 2022-2023 - The Cocov Authors"
	app.Commands = []*cli.Command{
		{
			Name:   "prepare",
			Usage:  "Starts a worker",
			Action: command.Prepare,
		},
		{
			Name:  "submit",
			Usage: "Submits recorded coverage results",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "url",
					Usage:    "Defines the Cocov API URL that will receive the report",
					Required: true,
					EnvVars:  envs("COCOV_API_URL"),
				},
				&cli.StringFlag{
					Name:     "commitish",
					Usage:    "The SHA identifying the commit that generated the report being sent",
					Required: false,
					EnvVars:  envs("GITHUB_SHA", "GIT_SHA", "COMMMIT_SHA", "COCOV_COMMIT_SHA"),
				},
			},
			Action: command.Submit,
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println()
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
}
