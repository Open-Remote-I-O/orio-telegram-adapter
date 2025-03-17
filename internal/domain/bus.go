package domain

type Command struct {
	DeviceID  string
	CommandID uint8
	Payload   []byte
}

type MessageBus struct {
	CommandCh chan Command
}
