package dbjumper

import (
	"github.com/jmoiron/sqlx"
	"net"
)

// Instance is a replica or master host
type Instance struct {
	ConnectionString string
	Address          string
	Master           bool
	Db               *sqlx.DB
	Type             string
	ConnCount        int
	Connections      []net.Conn
}

// Pool config
type Pool struct {
	// MaxClient     int
	// MaxClientIdle int
	MaxServer     int
	MaxServerIdle int
}

// Config for the db connections
type Config struct {
	ListenAddress string // Listen address in host:port notation (127.0.0.1:5432)
	Instances     map[string]Instance
	Pool
}
