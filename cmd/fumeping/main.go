package main

import (
	"context"
	"flag"
	"github.com/influxdata/influxdb-client-go/v2/log"
	fumePingConfig "github.com/mmichaelb/fumeping/internal/pkg/fumeping/config"
	"github.com/mmichaelb/fumeping/internal/pkg/fumeping/influx"
	logrus2 "github.com/mmichaelb/fumeping/internal/pkg/fumeping/logrus"
	"github.com/mmichaelb/fumeping/internal/pkg/fumeping/ping"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

// application parameters
var logLevel = flag.String("level", "info",
	"Set the logging level. See https://github.com/sirupsen/logrus#level-logging for more details.")
var configPath = flag.String("config", "./config.yml", "Set the config file path.")
var influxAuthTokenFile = flag.String("tokenfile", ".influxdb_token", "Set the InfluxDB token file to use.")

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
	token, err := loadInfluxToken()
	if err != nil {
		if os.IsNotExist(err) {
			token = ""
		} else {
			logrus.WithError(err).Fatalln("Could not load InfluxDB token!")
		}
	}
	ctx, cancelFunc := context.WithCancel(context.Background())
	influxHandler, err = influx.New(config.InfluxDb, token, *influxAuthTokenFile, executor, ctx)
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

func loadInfluxToken() (string, error) {
	tokenBytes, err := ioutil.ReadFile(*influxAuthTokenFile)
	if err != nil {
		return "", err
	}
	token := strings.Trim(string(tokenBytes), " \r\n")
	return token, nil
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
