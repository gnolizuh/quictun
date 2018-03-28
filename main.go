package main

import (
	"crypto/rand"
	"crypto/tls"
	"crypto/rsa"
	"crypto/x509"
	"math/big"
	"time"
	"log"
	"fmt"
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

// TODO: how to use TLS config?
func TLSConfig() (*tls.Config) {
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

func handleClient(p1, p2 io.ReadWriteCloser, quiet bool) {
	if !quiet {
		fmt.Println("stream opened")
		defer fmt.Println("stream closed")
	}
	defer p1.Close()
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

func main() {
	myApp := cli.NewApp()
	myApp.Name = "quic2tcp"
	myApp.Usage = "a forwarding proxy from QUIC to tcp."
	myApp.Version = VERSION
	myApp.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "listen, l",
			Value: ":6935",
			Usage: "quic2tcp server listen address",
		},
		cli.StringFlag{
			Name:  "target, t",
			Value: "127.0.0.1:80",
			Usage: "target server address",
		},
		cli.IntFlag{
			Name:  "timeout",
			Value: 5,
			Usage: "max time of waiting a connection to complete",
		},
	}
	myApp.Action = func(c *cli.Context) error {
		config := Config{}
		config.Listen = c.String("listen")
		config.Target = c.String("target")
		config.Timeout = c.Int("timeout")

		lis, err := quicconn.Listen("udp", config.Listen, TLSConfig())
		if err != nil {
			panic(err)
		}

		log.Println("version:", VERSION)
		log.Println("listening on:", lis.Addr())
		log.Println("target:", config.Target)

		for {
			if p1, err := lis.Accept(); err == nil {
				log.Println("remote address:", p1.RemoteAddr())

				p2, err := net.DialTimeout("tcp", config.Target, time.Duration(config.Timeout) * time.Second)
				if err != nil {
					p1.Close()
					log.Println(err)
					continue
				}

				go handleClient(p1, p2, false)
			} else {
				log.Printf("%+v", err)
			}
		}
	}
	myApp.Run(os.Args)
}
