package ping

import (
	"context"
	"github.com/go-ping/ping"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

type ResultHandler func(host string, pinger ping.Pinger, result Result)

type Executor struct {
	host          string
	context       context.Context
	waitGroup     *sync.WaitGroup
	resultHandler ResultHandler
	interval      time.Duration

	pinger *ping.Pinger
}

func New(host string, interval, packetInterval time.Duration, count int, timeout time.Duration, size int, context context.Context,
	waitGroup *sync.WaitGroup, resultHandler ResultHandler) (executor *Executor, err error) {
	executor = &Executor{
		host:          host,
		context:       context,
		waitGroup:     waitGroup,
		resultHandler: resultHandler,
		interval:      interval,
	}
	executor.pinger, err = ping.NewPinger(host)
	if err != nil {
		return nil, err
	}
	executor.pinger.SetPrivileged(true)
	if packetInterval != 0 {
		executor.pinger.Interval = packetInterval
	}
	if count != 0 {
		executor.pinger.Count = count
	}
	if timeout != 0 {
		executor.pinger.Timeout = timeout
	}
	if size != 0 {
		executor.pinger.Size = size
	}
	return
}

func (executor *Executor) Run() {
	defer func() {
		executor.withHost().Debugln("Stopped ping executor!")
		executor.waitGroup.Done()
	}()
	// run first ping sequence initially
	if !executor.RunPingSequence() {
		// ping sequence failed
		return
	}
	for {
		select {
		case <-executor.context.Done():
			executor.withHost().Debugln("Stopping ping executor...")
			return
		case <-time.After(executor.interval):
			// run pings after every interval
			if !executor.RunPingSequence() {
				// ping sequence failed
				return
			}
		}
	}
}

func (executor Executor) withHost() *logrus.Entry {
	return logrus.WithField("host", executor.host)
}
