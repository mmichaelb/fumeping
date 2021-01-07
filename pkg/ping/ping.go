package ping

import (
	"net"
	"time"
)

type Result struct {
	PacketsRecv int
	PacketsSent int
	PacketLoss  float64
	IPAddr      *net.IPAddr
	Addr        string
	Rtts        []time.Duration
	MinRtt      time.Duration
	MaxRtt      time.Duration
	AvgRtt      time.Duration
	StdDevRtt   time.Duration
	Time        time.Time
}

func (executor *Executor) RunPingSequence() bool {
	executor.withHost().WithField("count", executor.pinger.Count).
		WithField("interval", executor.pinger.Interval).
		WithField("size", executor.pinger.Size).
		WithField("timeout", executor.pinger.Timeout).Debug("Pinging host")
	err := executor.pinger.Run()
	if err != nil {
		executor.withHost().WithError(err).Errorln("Could not run pinger")
		return false
	}
	stats := executor.pinger.Statistics()
	if executor.resultHandler != nil {
		result := Result{
			PacketsRecv: stats.PacketsRecv,
			PacketsSent: stats.PacketsSent,
			PacketLoss:  stats.PacketLoss,
			IPAddr:      stats.IPAddr,
			Addr:        stats.Addr,
			Rtts:        stats.Rtts,
			MinRtt:      stats.MinRtt,
			MaxRtt:      stats.MaxRtt,
			AvgRtt:      stats.AvgRtt,
			StdDevRtt:   stats.StdDevRtt,
			Time:        time.Now(),
		}
		go executor.resultHandler(executor.host, *executor.pinger, result)
	}
	return true
}
