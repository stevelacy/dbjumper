package dbjumper

// Instance is a replica or master host
type Instance struct {
	ConnectionString string
	Address          string
	Name             string
	Online           bool
	Type             string
	ConnCount        int
}

// Config for the db connections
type Config struct {
	Address        string // Listen address in host:port notation (127.0.0.1:5432)
	Instances      map[string]Instance
	ConnTimeout    int // TODO Time in seconds before the connection is released
	MaxServerConns int // TODO Max connections to each db
}
