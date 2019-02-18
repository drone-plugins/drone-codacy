package main

import (
	"log"
	"os"

	"github.com/urfave/cli"
)

var (
	version = "unknown"
)

func main() {
	app := cli.NewApp()
	app.Name = "codacy plugin"
	app.Usage = "codacy plugin"
	app.Action = run
	app.Version = version
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "token",
			Usage:  "token for authentication",
			EnvVar: "PLUGIN_TOKEN,CODACY_TOKEN",
		},
		cli.StringFlag{
			Name:   "pattern",
			Usage:  "coverage file pattern",
			Value:  "**/*.out",
			EnvVar: "PLUGIN_PATTERN",
		},
		cli.StringFlag{
			Name:   "language",
			Usage:  "language of coverage",
			Value:  "go",
			EnvVar: "PLUGIN_LANGUAGE",
		},
		cli.StringFlag{
			Name:   "commit.sha",
			Usage:  "git commit sha",
			EnvVar: "DRONE_COMMIT_SHA",
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(c *cli.Context) error {
	plugin := Plugin{
		Build: Build{
			Commit: c.String("commit.sha"),
		},
		Config: Config{
			Token:    c.String("token"),
			Pattern:  c.String("pattern"),
			Language: c.String("language"),
			Debug:    c.Bool("debug"),
		},
	}

	return plugin.Exec()
}
