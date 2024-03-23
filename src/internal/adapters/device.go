package adapters

import (
	"bufio"
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"os"

	"github.com/rs/zerolog"
)

type DeviceHandler struct {
	logger *zerolog.Logger
	server net.Listener
}

const (
	certificatesPath = "/etc/ssl/certs/"
	privateKeyPath   = "/etc/ssl/private/"
)

func NewDeviceRemoteControlAdapter(
	logger *zerolog.Logger,
) (*DeviceHandler, error) {
	deviceControlPort, envIsPresent := os.LookupEnv("DEVICE_CONTROL_PORT")
	if !envIsPresent {
		fmt.Println("missing env variable")
		return nil, errors.New("missing DEVICE_CONTROL_PORT env variable")
	}

	// Load client certificate and key
	orioCert, err := tls.LoadX509KeyPair(certificatesPath+"orio-server.crt", privateKeyPath+"orio-server.key")
	if err != nil {
		logger.Error().AnErr("Failed to load client certificate/key:", err)
		return nil, err
	}

	// Load CA certificate
	caCert, err := os.ReadFile(certificatesPath + "orio-ca.crt")
	if err != nil {
		logger.Error().AnErr("Failed to read CA certificate:", err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Configure TLS client
	tlsConf := &tls.Config{
		Certificates:       []tls.Certificate{orioCert},
		RootCAs:            caCertPool,
		InsecureSkipVerify: false,
	}

	listen, err := tls.Listen(
		"tcp",
		net.JoinHostPort("0.0.0.0", deviceControlPort),
		tlsConf,
	)
	if err != nil {
		logger.Fatal().Err(err)
	}

	return &DeviceHandler{
		logger: logger,
		server: listen,
	}, nil
}

func (dh *DeviceHandler) StartServer(ctx context.Context) {
	for {
		conn, err := dh.server.Accept()
		if err != nil {
			dh.logger.Err(err).Msg("something went wrong while starting device control server")
			panic(err)
		}
		fmt.Println("handling connection from", conn.RemoteAddr())
		handleClient(conn)
	}
}

func handleClient(connection net.Conn) {
	remoteAddr := connection.RemoteAddr()

	defer func() {
		fmt.Printf("%v disconnected \n", remoteAddr)
	}()

	data, err := bufio.NewReader(connection).ReadString('\n')
	if err != nil {
		fmt.Println("error during data reading operation", err)
		return
	}
	fmt.Printf("%v <= %v: and %v", connection.LocalAddr(), connection.RemoteAddr(), data)

	randMessage := fmt.Sprintf("Message! %v\n", rand.Intn(100000))
	fmt.Printf("%v => %v: and %v", connection.LocalAddr(), connection.RemoteAddr(), randMessage)
	_, err = connection.Write([]byte(randMessage))
	if err != nil {
		fmt.Println("error writing data to connection", err)
		return
	}
}
