package main

import (
	"crypto/tls"
	"time"
	"log"
	"github.com/urfave/cli"
	"github.com/marten-seemann/quic-conn"
	"net"
	"io"
	"os"
)

var (
	VERSION = "1.0"
)

func main() {
	myApp := cli.NewApp()
	myApp.Name = "quictun"
	myApp.Usage = "client of QUIC tunnel."
	myApp.Version = VERSION
	myApp.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "localaddr, l",
			Value: ":1935",
			Usage: "client listen address",
		},
		cli.StringFlag{
			Name:  "remoteaddr, r",
			Value: "127.0.0.1:6935",
			Usage: "quic server address",
		},
		cli.IntFlag{
			Name:  "timeout",
			Value: 5,
			Usage: "max time of waiting a connection to complete",
		},
		cli.IntFlag{
			Name:  "retry",
			Value: 10,
			Usage: "max retry time for quic server connect",
		},
		cli.BoolFlag{
			Name:  "quiet",
			Usage: "to suppress the 'stream open/close' messages",
		},
	}
	myApp.Action = func(c *cli.Context) error {
		config := Config{}
		config.LocalAddr = c.String("localaddr")
		config.RemoteAddr = c.String("remoteaddr")
		config.Timeout = c.Int("timeout")
		config.Retry = c.Int("retry")
		config.Quiet = c.Bool("quiet")

		log.SetFlags(log.LstdFlags | log.Lmicroseconds)

		// TODO: how to use TLS config?
		TLSConfig := func() *tls.Config {
			return &tls.Config{InsecureSkipVerify: true}
		}

		addr, err := net.ResolveTCPAddr("tcp", config.LocalAddr)
		if err != nil {
			log.Println(err)
			return err
		}

		listener, err := net.ListenTCP("tcp", addr)
		if err != nil {
			log.Println(err)
			return err
		}

		log.Println("version:", VERSION)
		log.Println("listening on:", listener.Addr())
		log.Println("remote addr:", config.RemoteAddr)
		log.Println("timeout:", config.Timeout)
		log.Println("retry:", config.Retry)
		log.Println("quiet:", config.Quiet)

		// transfer data between p1(tcp side) and p2(quic side).
		transfer := func(p1 io.ReadWriteCloser) {
			if !config.Quiet {
				log.Println("stream opened")
				defer log.Println("stream closed")
			}
			defer p1.Close()

			max := config.Retry
			p2, err := quicconn.Dial(config.RemoteAddr, TLSConfig())
			for err != nil {
				log.Println(err)
				if max <= 0 {
					p1.Close()
					return
				} else {
					time.Sleep(1 * time.Second)
					max--
					p2, err = quicconn.Dial(config.RemoteAddr, TLSConfig())
				}
			}
			defer p2.Close()

			p1die := make(chan struct{})
			go func() { io.Copy(p1, p2); close(p1die) }()

			p2die := make(chan struct{})
			go func() { io.Copy(p2, p1); close(p2die) }()

			select {
			case <-p1die:
			case <-p2die:
			}
		}

		for {
			if p1, err := listener.AcceptTCP(); err == nil {
				log.Printf("accpet tcp addr:%s\n", p1.RemoteAddr())

				go transfer(p1)
			} else {
				log.Fatalln(err)
			}
		}
	}
	myApp.Run(os.Args)
}
