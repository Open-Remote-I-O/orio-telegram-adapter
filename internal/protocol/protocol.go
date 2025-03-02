// Package goprotocol had the message structure that the server expects for a specific version
// and Marshal and Unmarshal methos in order to generate and serialize protocol communication
package protocol

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
