package influx

import (
	"context"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/domain"
	"github.com/mmichaelb/fumeping/internal/pkg/fumeping/config"
	"github.com/mmichaelb/fumeping/internal/pkg/fumeping/ping"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

type ResultHandler struct {
	influxClient influxdb2.Client
	config       config.InfluxDb
	executor     *ping.Executor
	ctx          context.Context
}

func New(config config.InfluxDb, authToken, authTokenPath string, executor *ping.Executor, ctx context.Context) (*ResultHandler, error) {
	client := influxdb2.NewClient(config.ServerUrl, authToken)
	if authToken == "" {
		// InfluxDB seems to be fresh and not initialized yet
		authToken = setupInflux(client, config, authTokenPath)
	}
	handler := &ResultHandler{
		influxClient: client,
		config:       config,
		executor:     executor,
		ctx:          ctx,
	}
	return handler, nil
}

func setupInflux(client influxdb2.Client, influxConfig config.InfluxDb, authTokenPath string) string {
	retentionRule := int(time.Hour * time.Duration(influxConfig.RetentionPeriodHours))
	response, err := client.Setup(context.Background(), influxConfig.Username, influxConfig.Password, influxConfig.Organization, influxConfig.DefaultBucket, retentionRule)
	if err != nil {
		logrus.WithError(err).Fatalln("Could not setup InfluxDB!")
	}
	file, err := os.Create(authTokenPath)
	if err != nil {
		logrus.WithError(err).Fatalln("Could not create InfluxDB auth token file!")
	}
	token := *response.Auth.Token
	_, err = file.Write([]byte(token))
	if err != nil {
		logrus.WithError(err).Fatalln("Could not write InfluxDB auth token to file!")
	}
	return token
}

func (handler *ResultHandler) Run() {
	logrus.WithFields(logrus.Fields{"config": handler.config.Organization}).Infoln("Altering default retention policy...")
	bucketsApi := handler.influxClient.BucketsAPI()
	bucket, err := bucketsApi.FindBucketByName(context.Background(), handler.config.DefaultBucket)
	if err != nil {
		logrus.WithError(err).WithField("bucketName", handler.config.DefaultBucket).Fatalln("Could not find default bucket!")
	}
	bucket.RetentionRules = domain.RetentionRules{
		domain.RetentionRule{
			EverySeconds: handler.config.RetentionPeriodHours * 60 * 60,
		},
	}
	if _, err = bucketsApi.UpdateBucket(context.Background(), bucket); err != nil {
		logrus.WithError(err).WithField("bucketName", bucket.Name).Warnln("Could not update bucket retention rules!")
	}
	interval := time.Second * time.Duration(handler.config.GatherInterval)
	logrus.WithField("intervalSeconds", interval.String()).Infoln("Starting InfluxDB metrics gatherer...")
	for {
		select {
		case <-handler.ctx.Done():
			handler.Close()
			return
		case <-time.After(interval):
			handler.gather()
		}
	}
}

func (handler *ResultHandler) gather() {
	metrics := handler.executor.ExportAndClear()
	logrus.WithField("metrics", metrics).Debugln("Gathered metrics from ping monitor.")
	api := handler.influxClient.WriteAPIBlocking(handler.config.Organization, handler.config.DefaultBucket)
	for name, metric := range metrics {
		logrus.WithField("name", name).Debugln("Writing metrics entry to InfluxDB...")
		point := influxdb2.NewPointWithMeasurement("ping").
			AddTag("destination", name).
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
