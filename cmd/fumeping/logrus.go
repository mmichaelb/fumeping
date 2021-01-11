package main

import "github.com/sirupsen/logrus"

// wrapped logrus logger in order to use logrus for influxdb v2 client
type wrappedLogrus struct {
	*logrus.Logger
}

func (wrappedLogrus *wrappedLogrus) SetPrefix(prefix string) {
	panic("not yet implemented")
}

func (wrappedLogrus *wrappedLogrus) SetLogLevel(logLevel uint) {
	wrappedLogrus.Logger.Level = logrus.Level(logLevel)
}

func (wrappedLogrus *wrappedLogrus) Debug(message string) {
	wrappedLogrus.Logger.Debugln(message)
}

func (wrappedLogrus *wrappedLogrus) Warn(message string) {
	wrappedLogrus.Logger.Warnln(message)
}

func (wrappedLogrus *wrappedLogrus) Error(message string) {
	wrappedLogrus.Logger.Errorln(message)
}

func (wrappedLogrus *wrappedLogrus) Info(message string) {
	wrappedLogrus.Logger.Infoln(message)
}
