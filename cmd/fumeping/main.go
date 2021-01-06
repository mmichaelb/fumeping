package main

import (
	"flag"
	"github.com/sirupsen/logrus"
)

// application parameters
var logLevel = flag.String("level", "info",
	"Set the logging level. See https://github.com/sirupsen/logrus#level-logging for more details.")

// ldflags
var GitVersion string
var GitBranch string
var GitDefaultBranch string

func main() {
	flag.Parse()
	// set logrus level
	level, err := logrus.ParseLevel(*logLevel)
	if err != nil {
		panic(err)
	}
	logrus.SetLevel(level)
	// set custom logging formatter
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceQuote:       true,
		QuoteEmptyFields: true,
	})
	// startup message
	logrus.WithField("version", GitVersion).WithField("branch", GitBranch).Info("Starting FumePing...")
	logrus.Info("Shutting down...")
	logrus.Info("Goodbye! Thank you for using FumePing.")
}
