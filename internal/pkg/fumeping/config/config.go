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
	ServerUrl            string
	Organization         string
	DefaultBucket        string
	Username             string
	Password             string
	RetentionPeriodHours int
	GatherInterval       int
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
		ServerUrl:            "http://localhost:8086/",
		Organization:         "fumeping",
		DefaultBucket:        "fumeping",
		Username:             "admin",
		Password:             "mycrazypassword",
		RetentionPeriodHours: 24 * 7,
		GatherInterval:       30,
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
