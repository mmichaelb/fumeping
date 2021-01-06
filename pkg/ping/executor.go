package ping

import (
	"context"
	"github.com/go-ping/ping"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

type ResultHandler func(result Result)

type Executor struct {
	host          string
	context       context.Context
	waitGroup     *sync.WaitGroup
	resultHandler ResultHandler

	pinger *ping.Pinger
}

func New(host string, interval time.Duration, count int, timeout time.Duration, size int, context context.Context,
	waitGroup *sync.WaitGroup, resultHandler ResultHandler) (executor *Executor, err error) {
	executor = &Executor{
		host:          host,
		context:       context,
		waitGroup:     waitGroup,
		resultHandler: resultHandler,
	}
	executor.pinger, err = ping.NewPinger(host)
	if err != nil {
		return nil, err
	}
	executor.pinger.Interval = interval
	executor.pinger.Count = count
	executor.pinger.Timeout = timeout
	executor.pinger.Size = size
	return
}

func (executor *Executor) Run() {
	defer func() {
		executor.withHost().Debug("Stopped ping executor!")
		executor.waitGroup.Done()
	}()
	if !executor.setupPinger() {
		// pinger setup failed
		return
	}
	// run first ping sequence initially
	if !executor.RunPingSequence() {
		// ping sequence failed
		return
	}
	for {
		select {
		case <-executor.context.Done():
			executor.withHost().Debug("Stopping ping executor...")
			return
		case <-time.After(executor.pinger.Interval):
			// run pings after every interval
			if !executor.RunPingSequence() {
				// ping sequence failed
				return
			}
		}
	}
}

func (executor *Executor) setupPinger() bool {
	var err error
	executor.pinger, err = ping.NewPinger(executor.host)
	if err != nil {
		executor.withHost().WithError(err).Errorln("Could not initiate pinger")
		return false
	}
	return true
}

func (executor Executor) withHost() *logrus.Entry {
	return logrus.WithField("host", executor.host)
}
