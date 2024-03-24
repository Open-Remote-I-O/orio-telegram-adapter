package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net"
	"os"
)

const (
	certificatesPath = "/etc/ssl/certs/"
	privateKeyPath   = "/etc/ssl/private/"
)

func main() {
	serverAddr, ok := os.LookupEnv("SERVER_ADDR_NAME")
	if !ok {
		panic("env not present")
	}

	deviceControlPort, ok := os.LookupEnv("DEVICE_CONTROL_PORT")
	if !ok {
		panic("env not present")
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
		"tcp",
		net.JoinHostPort(serverAddr, deviceControlPort),
		tlsConf,
	)
	if err != nil {
		fmt.Println("Something went wrong while connecting to server:", err)
		panic(err)
	}
	// TODO: for local test do not close but send a command each 30 sec?
	defer conn.Close()

	// Write and read data over the secure connection
	data := "Hello from mTLS client!\n"
	_, err = conn.Write([]byte(data))
	if err != nil {
		log.Fatal("Failed to write data:", err)
	}

	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		log.Fatal("Failed to read data:", err)
	}

	fmt.Println("Server response:", string(buf[:n]))
}
