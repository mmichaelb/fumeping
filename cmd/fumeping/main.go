package main

import (
	"context"
	"flag"
	"github.com/influxdata/influxdb-client-go/v2/log"
	fumePingConfig "github.com/mmichaelb/fumeping/internal/pkg/fumeping/config"
	logrus2 "github.com/mmichaelb/fumeping/internal/pkg/fumeping/logrus"
	"github.com/mmichaelb/fumeping/internal/pkg/fumeping/ping"
	"github.com/mmichaelb/fumeping/pkg/influx"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// application parameters
var logLevel = flag.String("level", "info",
	"Set the logging level. See https://github.com/sirupsen/logrus#level-logging for more details.")
var configPath = flag.String("config", "./config.yml", "Set the config file path.")

// ldflags
var GitVersion string
var GitBranch string
var GitDefaultBranch string
var config *fumePingConfig.Config
var executor *ping.Executor
var influxHandler *influx.ResultHandler

func main() {
	flag.Parse()
	setupLogrus()
	// startup message
	logrus.WithField("version", GitVersion).WithField("branch", GitBranch).Info("Starting FumePing...")
	loadConfig()
	// start ping executors
	startPingExecutors()
	// setup influx handler
	setupInfluxHandler()
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signalChannel
		logrus.Exit(0)
	}()
	// start blocking InfluxDB gatherer
	influxHandler.Run()
}

func startPingExecutors() {
	destinationCount := len(config.PingMonitor.Destinations)
	if destinationCount == 0 {
		logrus.WithField("destinationNumber", destinationCount).Fatalln("Could not start ping executors - no destination specified!")
	}
	logrus.WithField("destinationNumber", destinationCount).Infoln("Starting ping executors...")
	timeout := time.Second * time.Duration(config.PingMonitor.Timeout)
	interval := time.Second * time.Duration(config.PingMonitor.PingInterval)
	payloadSize := config.PingMonitor.PayloadSize
	var err error
	executor, err = ping.New(timeout, interval, payloadSize)
	if err != nil {
		logrus.WithField("timeout", timeout).
			WithField("interval", interval).WithField("payloadSize", payloadSize).
			WithError(err).Fatalln("Could not setup ping executor!")
	}
	logrus.DeferExitHandler(func() {
		logrus.Infoln("Stopping ping executor...")
		executor.Stop()
	})
	for name, destination := range config.PingMonitor.Destinations {
		if err := executor.AddHostTarget(name, destination.Network, destination.Host); err != nil {
			logrus.WithField("name", name).WithField("destination", destination).Fatalln("Could not add host target.")
		}
	}
	logrus.Infoln("Ping executors started in background!")
}

func setupInfluxHandler() {
	logrus.Infoln("Setting up InfluxDB metrics gatherer...")
	var err error
	serverUrl := config.InfluxDb.ServerUrl
	databaseName := config.InfluxDb.DatabaseName
	interval := time.Second * time.Duration(config.InfluxDb.GatherInterval)
	ctx, cancelFunc := context.WithCancel(context.Background())
	if config.InfluxDb.AuthEnabled {
		logrus.Debugln("Using auth enabled InfluxDB connection.")
		username := config.InfluxDb.Username
		password := config.InfluxDb.Password
		influxHandler, err = influx.NewWithAuth(serverUrl, databaseName, username, password, executor, interval, ctx)
	} else {
		logrus.Debugln("Using auth disabled InfluxDB connection.")
		influxHandler, err = influx.New(serverUrl, databaseName, executor, interval, ctx)
	}
	log.Log = &logrus2.WrappedLogrus{Logger: logrus.StandardLogger()}
	if err != nil {
		logrus.WithError(err).Fatalln("Could not instantiate new InfluxDB handler!")
	}
	logrus.DeferExitHandler(func() {
		logrus.Infoln("Stopping InfluxDB metrics gatherer...")
		cancelFunc()
		logrus.Infoln("InfluxDB metrics gatherer stopped!")
	})
	logrus.Infoln("Successfully set up InfluxDB metrics gatherer!")
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
