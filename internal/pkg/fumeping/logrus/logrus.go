package logrus

import "github.com/sirupsen/logrus"

// wrapped logrus logger in order to use logrus for influxdb v2 client
type WrappedLogrus struct {
	*logrus.Logger
}

func (wrappedLogrus *WrappedLogrus) SetPrefix(prefix string) {
	panic("not yet implemented")
}

func (wrappedLogrus *WrappedLogrus) SetLogLevel(logLevel uint) {
	wrappedLogrus.Logger.Level = logrus.Level(logLevel)
}

func (wrappedLogrus *WrappedLogrus) Debug(message string) {
	wrappedLogrus.Logger.Debugln(message)
}

func (wrappedLogrus *WrappedLogrus) Warn(message string) {
	wrappedLogrus.Logger.Warnln(message)
}

func (wrappedLogrus *WrappedLogrus) Error(message string) {
	wrappedLogrus.Logger.Errorln(message)
}

func (wrappedLogrus *WrappedLogrus) Info(message string) {
	wrappedLogrus.Logger.Infoln(message)
}
