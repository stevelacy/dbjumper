package cli

import (
	"github.com/stevelacy/dbjumper/pkg"
	"github.com/stevelacy/dbjumper/pooler"
	"github.com/stevelacy/dbjumper/proxy"
	"log"
	"net"
	"time"
)

// Init creates the cli
func Init() error {

	config := dbjumper.Config{
		ListenAddress: "127.0.0.1:6543",
		Instances:     map[string]dbjumper.Instance{},
	}
	config.Instances["node1"] = dbjumper.Instance{
		ConnectionString: "postgres://postgres@127.0.0.1:5432/stae?sslmode=disable",
		Master:           true,
		Type:             "postgres",
		Connections:      []net.Conn{},
	}
	config.Instances["node2"] = dbjumper.Instance{
		ConnectionString: "postgres://postgres@127.0.0.1:5432/stae?sslmode=disable",
		Master:           true,
		Type:             "postgres",
		Connections:      []net.Conn{},
	}

	log.Printf("starting on %s", config.ListenAddress)
	pooler.Open(&config)
	proxy.Start(&config)
	go forever()
	select {}
}

func forever() {
	for {
		time.Sleep(1)
	}
}
