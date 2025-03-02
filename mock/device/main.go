package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net"
	"orio-telegram-adapter/internal/protocol"
	"os"
	"time"
)

const (
	certificatesPath = "/etc/ssl/certs/"
	privateKeyPath   = "/etc/ssl/private/"
)

func main() {
	serverAddr, ok := os.LookupEnv("SERVER_ADDR_NAME")
	if !ok {
		panic("SERVER_ADDR_NAME env not present")
	}

	deviceControlPort, ok := os.LookupEnv("LOCAL_DEVICE_CONTROL_PORT")
	if !ok {
		panic("LOCAL_DEVICE_CONTROL_PORT env not present")
	}

	// Load client certificate and key
	mockClientCert, err := tls.LoadX509KeyPair(
		certificatesPath+"mock-device.crt",
		privateKeyPath+"mock-device.key",
	)
	if err != nil {
		log.Fatal("Failed to load client certificate/key:", err)
	}

	// Load CA certificate
	caCert, err := os.ReadFile(certificatesPath + "orio-ca.crt")
	if err != nil {
		log.Fatal("Failed to read CA certificate:", err)
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Configure TLS client
	tlsConf := &tls.Config{
		Certificates:     []tls.Certificate{mockClientCert},
		RootCAs:          caCertPool,
		ClientAuth:       tls.RequireAndVerifyClientCert,
		CurvePreferences: []tls.CurveID{tls.CurveP256},
	}

	conn, err := tls.Dial(
		"tcp6",
		net.JoinHostPort(serverAddr, deviceControlPort),
		tlsConf,
	)
	if err != nil {
		fmt.Println("Something went wrong while connecting to server:", err)
		panic(err)
	}
	defer conn.Close()

	t := time.NewTicker(10 * time.Second)
	defer t.Stop()
	for {
		select {
		case <-t.C:
			var rawPayload protocol.OrioHeader
			rawPayload.DeviceID = 6969
			rawPayload.Version = 0
			rawPayload.PayloadLen = 0
			payload, err := rawPayload.MarshalBinary()
			if err != nil {
				log.Printf("client write marshal header error: %s", err)
				continue
			}

			_, err = conn.Write(payload)
			if err != nil {
				log.Printf("client socker write errors: %s", err)
				continue
			}
			log.Printf("client: sent: %#v", payload)
		}
	}
}
