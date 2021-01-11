package influx

import (
	"context"
	"fmt"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/mmichaelb/fumeping/internal/pkg/fumeping/ping"
	"github.com/sirupsen/logrus"
	"time"
)

const authTokenSyntax = "%s:%s"

type ResultHandler struct {
	influxClient influxdb2.Client
	databaseName string
	executor     *ping.Executor
	interval     time.Duration
	ctx          context.Context
}

func New(serverUrl, databaseName string, executor *ping.Executor, interval time.Duration, ctx context.Context) (*ResultHandler, error) {
	return NewWithAuth(serverUrl, databaseName, "", "", executor, interval, ctx)
}

func NewWithAuth(serverUrl, databaseName, username, password string, executor *ping.Executor, interval time.Duration, ctx context.Context) (*ResultHandler, error) {
	var authToken string
	if username == "" {
		authToken = fmt.Sprintf(authTokenSyntax, username, password)
	}
	client := influxdb2.NewClient(serverUrl, authToken)
	handler := &ResultHandler{
		influxClient: client,
		databaseName: databaseName,
		executor:     executor,
		interval:     interval,
		ctx:          ctx,
	}
	return handler, nil
}

func (handler *ResultHandler) Run() {
	logrus.WithField("interval", handler.interval.String()).Infoln("Starting InfluxDB metrics gatherer...")
	for {
		select {
		case <-handler.ctx.Done():
			handler.Close()
			return
		case <-time.After(handler.interval):
			handler.gather()
		}
	}
}

func (handler *ResultHandler) gather() {
	metrics := handler.executor.ExportAndClear()
	logrus.WithField("metrics", metrics).Debugln("Gathered metrics from ping monitor.")
	api := handler.influxClient.WriteAPIBlocking("", handler.databaseName)
	for name, metric := range metrics {
		logrus.WithField("name", name).Debugln("Writing metrics entry to InfluxDB...")
		point := influxdb2.NewPointWithMeasurement("ping").
			AddTag("name", name).
			AddField("packetsSent", metric.PacketsSent).
			AddField("packetsLost", metric.PacketsLost).
			AddField("minRtt", metric.Best).
			AddField("maxRtt", metric.Worst).
			AddField("avgRtt", metric.Mean).
			AddField("medianRtt", metric.Median).
			AddField("stdDevRtt", metric.StdDev).
			SetTime(time.Now())
		if err := api.WritePoint(context.Background(), point); err != nil {
			logrus.WithField("name", name).WithError(err).Errorln("Could not write ping result point to InfluxDB!")
			continue
		}
		logrus.WithField("name", name).Debugln("Successfully written metrics to InfluxDB.")
	}
	logrus.WithField("targetCount", len(metrics)).Infoln("Wrote ping monitor metrics to InfluxDB.")
}

func (handler *ResultHandler) Close() {
	handler.influxClient.Close()
}
