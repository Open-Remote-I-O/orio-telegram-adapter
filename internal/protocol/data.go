package protocol

import (
	"bytes"
	"encoding/binary"
)

type Command uint8

const (
	TestLedColor Command = iota
)

// OrioData is the body sent with expected command and eventual data in order to give detail about command
type OrioData struct {
	CommandID Command
	Len       uint16
	Data      []byte
}

func (o OrioData) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	if err := binary.Write(buf, binary.BigEndian, o.CommandID); err != nil {
		return nil, err
	}

	if err := binary.Write(buf, binary.BigEndian, uint16(len(o.Data))); err != nil {
		return nil, err
	}

	if _, err := buf.Write(o.Data); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (o *OrioData) UnmarshalBinary(data []byte) error {
	buf := bytes.NewReader(data)

	if err := binary.Read(buf, binary.BigEndian, &o.CommandID); err != nil {
		return err
	}

	if err := binary.Read(buf, binary.BigEndian, &o.Len); err != nil {
		return err
	}

	o.Data = make([]byte, o.Len)

	if _, err := buf.Read(o.Data); err != nil {
		return err
	}

	return nil
}
