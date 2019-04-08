package dbjumper

import (
	// "github.com/jmoiron/sqlx"
	"net"
)

// Instance is a replica or master host
type Instance struct {
	ConnectionString string
	Address          string
	Master           bool
	Name             string
	Type             string
	ConnCount        int
	Connections      []*net.TCPConn // Available connections for use
	// Add a mutex here
}

// Config for the db connections
type Config struct {
	ListenAddress  string // Listen address in host:port notation (127.0.0.1:5432)
	Instances      map[string]Instance
	ConnTimeout    int // Time in seconds before the connection is released
	MaxServerConns int // Max connections to each db
	// TODO: MaxClientConns int // Max client connections to dbjumper
}
