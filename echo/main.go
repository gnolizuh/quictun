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
	"io"
	"os"
)

var (
	VERSION = "1.0"
)

func main() {
	myApp := cli.NewApp()
	myApp.Name = "QUIC echo server"
	myApp.Usage = "Echo QUIC data immediately."
	myApp.Version = VERSION
	myApp.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "listen, l",
			Value: ":6935",
			Usage: "server listen address",
		},
		cli.BoolFlag{
			Name:  "quiet",
			Usage: "to suppress the 'stream open/close' messages",
		},
	}
	myApp.Action = func(c *cli.Context) error {
		config := Config{}
		config.Listen = c.String("listen")
		config.Quiet = c.Bool("quiet")

		log.SetFlags(log.LstdFlags | log.Lmicroseconds)

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
		log.Println("quiet:", config.Quiet)

		// echo data from read side to write side.
		echo := func(p1 io.ReadWriteCloser) {
			if !config.Quiet {
				log.Println("stream opened")
				defer log.Println("stream closed")
			}
			defer p1.Close()

			n, err := io.Copy(p1, p1)
			if err != nil {
				log.Println(err)
			}

			log.Printf("echo %d bytes", n)
		}

		for {
			if p1, err := listener.Accept(); err == nil {
				log.Printf("accpet quic addr:%s\n", p1.RemoteAddr())

				go echo(p1)
			} else {
				log.Fatalln(err)
			}
		}
	}
	myApp.Run(os.Args)
}
