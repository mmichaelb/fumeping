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

var LoadDefaultConfig = Config{
	InfluxDb: DefaultConfig.InfluxDb,
	PingMonitor: PingMonitor{
		Timeout:      10,
		PingInterval: 1,
		PayloadSize:  56,
		Destinations: map[string]Destination{},
	},
}

var DefaultConfig = Config{
	InfluxDb: InfluxDb{
		DatabaseName:   "fumeping",
		ServerUrl:      "http://localhost:8086/",
		AuthEnabled:    true,
		Username:       "admin",
		Password:       "mycrazypassword",
		GatherInterval: 30,
	},
	PingMonitor: PingMonitor{
		Timeout:      10,
		PingInterval: 1,
		PayloadSize:  56,
		Destinations: map[string]Destination{
			"First": {
				Host: "mycustomhost",
			},
			"Second": {
				Host:    "mycustomipv6host",
				Network: "ip6",
			},
		},
	},
}
