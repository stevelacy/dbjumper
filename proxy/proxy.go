package proxy

import (
	"errors"
	"github.com/stevelacy/dbjumper/pkg"
	"io"
	"log"
	"net"
	"sort"
)

// NetPair is a local and remote connection pair
type NetPair struct {
	Address *net.TCPAddr
	Conn    *net.TCPConn
}

// Connection is a network connection from the local proxy to a remote host
type Connection struct {
	local  NetPair
	remote NetPair
}

// Start a new proxy listener
func Start(config *dbjumper.Config) (net.Listener, error) {

	host, err := net.ResolveTCPAddr("tcp", config.ListenAddress)
	if err != nil {
		log.Println(err)
	}
	listener, err := net.ListenTCP("tcp", host)

	remote, err := getRemote(config)
	if err != nil {
		log.Println(err)
		return listener, err
	}

	if err != nil {
		return listener, err
	}
	for {
		lconn, err := listener.AcceptTCP()
		if err != nil {
			log.Println(err)
			continue
		}

		rconn, err := net.DialTCP("tcp", nil, remote.Address)
		if err != nil {
			log.Println(err)
			continue
		}

		if err != nil {
			log.Println(err)
			continue
		}
		proxy := Connection{
			local: NetPair{
				Address: host,
				Conn:    lconn,
			},
			remote: NetPair{
				Address: remote.Address,
				Conn:    rconn,
			},
		}

		defer func() {
			proxy.local.Conn.Close()
			proxy.remote.Conn.Close()
		}()

		go proxy.execute(lconn, rconn)
	}
}

func (c *Connection) execute(local, remote *net.TCPConn) {
	go pipe(local, remote)
	go pipe(remote, local)
}

func pipe(src, dest *net.TCPConn) {
	io.Copy(src, dest)
	io.Copy(dest, src)

	// for {
	// 	val, err := src.Read(buff)
	//
	// 	if err != nil {
	// 		if err.Error() != "EOF" {
	// 			log.Println(err)
	// 		}
	// 		return
	// 	}
	//
	// 	b := buff[:val]
	// 	_, err = dest.Write(b)
	// 	if err != nil {
	// 		if err.Error() != "EOF" {
	// 			log.Println(err)
	// 		}
	// 		return
	// 	}
	// }
}

func getRemote(config *dbjumper.Config) (NetPair, error) {
	if len(config.Instances) == 0 {
		return NetPair{}, errors.New("0 instances configured")
	}

	var list []dbjumper.Instance
	for _, v := range config.Instances {
		list = append(list, v)
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].ConnCount < list[j].ConnCount
	})

	addr, err := net.ResolveTCPAddr("tcp", list[0].Address)
	if err != nil {
		return NetPair{}, err
	}

	return NetPair{Address: addr}, nil
}
