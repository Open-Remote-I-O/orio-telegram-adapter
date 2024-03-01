package adapters

import (
	"bufio"
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"math/rand"
	"net"
	"os"
	"time"

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
) (DeviceHandler, error) {
	_, envIsPresent := os.LookupEnv("DEVICE_CONTROL_PORT")
	if !envIsPresent {
		fmt.Println("missing env variable")
	}

	// Load client certificate and key
	orioCert, err := tls.LoadX509KeyPair(certificatesPath+"orio-server.crt", privateKeyPath+"orio-server.key")
	if err != nil {
		logger.Error().AnErr("Failed to load client certificate/key:", err)
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

	fmt.Println(net.JoinHostPort("0.0.0.0", "23333"))

	listen, err := tls.Listen(
		"tcp",
		net.JoinHostPort("0.0.0.0", "23333"),
		tlsConf,
	)
	if err != nil {
		logger.Fatal().Err(err)
		os.Exit(1)
	}
	defer listen.Close()

	return DeviceHandler{
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
		dh.logger.Debug().Msg(fmt.Sprintf("listening connection on: %s", dh.server.Addr().String()))
		go handleClient(conn)
	}
}

func handleClient(connection net.Conn) {
	remoteAddr := connection.RemoteAddr()

	defer func() {
		fmt.Printf("%v disconnected", remoteAddr)
	}()

	for {
		data, err := bufio.NewReader(connection).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("%v <= %v: and %v", connection.LocalAddr(), connection.RemoteAddr(), data)

		time.Sleep(time.Duration(rand.Int31n(15)) * time.Second)

		randMessage := fmt.Sprintf("Message! %v\n", rand.Intn(100000))
		fmt.Printf("%v => %v: and %v", connection.LocalAddr(), connection.RemoteAddr(), randMessage)
		_, err = connection.Write([]byte(randMessage))
		if err != nil {
			fmt.Println("test")
		}
	}
}
