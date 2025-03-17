// Package goprotocol had the message structure that the server expects for a specific version
// and Marshal and Unmarshal methos in order to generate and serialize protocol communication
package protocol

import (
	"bytes"
	"encoding/binary"
)

// NOTE: currently commented other protocol values in order to test basic implementation of the unmarshalling

const version = uint16(0)

const (
	headerByteSize  = 8
	mindataByteSize = 3
	chunkSize       = 128
)

// OrioPayload is the complete payload that a client will be sending to server
type OrioPayload struct {
	Header OrioHeader
	Data   []OrioData
}

func (o OrioPayload) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	val, err := o.Header.MarshalBinary()
	if err != nil {
		panic(err)
	}

	if err := binary.Write(buf, binary.BigEndian, val); err != nil {
		return nil, err
	}

	// val2, err := o.Data.MarshalBinary()
	// if err := binary.Write(buf, binary.BigEndian, val2); err != nil {
	// 	return nil, err
	// }

	return buf.Bytes(), nil
}
