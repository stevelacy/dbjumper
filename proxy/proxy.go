package proxy

import (
	"github.com/stevelacy/dbjumper/pkg"
	"github.com/stevelacy/dbjumper/pooler"
	// "io"
	"errors"
	"log"
	"net"
	"net/url"
)

// NetPair is a local and remote connection pair
type NetPair struct {
	// Address *net.TCPAddr
	Conn *net.TCPConn
}

// Connection is a network connection from the local proxy to a remote host
type Connection struct {
	local  NetPair
	remote NetPair
}

type closer func(*net.TCPConn)

// Start a new proxy listener
func Start(config *dbjumper.Config) error {
	if len(config.Instances) == 0 {
		return errors.New("No instances are configured")
	}

	host, err := net.ResolveTCPAddr("tcp", config.ListenAddress)
	if err != nil {
		log.Println(err)
	}
	listener, err := net.ListenTCP("tcp", host)
	if err != nil {
		return err
	}

	for k, v := range config.Instances {
		// Assign the Address to an instance
		parsed, err := url.Parse(v.ConnectionString)
		if err != nil {
			log.Println(err)
			continue
		}
		v.Address = parsed.Host
		v.Name = k
		log.Printf("instance connected: %s %s\n", k, parsed.Host)
		config.Instances[k] = v
	}

	for {
		lconn, err := listener.AcceptTCP()
		if err != nil {
			log.Println(err)
			continue
		}

		rconn, inst, err := pooler.Open(config)

		if err != nil {
			log.Println(err)
			continue
		}
		proxy := Connection{
			local: NetPair{
				// Address: host,
				Conn: lconn,
			},
			remote: NetPair{
				// Address: instance.Address,
				Conn: rconn,
			},
		}

		defer func() {
			proxy.local.Conn.Close()
			proxy.remote.Conn.Close()
		}()

		go proxy.execute(lconn, rconn, func(conn *net.TCPConn) {
			pooler.Close(config, inst, conn)
		})
	}
}

func (c *Connection) execute(local, remote *net.TCPConn, end closer) {
	go pipe(local, remote, end)
	go pipe(remote, local, end)
}

// func pipe(src, dest *net.TCPConn) {
// 	io.Copy(src, dest)
// 	io.Copy(dest, src)
// }

func pipe(src, dest *net.TCPConn, end closer) {
	buff := make([]byte, 65535)
	for {
		n, err := src.Read(buff)
		if err != nil {
			if err.Error() == "EOF" {
				end(dest)
				return
			}
			log.Println(err)
			return
		}
		b := buff[:n]

		_, err = dest.Write(b)
		if err != nil {
			if err.Error() == "EOF" {
				// end(dest)
				return
			}
			log.Println(err)
			return
		}
	}
}
