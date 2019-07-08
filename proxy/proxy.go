package proxy

import (
	"errors"

	dbjumper "github.com/stevelacy/dbjumper/pkg"
	"github.com/stevelacy/dbjumper/pooler"

	"net"
	"net/url"
)

// NetPair is a local and remote connection pair
type NetPair struct {
	// Address *net.TCPAddr
	Conn *net.Conn
}

// Connection is a network connection from the local proxy to a remote host
type Connection struct {
	local  NetPair
	remote NetPair
}

type closer func(net.Conn)

// Start a new proxy listener
func Start(config *dbjumper.Config) error {
	if len(config.Instances) == 0 {
		return errors.New("No instances are configured")
	}

	listener, err := net.Listen("tcp", config.Address)
	if err != nil {
		return err
	}

	for k, v := range config.Instances {
		// Assign the Address to an instance
		parsed, err := url.Parse(v.ConnectionString)
		if err != nil {
			dbjumper.Error(err)
			continue
		}
		v.Address = parsed.Host
		v.Name = k
		host, _ := net.ResolveTCPAddr("tcp", parsed.Host)
		c, err := net.DialTCP("tcp", nil, host)
		if err != nil {
			dbjumper.Log("instance offline: %s %s\n", k, parsed.Host)
			v.Online = false
			c.Close()
		} else {
			dbjumper.Log("instance connected: %s %s\n", k, parsed.Host)
			v.Online = true
			c.Close()
		}
		config.Instances[k] = v
	}

	for {
		lconn, err := listener.Accept()
		if err != nil {
			dbjumper.Error(err)
			continue
		}

		rconn, inst, err := pooler.Open(config)

		if err != nil {
			dbjumper.Error(err)
			continue
		}
		proxy := Connection{
			local: NetPair{
				// Address: host,
				Conn: &lconn,
			},
			remote: NetPair{
				// Address: instance.Address,
				Conn: &rconn,
			},
		}

		defer func() {
			lconn.Close()
			rconn.Close()
		}()

		go proxy.execute(lconn, rconn, func(conn net.Conn) {
			pooler.Free(config, inst, conn)
		})
	}
}

func (c *Connection) execute(local, remote net.Conn, end closer) {
	go pipe(local, remote, end)
	go pipe(remote, local, end)
}

// dummy pipe
// func pipe(src, dest net.Conn, end closer) {
// 	io.Copy(src, dest)
// 	io.Copy(dest, src)
// }

func pipe(src, dest net.Conn, end closer) {
	buff := make([]byte, 65535)
	for {
		n, err := src.Read(buff)
		if err != nil {
			if err.Error() == "EOF" {
				end(dest)
				return
			}
			dbjumper.Error(err)
			return
		}
		b := buff[:n]

		_, err = dest.Write(b)
		if err != nil {
			if err.Error() == "EOF" {
				end(dest)
				return
			}
			dbjumper.Error(err)
			return
		}
	}
}
