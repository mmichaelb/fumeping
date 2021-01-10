package influx

import (
	"context"
	"fmt"
	ping2 "github.com/go-ping/ping"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/domain"
	"github.com/mmichaelb/fumeping/pkg/ping"
	"github.com/sirupsen/logrus"
)

const organization = "fumeping"
const bucketSyntax = "stats-%s"
const authTokenSyntax = "%s:%s"

type ResultHandler struct {
	influxClient influxdb2.Client
	organization *domain.Organization
}

func New(serverUrl, username, password string) (*ResultHandler, error) {
	client := influxdb2.NewClient(serverUrl, fmt.Sprintf(authTokenSyntax, username, password))
	handler := &ResultHandler{
		influxClient: client,
	}
	var err error
	organizationsApi := client.OrganizationsAPI()
	if handler.organization, err = organizationsApi.FindOrganizationByName(context.Background(), organization); err != nil {
		logrus.WithError(err).WithField("organizationName", organization).Errorln("Could not run InfluxDB find organization by name!")
	} else if handler.organization == nil {
		logrus.WithField("organizationName", organization).Debugln("Creating InfluxDB organization...")
		handler.organization, err = organizationsApi.CreateOrganizationWithName(context.Background(), organization)
		if err != nil {
			logrus.WithError(err).WithField("organizationName", organization).Errorln("Could not create InfluxDB organization!")
		}
		logrus.WithField("organizationName", organization).Debugln("Created InfluxDB organization!")
	}
	if err != nil {
		return nil, err
	}
	return handler, nil
}

func (handler *ResultHandler) Handle(host string, pinger ping2.Pinger, result ping.Result) {
	logrus.WithField("host", host).Debugln("Storing ping result in influxdb...")
	bucketName := fmt.Sprintf(bucketSyntax, host)
	bucketsApi := handler.influxClient.BucketsAPI()
	if bucket, err := bucketsApi.FindBucketByName(context.Background(), bucketName); err != nil {
		logrus.WithError(err).WithField("host", host).Errorln("Could not run InfluxDB find bucket by name!")
		return
	} else if bucket == nil {
		logrus.WithField("host", host).Infof("Creating bucket %s in influxdb...", bucketName)
		_, err := bucketsApi.CreateBucketWithName(context.Background(), handler.organization, bucketName)
		if err != nil {
			logrus.WithError(err).WithField("host", host).Errorln("Could not create new InfluxDB bucket!")
			return
		}
		logrus.WithField("host", host).Infof("Created new bucket %s in influxdb!", bucketName)
	}
	api := handler.influxClient.WriteAPIBlocking(organization, bucketName)
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
