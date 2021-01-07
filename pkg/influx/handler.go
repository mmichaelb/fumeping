package influx

import (
	"context"
	"fmt"
	ping2 "github.com/go-ping/ping"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/mmichaelb/fumeping/pkg/ping"
	"github.com/sirupsen/logrus"
)

const organization = "fumeping"
const bucketSyntax = "stats-%s"
const authTokenSyntax = "%s:%s"

type ResultHandler struct {
	influxClient influxdb2.Client
}

func New(serverUrl, username, password string) (*ResultHandler, error) {
	client := influxdb2.NewClient(serverUrl, fmt.Sprintf(authTokenSyntax, username, password))
	_, err := client.Health(context.Background())
	return &ResultHandler{
		influxClient: client,
	}, err
}

func (handler *ResultHandler) Handle(host string, pinger ping2.Pinger, result ping.Result) {
	logrus.WithField("host", host).Debugln("Storing ping result in influxdb...")
	api := handler.influxClient.WriteAPIBlocking(organization, fmt.Sprintf(bucketSyntax, host))
	point := influxdb2.NewPointWithMeasurement("ping").
		AddTag("host", host).
		AddField("interval", pinger.Interval).
		AddField("timeout", pinger.Timeout).
		AddField("count", pinger.Count).
		AddField("size", pinger.Size).
		AddField("packetsRecv", result.PacketsRecv).
		AddField("packetsSent", result.PacketsSent).
		AddField("packetLoss", result.PacketLoss).
		AddField("address", result.Addr).
		AddField("minRtt", result.MinRtt).
		AddField("maxRtt", result.MaxRtt).
		AddField("avgRtt", result.AvgRtt).
		AddField("stdDevRtt", result.StdDevRtt).
		SetTime(result.Time)
	if err := api.WritePoint(context.Background(), point); err != nil {
		logrus.WithField("host", host).WithError(err).Errorln("Could not write ping result point to InfluxDB!")
		return
	}
	logrus.WithField("host", host).Debugln("Successfully written ping result point to InfluxDB.")
}

func (handler *ResultHandler) Close() {
	handler.influxClient.Close()
}
