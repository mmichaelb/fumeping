package config

import "time"

type Config struct {
	InfluxDb     InfluxDb
	Destinations []DestinationHost
}

type DestinationHost struct {
	Host           string
	Interval       time.Duration
	PacketInterval time.Duration
	Count          int
	Timeout        time.Duration
	PacketSize     int
}

type InfluxDb struct {
	Organization string
	BucketSyntax string
	ServerUrl    string
	AuthEnabled  bool
	Username     string
	Password     string
}
