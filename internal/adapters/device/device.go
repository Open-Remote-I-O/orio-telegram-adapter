package adapters

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"sync"

	"orio-telegram-adapter/internal/config"
	"orio-telegram-adapter/internal/protocol"

	"github.com/rs/zerolog"
)

type DeviceHandler struct {
	logger  *zerolog.Logger
	server  net.Listener
	wg      sync.WaitGroup
	clients clientsStatus
}

func NewDeviceRemoteControlAdapter(
	logger *zerolog.Logger,
	conf config.DeviceConfig,
) (*DeviceHandler, error) {
	deviceControlPort, envIsPresent := os.LookupEnv("LOCAL_DEVICE_CONTROL_PORT")
	if !envIsPresent {
		logger.Warn().Msg("missing DEVICE_CONTROL_PORT env variable")
		return nil, errors.New("missing DEVICE_CONTROL_PORT env variable")
	}

	// Load client certificate and key
	orioCert, err := tls.LoadX509KeyPair(conf.Orio_tls_cert_path, conf.Orio_tls_key_path)
	if err != nil {
		logger.Error().AnErr("Failed to load client certificate/key:", err)
		return nil, err
	}

	// Load CA certificate
	caCert, err := os.ReadFile(conf.Orio_ca_cert_path)
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
		"tcp6",
		net.JoinHostPort("::", deviceControlPort),
		tlsConf,
	)
	if err != nil {
		logger.Fatal().Err(err)
		return nil, err
	}

	return &DeviceHandler{
		logger:  logger,
		server:  listen,
		wg:      sync.WaitGroup{},
		clients: clientsStatus{},
	}, nil
}

func (dh *DeviceHandler) Stop(ctx context.Context) error {
	dh.wg.Wait()
	return nil
}

func (dh *DeviceHandler) Start(ctx context.Context) error {
	defer func() {
		err := dh.server.Close()
		if err != nil {
			dh.logger.Err(err).Msg("failed to gracefully shut down device control service tcp server")
		}
	}()

	dh.logger.Info().Msg("starting device control service")
	for {
		conn, err := dh.server.Accept()
		if err != nil {
			dh.logger.Err(err).Msg("something went wrong while starting device control server, closing connection")
			conn.Close()
		}
		dh.logger.Debug().Msgf("handling connection from %s", conn.RemoteAddr())

		// parse header to deduct if the connection exists or not?
		// but accept returns the next connection in listener so we can assume that it will be a new one

		// so get first message that should our handshake with client, can directly buffer header only size for now
		// create new connection on success, enable keep alive and then start looping throught connection to next msgs

		fmt.Printf("\n %v --> %v \n", conn.RemoteAddr(), conn.LocalAddr())

		data, err := protocol.Unmarshal(conn)
		if err != nil {
			if err == io.EOF {
				fmt.Println("Connection closed by client (keep-alive failure or normal close)")
				continue
			}
			// fix this unreadable mess
			if netErr, ok := err.(net.Error); ok {
				fmt.Printf("Network error: %s, Timeout: %v", netErr.Error(), netErr.Timeout())
				continue
			} else {
				fmt.Printf("Other error: %v", err)
				continue
			}
		}

		fmt.Printf("\n Received: %+v \n", data)
		dh.logger.Info().Msg("received first handshake, persisting connection and reading from it")

		deviceConn := NewDeviceConnection(conn, string(data.Header.DeviceID))

		dh.wg.Add(1)
		go func() {
			dh.handleClient(deviceConn.conn)
			dh.wg.Done()
		}()
	}
}

// move to device conn since the adapter here will only do first setup and handle spin up and down of tcp server
func (dh *DeviceHandler) handleClient(con net.Conn) {
	remoteAddr := con.RemoteAddr()

	defer func() {
		fmt.Printf("%v disconnected \n", remoteAddr)
	}()

	enableTcpKeepAlive(con)

	fmt.Printf("\n %v --> %v \n", con.RemoteAddr(), con.LocalAddr())

	for {
		data, err := protocol.Unmarshal(con)
		if err != nil {
			if err == io.EOF {
				fmt.Println("Connection closed by client (keep-alive failure or normal close)")
				return
			}
			if netErr, ok := err.(net.Error); ok {
				fmt.Printf("Network error: %s, Timeout: %v", netErr.Error(), netErr.Timeout())
			} else {
				fmt.Printf("Other error: %v", err)
			}
			fmt.Println("Error reading:", err)
			return
		}

		fmt.Printf("\n Received: %#v \n", data)
	}
}
