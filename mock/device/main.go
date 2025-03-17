package main

import (
	"bytes"
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

func connect(serverAddr string, deviceControlPort string, tlsConf *tls.Config) net.Conn {
	for {
		time.Sleep(time.Second)
		conn, err := tls.Dial(
			"tcp4",
			net.JoinHostPort(serverAddr, deviceControlPort),
			tlsConf,
		)
		if err != nil {
			fmt.Println("Something went wrong while connecting to server:", err)
			continue
		}
		return conn
	}
}

func main() {
	serverAddr, ok := os.LookupEnv("SERVER_ADDR_NAME")
	if !ok {
		panic("SERVER_ADDR_NAME env not present")
	}

	deviceControlPort, ok := os.LookupEnv("DEVICE_CONTROL_PORT")
	if !ok {
		panic("DEVICE_CONTROL_PORT env not present")
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

	for {
		conn := connect(serverAddr, deviceControlPort, tlsConf)

		fmt.Println("handled first handshake stopping looping")

		scheduleSendMockData(&conn)
		conn.Close()
	}
}

func scheduleSendMockData(conn *net.Conn) {
	t := time.NewTicker(10 * time.Second)
	defer t.Stop()

	t2 := time.NewTicker(5 * time.Second)
	defer t2.Stop()

	go func() {
		fmt.Println("Unmarshalling from G routine")
		res, err := protocol.Unmarshal(*conn)
		fmt.Println("something??")
		if err != nil {
			fmt.Printf("\nError catched on return: %s\n", err)
		}
		fmt.Printf("\nResponse fetched from mock: %+v\n", res)
	}()

	for {
		select {
		case <-t2.C:
			fmt.Println("trying to read from tcp socket")
			fmt.Println("debug info from conn", (*conn).RemoteAddr())
			fmt.Println("local addr", (*conn).LocalAddr())
			var gne bytes.Buffer
			n, err := (*conn).Read(gne.Bytes())
			if err != nil {
				fmt.Println("Failed to read from conn", err)
			}
			fmt.Printf("read %v bytes \n", n)
			fmt.Printf("read %x data \n", gne.Bytes())
		case <-t.C:
			var rawPayload protocol.OrioHeader
			rawPayload.DeviceID = 6969
			rawPayload.Version = 0
			rawPayload.PayloadLen = 0
			payload, err := rawPayload.MarshalBinary()
			if err != nil {
				log.Printf("client write marshal header error: %s", err)
				return
			}

			_, err = (*conn).Write(payload)
			if err != nil {
				log.Printf("client socker write errors: %s", err)
				return
			}
			log.Printf("mockmock  client: sent: %#v", payload)
		}
	}
}
