package pooler

import (
	"github.com/jmoiron/sqlx"
	"net/url"
	// postgres
	_ "github.com/lib/pq"
	"github.com/stevelacy/dbjumper/pkg"
	"log"
	// "net"
)

// Open a connection to a host db
func Open(config *dbjumper.Config) error {

	for key, inst := range config.Instances {
		var err error
		inst.Db, err = sqlx.Connect(inst.Type, inst.ConnectionString)
		if err != nil {
			log.Println(err)
			continue
		}
		parsed, err := url.Parse(inst.ConnectionString)
		if err != nil {
			log.Println(err)
			continue
		}
		inst.Address = parsed.Host
		inst.Db.SetMaxOpenConns(config.Pool.MaxServer)
		inst.Db.SetMaxIdleConns(config.Pool.MaxServerIdle)
		config.Instances[key] = inst
		log.Printf("connected to %s %s\n", key, inst.Address)
		defer inst.Db.Close()
	}
	return nil
}
