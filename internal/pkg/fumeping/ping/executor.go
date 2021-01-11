package ping

import (
	"github.com/digineo/go-ping"
	"github.com/digineo/go-ping/monitor"
	"github.com/sirupsen/logrus"
	"net"
	"time"
)

const defaultNetwork = "ip"

type Executor struct {
	*monitor.Monitor
}

func New(timeout time.Duration, interval time.Duration, size uint16) (executor *Executor, err error) {
	// update logger for ping framework
	ping.SetLogger(logrus.StandardLogger())
	var bind4, bind6 string
	if ln, err := net.Listen("tcp4", "127.0.0.1:0"); err == nil {
		ln.Close()
		bind4 = "0.0.0.0"
	}
	if ln, err := net.Listen("tcp6", "[::1]:0"); err == nil {
		ln.Close()
		bind6 = "::"
	}
	pinger, err := ping.New(bind4, bind6)
	if err != nil {
		return nil, err
	}
	if size != 0 {
		pinger.SetPayloadSize(size)
	} else {
		pinger.SetPayloadSize(24)
	}
	executor = &Executor{
		Monitor: monitor.New(pinger, interval, timeout),
	}
	return
}

func (executor *Executor) AddHostTarget(key, network, host string) error {
	if network == "" {
		network = defaultNetwork
	}
	ip, err := net.ResolveIPAddr(network, host)
	if err != nil {
		return err
	}
	if err = executor.Monitor.AddTarget(key, *ip); err != nil {
		return err
	}
	return nil
}
