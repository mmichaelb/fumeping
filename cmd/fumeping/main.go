package main

import (
	"flag"
	fumePingConfig "github.com/mmichaelb/fumeping/internal/pkg/fumeping/config"
	"github.com/mmichaelb/fumeping/pkg/influx"
	"github.com/mmichaelb/fumeping/pkg/ping"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// application parameters
var logLevel = flag.String("level", "info",
	"Set the logging level. See https://github.com/sirupsen/logrus#level-logging for more details.")
var configPath = flag.String("config", "./config.toml", "Set the config file path.")

// ldflags
var GitVersion string
var GitBranch string
var GitDefaultBranch string
var config *fumePingConfig.Config
var influxHandler *influx.ResultHandler

func main() {
	flag.Parse()
	setupLogrus()
	// startup message
	logrus.WithField("version", GitVersion).WithField("branch", GitBranch).Info("Starting FumePing...")
	loadConfig()
	setupInfluxHandler()
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signalChannel
		logrus.Exit(0)
	}()
	// start ping executors
	startPingExecutors()
}

func startPingExecutors() {
	waitGroup := &sync.WaitGroup{}
	cancelContext, cancelFunction := context.WithCancel(context.Background())
	logrus.DeferExitHandler(func() {
		logrus.Infoln("Sending stop signal to ping executors...")
		cancelFunction()
	})
	logrus.WithField("destinationNumber", len(config.Destinations)).Infoln("Starting ping executors...")
	for _, destination := range config.Destinations {
		executor, err := ping.New(destination.Host, destination.Interval, destination.PacketInterval, destination.Count, destination.Timeout, destination.PacketSize, cancelContext, waitGroup, influxHandler.Handle)
		if err != nil {
			logrus.WithField("host", destination.Host).WithError(err).Fatalln("Could not setup executor!")
		}
		// increment waitgroup after pinger initialization was successful
		waitGroup.Add(1)
		go executor.Run()
	}
	logrus.Infoln("Ping executors started in background!")
	waitGroup.Wait()
	logrus.Infoln("Stopped ping executors!")
}

func setupInfluxHandler() {
	var err error
	if config.InfluxDb.AuthEnabled {
		influxHandler, err = influx.NewWithAuth(config.InfluxDb.ServerUrl, config.InfluxDb.DatabaseName, config.InfluxDb.Username, config.InfluxDb.Password)
	} else {
		influxHandler, err = influx.New(config.InfluxDb.ServerUrl, config.InfluxDb.DatabaseName)
	}
	logrus.DeferExitHandler(func() {
		logrus.Infoln("Closing InfluxDB connection...")
		influxHandler.Close()
		logrus.Infoln("InfluxDB connection closed!")
	})
	if err != nil {
		logrus.WithError(err).Fatalln("Could not instantiate new InfluxDB handler!")
	}
}

func setupLogrus() {
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
	logrus.RegisterExitHandler(func() {
		logrus.Info("Goodbye! Thank you for using FumePing.")
	})
}

func loadConfig() {
	logrus.WithField("configPath", *configPath).Infoln("Loading config...")
	var err error
	config, err = fumePingConfig.ReadConfig(*configPath)
	if err != nil {
		logrus.WithError(err).Fatalln("Could not load config!")
	}
	logrus.Infoln("Config loaded!")
}
