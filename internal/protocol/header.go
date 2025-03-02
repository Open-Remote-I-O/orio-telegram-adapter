package protocol

import (
	"bytes"
	"encoding/binary"
)

// OrioHeader has all metadata needed before handling actual protocol data
type OrioHeader struct {
	Version    uint16
	DeviceID   uint32
	PayloadLen uint16
}

func (m OrioHeader) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	if err := binary.Write(buf, binary.BigEndian, m.Version); err != nil {
		return nil, err
	}

	if err := binary.Write(buf, binary.BigEndian, m.DeviceID); err != nil {
		return nil, err
	}

	if err := binary.Write(buf, binary.BigEndian, m.PayloadLen); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (m OrioHeader) UnmarshalBinary(data []byte) error {
	buf := bytes.NewReader(data)

	if err := binary.Read(buf, binary.BigEndian, &m.Version); err != nil {
		return err
	}

	if err := binary.Read(buf, binary.BigEndian, &m.DeviceID); err != nil {
		return err
	}

	if err := binary.Read(buf, binary.BigEndian, &m.PayloadLen); err != nil {
		return err
	}

	return nil
}
