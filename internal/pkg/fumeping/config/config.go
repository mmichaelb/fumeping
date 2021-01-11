package config

type Config struct {
	InfluxDb    InfluxDb
	PingMonitor PingMonitor
}

type PingMonitor struct {
	Timeout      int
	PingInterval int
	PayloadSize  uint16
	Destinations map[string]Destination
}

type Destination struct {
	Host    string
	Network string
}

type InfluxDb struct {
	DatabaseName   string
	ServerUrl      string
	AuthEnabled    bool
	Username       string
	Password       string
	GatherInterval int
}
