package pooler

import (
	"errors"
	"fmt"
	"log"
	"sort"

	// "github.com/jmoiron/sqlx"
	"net"
	// "net/url"
	// postgres
	// _ "github.com/lib/pq"
	"github.com/stevelacy/dbjumper/pkg"
)

// OpenA a connection to a host db
// func OpenA(config *dbjumper.Config) error {
//
// 	for key, inst := range config.Instances {
// 		var err error
// 		inst.Db, err = sqlx.Connect(inst.Type, inst.ConnectionString)
// 		if err != nil {
// 			log.Println(err)
// 			continue
// 		}
// 		parsed, err := url.Parse(inst.ConnectionString)
// 		if err != nil {
// 			log.Println(err)
// 			continue
// 		}
// 		inst.Address = parsed.Host
// 		inst.Db.SetMaxOpenConns(config.Pool.MaxServer)
// 		inst.Db.SetMaxIdleConns(config.Pool.MaxServerIdle)
// 		config.Instances[key] = inst
// 		log.Printf("connected to %s %s\n", key, inst.Address)
// 		defer inst.Db.Close()
// 	}
// 	return nil
// }

// Open a connection to a db
func Open(config *dbjumper.Config) (*net.TCPConn, dbjumper.Instance, error) {

	inst, err := getAvailableInstance(config)
	if err != nil {
		log.Println(err)
		return &net.TCPConn{}, inst, err
	}

	fmt.Printf("available connections: %d \n", len(inst.Connections))

	if len(inst.Connections) > 0 {
		// Lock RwMutext here
		conn := inst.Connections[0]
		fmt.Printf("using: %v\n", conn.LocalAddr())
		inst.Connections = inst.Connections[1:]
		config.Instances[inst.Name] = inst
		return conn, inst, nil
	}

	fmt.Printf("opened: %v\n", inst.ConnCount)

	addr, err := net.ResolveTCPAddr("tcp", inst.Address)
	if err != nil {
		return &net.TCPConn{}, inst, err
	}
	conn, err := net.DialTCP("tcp", nil, addr)
	inst.ConnCount++
	config.Instances[inst.Name] = inst

	if err != nil {
		return &net.TCPConn{}, inst, err
	}

	return conn, inst, nil
}

// Close frees a connection to make it available
func Close(config *dbjumper.Config, inst dbjumper.Instance, conn *net.TCPConn) {
	addr := conn.LocalAddr()
	// local connection closed, ignore
	if addr.String() != config.ListenAddress {
		return
	}
	fmt.Printf("closed: %s \n", addr)
	inst.Connections = append(inst.Connections, conn)
	inst.ConnCount--
	config.Instances[inst.Name] = inst
}

func getAvailableInstance(config *dbjumper.Config) (dbjumper.Instance, error) {
	if len(config.Instances) == 0 {
		return dbjumper.Instance{}, errors.New("0 instances configured")
	}

	var list []dbjumper.Instance
	for _, v := range config.Instances {
		list = append(list, v)
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].ConnCount < list[j].ConnCount
	})

	return list[0], nil
}
