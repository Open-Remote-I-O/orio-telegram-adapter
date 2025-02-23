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
	"sync"

	"orio-telegram-adapter/src/internal/config"

	"github.com/rs/zerolog"
)

type DeviceHandler struct {
	logger *zerolog.Logger
	server net.Listener
	wg     sync.WaitGroup
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
		net.JoinHostPort("::1", deviceControlPort),
		tlsConf,
	)
	if err != nil {
		logger.Fatal().Err(err)
		return nil, err
	}

	return &DeviceHandler{
		logger: logger,
		server: listen,
		wg:     sync.WaitGroup{},
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
		dh.wg.Add(1)
		go func() {
			handleClient(conn)
			dh.wg.Done()
		}()
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
