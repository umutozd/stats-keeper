package main

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/umutozd/stats-keeper/server"
	"github.com/urfave/cli/v2"
)

var config = server.NewConfig()

func main() {
	app := cli.NewApp()
	app.Name = "StatsKeeper Rest API"
	app.Flags = []cli.Flag{
		&cli.IntFlag{
			Name:        "http-port",
			Value:       config.HttpPort,
			Destination: &config.HttpPort,
			EnvVars:     []string{"SKEEPER_HTTP_PORT"},
			Usage:       "the port to listen for http requests",
		},
		&cli.StringFlag{
			Name:        "database-url",
			Destination: &config.DatabaseUrl,
			EnvVars:     []string{"SKEEPER_DATABASE_URL"},
			Usage:       "the url of the database server to connect to",
		},
	}
	app.Action = actionFunc

	if err := app.Run(os.Args); err != nil {
		logrus.WithError(err).Fatal("server exited with error")
	}
}

// actionFunc is the function called when cli app is run
func actionFunc(c *cli.Context) error {
	srv, err := server.NewServer(config)
	if err != nil {
		return err
	}
	return srv.ListenHTTP()
}
