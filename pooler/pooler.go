package pooler

import (
	"errors"
	"sort"

	"net"

	dbjumper "github.com/stevelacy/dbjumper/pkg"
)

// Open a connection to a db
func Open(config *dbjumper.Config) (net.Conn, dbjumper.Instance, error) {

	inst, err := getAvailableInstance(config)
	if err != nil {
		dbjumper.Error(err)
		return nil, inst, err
	}

	addr, err := net.ResolveTCPAddr("tcp", inst.Address)
	if err != nil {
		dbjumper.Log("error connecting: %v %s\n", inst.Address, inst.Name)
		inst.Online = false
		config.Instances[inst.Name] = inst
		return nil, inst, err
	}
	conn, err := net.DialTCP("tcp", nil, addr)

	dbjumper.Log("opened: %v %s\n", inst.ConnCount, inst.Name)

	inst.ConnCount++
	config.Instances[inst.Name] = inst

	if err != nil {
		return nil, inst, err
	}

	return conn, inst, nil
}

// Free a connection to make it available
func Free(config *dbjumper.Config, inst dbjumper.Instance, conn net.Conn) {
	addr := conn.LocalAddr()
	// local connection closed, ignore
	if addr.String() != config.Address {
		return
	}
	dbjumper.Log("closed: %s %s \n", addr, inst.Name)
	inst.ConnCount--
	config.Instances[inst.Name] = inst
	conn.Close()
}

func getAvailableInstance(config *dbjumper.Config) (dbjumper.Instance, error) {
	if len(config.Instances) == 0 {
		return dbjumper.Instance{}, errors.New("0 instances configured")
	}

	var list []dbjumper.Instance
	for _, v := range config.Instances {
		if v.Online {
			list = append(list, v)
		} else {
			// Check if it is now onlne
			addr, _ := net.ResolveTCPAddr("tcp", v.Address)
			c, err := net.DialTCP("tcp", nil, addr)
			if err == nil {
				v.Online = true
				list = append(list, v)
				c.Close()
			}

		}
	}
	if len(list) == 0 {
		return dbjumper.Instance{}, errors.New("0 instances online")
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].ConnCount < list[j].ConnCount
	})

	return list[0], nil
}
