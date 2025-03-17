package adapters

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"orio-telegram-adapter/internal/protocol"
	"strconv"
)

type deviceConnection struct {
	conn   *net.Conn
	id     string
	SendCh chan protocol.OrioPayload
}

func NewDeviceConnection(
	conn *net.Conn,
	id string,
) deviceConnection {
	return deviceConnection{
		conn: conn, //TODO: evaluate if this should be a pointer
		id:   id,
		// TODO: evaluate how to buffer it or not
		SendCh: make(chan protocol.OrioPayload, 100),
	}
}

func enableTcpKeepAlive(con *net.Conn) {
	tcpCon, ok := (*con).(*tls.Conn).NetConn().(*net.TCPConn)
	if !ok {
		// not tcp connection, shoudln't really happen in first place but safeguards are nice
		panic("implement me")
	}

	tcpCon.SetKeepAliveConfig(net.KeepAliveConfig{
		Enable:   true,
		Idle:     15, //NOTE: default value pick a more thought one
		Interval: 15, //NOTE: default value pick a more thought one
		Count:    9,  //NOTE: default value pick a more thought one
	})
}

// move to device conn since the adapter here will only do first setup and handle spin up and down of tcp server
func (dc deviceConnection) HandleClient(con *net.Conn) {
	remoteAddr := (*con).RemoteAddr()

	defer func() {
		fmt.Printf("%v disconnected \n", remoteAddr)
	}()

	fmt.Println("setting up tcp keepalive")

	enableTcpKeepAlive(con)

	fmt.Println("tcp keepalive configured")

	fmt.Printf("Handling persistent tcp socket \n %v --> %v \n", (*con).RemoteAddr(), (*con).LocalAddr())

	// plan connection shutdown in case of server restart, basically propagate
	// temination signal to properly clean connections later on
	// placeholder for cancelation, context should be propagated for entire device server
	ctx, cancel := context.WithCancel(context.Background())
	for {
		select {
		case <-ctx.Done():
			cancel()
		case val := <-dc.SendCh:
			gugu, err := strconv.Atoi(dc.id)
			if err != nil {
				fmt.Println("printato errore gigetto nel convertire dc id", err.Error())
			}
			fmt.Printf("Connection handler %v, Received from message bus; sending data to device %#v\n", gugu, val)
			binVal, err := val.MarshalBinary()
			if err != nil {
				fmt.Println("\n failed to send stuff in this connection %+v\n", binVal)
				continue
			}
			// TODO: add random color generation to send to each connection,
			// this will be then handled by ui or color code will be used to rappresent
			// device status?
			n, err := (*dc.conn).Write(binVal)
			if err != nil {
				fmt.Println("\n not able to send %+v \n %v --> %v \n", val, (*con).RemoteAddr(), (*con).LocalAddr())
				continue
			}
			fmt.Printf("\n successfully written %d bytes to tcp socket of device %s \n", n, (*con).RemoteAddr())
			// default:
			// 	fmt.Println("default not doing anything")
			// fmt.Println("going into default device connection case")
			// data, err := protocol.Unmarshal(con)
			// fmt.Println("finished reading data", data, err.Error())
			// if err != nil {
			// 	if err == io.EOF {
			// 		fmt.Println("Connection closed by client (keep-alive failure or normal close)")
			// 		return
			// 	}
			// 	if netErr, ok := err.(net.Error); ok {
			// 		fmt.Printf("\n Network error: %s, Timeout: %v \n", netErr.Error(), netErr.Timeout())
			// 	} else {
			// 		fmt.Printf("\n Other error: %v \n", err)
			// 	}
			// 	fmt.Printf("\n Error reading: \n %w \n", err)
			// 	return
			// }

			// fmt.Printf("\n Connection handler received: %#v \n", data)
		}
	}
}
