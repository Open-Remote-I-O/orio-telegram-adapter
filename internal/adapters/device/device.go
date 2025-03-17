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

	"golang.org/x/exp/maps"
	"orio-telegram-adapter/internal/config"
	"orio-telegram-adapter/internal/domain"
	"orio-telegram-adapter/internal/protocol"

	"github.com/rs/zerolog"
)

type DeviceHandler struct {
	logger     *zerolog.Logger
	server     net.Listener
	wg         sync.WaitGroup
	clientsMx  sync.RWMutex
	clients    map[string]deviceConnection
	messageBus domain.MessageBus
}

func NewDeviceRemoteControlAdapter(
	logger *zerolog.Logger,
	conf config.DeviceConfig,
	messageBus domain.MessageBus,
) (*DeviceHandler, error) {
	deviceControlPort, envIsPresent := os.LookupEnv("DEVICE_CONTROL_PORT")
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
		"tcp4",
		net.JoinHostPort("0.0.0.0", deviceControlPort),
		tlsConf,
	)
	if err != nil {
		logger.Fatal().Err(err)
		return nil, err
	}

	return &DeviceHandler{
		logger:     logger,
		server:     listen,
		wg:         sync.WaitGroup{},
		clients:    map[string]deviceConnection{},
		messageBus: messageBus,
	}, nil
}

func (dh *DeviceHandler) Stop(ctx context.Context) error {
	dh.wg.Wait()
	return nil
}

func (dh *DeviceHandler) receiveEvents() {
	for val := range dh.messageBus.CommandCh {
		fmt.Println("number of connected device", len(maps.Keys(dh.clients)))
		deviceConn, ok := dh.clients[val.DeviceID]
		if !ok {
			continue
		}
		fmt.Printf("Sending to device connection %v; chan length %v; \n", val.DeviceID, len(deviceConn.SendCh))
		deviceConn.SendCh <- protocol.OrioPayload{
			Header: protocol.OrioHeader{
				Version:    0,
				DeviceID:   6969,
				PayloadLen: 0,
			},
		}
	}
}

func (dh *DeviceHandler) Start(ctx context.Context) error {
	defer func() {
		err := dh.server.Close()
		if err != nil {
			dh.logger.Err(err).Msg("failed to gracefully shut down device control service tcp server")
		}
	}()

	go dh.receiveEvents()

	dh.logger.Info().Msg("starting device control service")
	for {
		conn, err := dh.server.Accept()
		if err != nil {
			dh.logger.Err(err).Msg("something went wrong while starting device control server, closing connection")
			conn.Close()
		}
		dh.logger.Debug().Msgf("handling connection from %s", conn.RemoteAddr().String())

		// parse header to deduct if the connection exists or not?
		// but accept returns the next connection in listener so we can assume that it will be a new one

		// so get first message that should our handshake with client, can directly buffer header only size for now
		// create new connection on success, enable keep alive and then start looping throught connection to next msgs

		fmt.Printf("\n%v --> %v\n", conn.RemoteAddr(), conn.LocalAddr())

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

		fmt.Printf("\nFirst connection payload received: %+v\n", data)
		fmt.Println("sanity check", strconv.FormatUint(uint64(data.Header.DeviceID), 10))
		dh.logger.Info().Msg("received first handshake, persisting connection and reading from it")

		deviceConn := NewDeviceConnection(&conn, string(data.Header.DeviceID))

		dh.clientsMx.Lock()
		// for now ok since it's still poc state, later on a more appropriate key should be formed
		// considering also how the user side will be done, since in state like this it could cause
		// conflicts due to low number of devices
		fmt.Println("persisting in clients miao", strconv.FormatUint(uint64(data.Header.DeviceID), 10))
		dh.clients[strconv.FormatUint(uint64(data.Header.DeviceID), 10)] = deviceConn

		dh.clientsMx.Unlock()

		dh.wg.Add(1)
		go func() {
			deviceConn.HandleClient(deviceConn.conn)
			dh.wg.Done()
		}()
	}
}
