package adapters

import (
	"crypto/tls"
	"net"
	"orio-telegram-adapter/internal/protocol"
)

type clientsStatus map[string]deviceConnection

type deviceConnection struct {
	conn   net.Conn
	id     string
	sendCh <-chan protocol.OrioPayload
}

func NewDeviceConnection(
	conn net.Conn,
	id string,
) deviceConnection {
	return deviceConnection{
		conn:   conn, //TODO: evaluate if this should be a pointer
		id:     id,
		sendCh: make(<-chan protocol.OrioPayload),
	}
}

func enableTcpKeepAlive(con net.Conn) {
	tcpCon, ok := con.(*tls.Conn).NetConn().(*net.TCPConn)
	if !ok {
		// not tcp connection, shoudln't really happen in first place but safeguards are nice
		panic("implement me")
	}

	tcpCon.SetKeepAliveConfig(net.KeepAliveConfig{
		Enable:   true,
		Idle:     15, //NOTE: default value pick and reasona about a more thought one
		Interval: 15, //NOTE: default value pick and reasona about a more thought one
		Count:    9,  //NOTE: default value pick and reasona about a more thought one
	})
}
