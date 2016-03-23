package main

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"

	"github.com/itsyouonline/identityserver/db"
	"github.com/itsyouonline/identityserver/https"
	"github.com/itsyouonline/identityserver/routes"
)

func main() {

	app := cli.NewApp()
	app.Name = "Identity server"
	app.Version = "0.1-Dev"

	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})

	var debugLogging bool
	var bindAddress, dbConnectionString string
	var tlsCert, tlsKey string

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:        "debug, d",
			Usage:       "Enable debug logging",
			Destination: &debugLogging,
		},
		cli.StringFlag{
			Name:        "bind, b",
			Usage:       "Bind address",
			Value:       ":8080",
			Destination: &bindAddress,
		},
		cli.StringFlag{
			Name:        "connectionstring, c",
			Usage:       "Mongodb connection string",
			Value:       "127.0.0.1:27017",
			Destination: &dbConnectionString,
		},
		cli.StringFlag{
			Name:        "cert, s",
			Usage:       "TLS certificate path",
			Value:       "",
			Destination: &tlsCert,
		},
		cli.StringFlag{
			Name:        "key, k",
			Usage:       "TLS private key path",
			Value:       "",
			Destination: &tlsKey,
		},
	}

	app.Before = func(c *cli.Context) error {
		if debugLogging {
			log.SetLevel(log.DebugLevel)
			log.Debug("Debug logging enabled")
			log.Debug(app.Name, "-", app.Version)
		}
		return nil
	}

	app.Action = func(c *cli.Context) {
		// Connect to DB!
		go db.Connect(dbConnectionString)
		defer db.Close()

		// Get router!
		r := routes.GetRouter()

		// Go make magic over HTTPS
		log.Fatal(https.ListenAndServeTLS(bindAddress, tlsCert, tlsKey, r))
	}

	app.Run(os.Args)
}
