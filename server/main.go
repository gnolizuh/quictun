package main

import (
	"crypto/rand"
	"crypto/tls"
	"crypto/rsa"
	"crypto/x509"
	"math/big"
	"time"
	"log"
	"encoding/pem"
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
	myApp.Usage = "server of QUIC tunnel."
	myApp.Version = VERSION
	myApp.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "listen, l",
			Value: ":6935",
			Usage: "server listen address",
		},
		cli.StringFlag{
			Name:  "target, t",
			Value: "127.0.0.1:1935",
			Usage: "target server address",
		},
		cli.IntFlag{
			Name:  "timeout",
			Value: 5,
			Usage: "max time of waiting a connection to complete",
		},
		cli.IntFlag{
			Name:  "retry",
			Value: 10,
			Usage: "max retry time for target connect",
		},
		cli.BoolFlag{
			Name:  "quiet",
			Usage: "to suppress the 'stream open/close' messages",
		},
	}
	myApp.Action = func(c *cli.Context) error {
		config := Config{}
		config.Listen = c.String("listen")
		config.Target = c.String("target")
		config.Timeout = c.Int("timeout")
		config.Retry = c.Int("retry")
		config.Quiet = c.Bool("quiet")

		log.SetFlags(log.LstdFlags | log.Lmicroseconds)

		// TODO: how to use TLS config?
		TLSConfig := func() *tls.Config {
			key, err := rsa.GenerateKey(rand.Reader, 2048)
			if err != nil {
				log.Println(err)
				return nil
			}

			template := x509.Certificate{
				SerialNumber: big.NewInt(1),
				NotBefore:    time.Now(),
				NotAfter:     time.Now().Add(time.Hour),
				KeyUsage:     x509.KeyUsageCertSign | x509.KeyUsageDigitalSignature,
			}

			certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
			if err != nil {
				log.Println(err)
				return nil
			}

			keyPEM := pem.EncodeToMemory(&pem.Block{
				Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key),
			})
			b := pem.Block{Type: "CERTIFICATE", Bytes: certDER}
			certPEM := pem.EncodeToMemory(&b)

			tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
			if err != nil {
				log.Println(err)
				return nil
			}

			return &tls.Config{ Certificates: []tls.Certificate{tlsCert} }
		}

		listener, err := quicconn.Listen("udp", config.Listen, TLSConfig())
		if err != nil {
			panic(err)
		}

		log.Println("version:", VERSION)
		log.Println("listening on:", listener.Addr())
		log.Println("target:", config.Target)
		log.Println("timeout:", config.Timeout)
		log.Println("retry:", config.Retry)
		log.Println("quiet:", config.Quiet)

		// transfer data between p1(quic side) and p2(tcp side).
		transfer := func(p1 io.ReadWriteCloser) {
			if !config.Quiet {
				log.Println("stream opened")
				defer log.Println("stream closed")
			}
			defer p1.Close()

			max := config.Retry
			p2, err := net.DialTimeout("tcp", config.Target, time.Duration(config.Timeout) * time.Second)
			for err != nil {
				log.Println(err)
				if max <= 0 {
					p1.Close()
					return
				} else {
					time.Sleep(1 * time.Second)
					max--
					p2, err = net.DialTimeout("tcp", config.Target, time.Duration(config.Timeout) * time.Second)
				}
			}
			defer p2.Close()

			p1die := make(chan struct{})
			go func() {
				n, err := io.Copy(p1, p2)
				if err != nil {
					log.Println(err)
				}

				log.Printf("<- wrie %d bytes", n)

				close(p1die)
			}()

			p2die := make(chan struct{})
			go func() {
				n, err := io.Copy(p2, p1)
				if err != nil {
					log.Println(err)
				}

				log.Printf("-> write %d bytes", n)

				close(p2die)
			}()

			select {
			case <-p1die:
			case <-p2die:
			}
		}

		for {
			if p1, err := listener.Accept(); err == nil {
				log.Printf("accpet quic addr:%s\n", p1.RemoteAddr())

				go transfer(p1)
			} else {
				log.Fatalln(err)
			}
		}
	}
	myApp.Run(os.Args)
}
