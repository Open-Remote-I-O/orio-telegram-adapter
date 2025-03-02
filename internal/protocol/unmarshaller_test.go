package protocol

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"math/rand/v2"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testCase struct {
	name                string
	mockProtocolPayload *OrioPayload
	expected            *OrioPayload
	expectedError       error
}

func randUint16() uint16 {
	return uint16(rand.UintN(math.MaxUint16))
}

func generateMockProtocolBuffer(t testing.TB, val *OrioPayload) io.Reader {
	var w bytes.Buffer
	if val == nil {
		return &w
	}

	// Write header fields
	err := binary.Write(&w, binary.BigEndian, val.Header.Version)
	if err != nil {
		t.Fatalf("error writing header version: %v", err)
	}
	err = binary.Write(&w, binary.BigEndian, val.Header.DeviceID)
	if err != nil {
		t.Fatalf("error writing header device ID: %v", err)
	}
	err = binary.Write(&w, binary.BigEndian, val.Header.PayloadLen)
	if err != nil {
		t.Fatalf("error writing header payload length: %v", err)
	}

	// Write each OrioData element
	for _, dataItem := range val.Data {
		err = binary.Write(&w, binary.BigEndian, uint8(dataItem.CommandID))
		if err != nil {
			t.Fatalf("error writing OrioData command ID: %v", err)
		}
		err = binary.Write(&w, binary.BigEndian, uint16(dataItem.Len))
		if err != nil {
			t.Fatalf("error writing OrioData data length: %v", err)
		}
		err = binary.Write(&w, binary.BigEndian, dataItem.Data)
		if err != nil {
			t.Fatalf("error writing OrioData data: %v", err)
		}
	}

	return &w
}

func generateMockByteSlice(len int, fillVal byte) []byte {
	data := make([]byte, len) // Adjust length as needed
	for i := range data {
		data[i] = fillVal
	}
	return data
}

func Test_Unmarshal(t *testing.T) {
	tests := []testCase{
		{
			name: "ok",
			mockProtocolPayload: &OrioPayload{
				Header: OrioHeader{
					Version:    version,
					DeviceID:   uint32(10),
					PayloadLen: uint16(1),
				},
				Data: []OrioData{{
					CommandID: uint8(10),
					Len:       uint16(1),
					Data:      []byte{0xAA},
				}},
			},
			expected: &OrioPayload{
				Header: OrioHeader{Version: version, DeviceID: 10, PayloadLen: 1},
				Data: []OrioData{{
					CommandID: uint8(10),
					Len:       uint16(1),
					Data:      []byte{0xAA},
				}},
			},
			expectedError: nil,
		},
		{
			name: "ok multiple data payload",
			mockProtocolPayload: &OrioPayload{
				Header: OrioHeader{
					Version:    version,
					DeviceID:   uint32(10),
					PayloadLen: uint16(3),
				},
				Data: []OrioData{
					{
						CommandID: uint8(10),
						Len:       uint16(1),
						Data:      []byte{0xAA},
					},
					{
						CommandID: uint8(10),
						Len:       uint16(1),
						Data:      []byte{0xAA},
					},
					{
						CommandID: uint8(10),
						Len:       uint16(1),
						Data:      []byte{0xAA},
					},
				},
			},
			expected: &OrioPayload{
				Header: OrioHeader{Version: version, DeviceID: 10, PayloadLen: 3},
				Data: []OrioData{
					{
						CommandID: uint8(10),
						Len:       uint16(1),
						Data:      []byte{0xAA},
					},
					{
						CommandID: uint8(10),
						Len:       uint16(1),
						Data:      []byte{0xAA},
					},
					{
						CommandID: uint8(10),
						Len:       uint16(1),
						Data:      []byte{0xAA},
					},
				},
			},
			expectedError: nil,
		},
		{
			name: "ok no data sent",
			mockProtocolPayload: &OrioPayload{
				Header: OrioHeader{
					Version:    version,
					DeviceID:   uint32(10),
					PayloadLen: uint16(0),
				},
				Data: nil,
			},
			expected: &OrioPayload{
				Header: OrioHeader{Version: version, DeviceID: 10, PayloadLen: 0},
				Data:   nil,
			},
			expectedError: nil,
		},
		{
			name: "ok big paylaod",
			mockProtocolPayload: &OrioPayload{
				Header: OrioHeader{
					Version:    version,
					DeviceID:   uint32(10),
					PayloadLen: uint16(1),
				},
				Data: []OrioData{{
					CommandID: 0,
					Len:       chunkSize + 20,
					Data:      generateMockByteSlice(chunkSize+20, 0xFF),
				}},
			},
			expected: &OrioPayload{
				Header: OrioHeader{Version: version, DeviceID: 10, PayloadLen: 1},
				Data: []OrioData{{
					CommandID: 0,
					Len:       chunkSize + 20,
					Data:      generateMockByteSlice(chunkSize+20, 0xFF),
				}},
			},
			expectedError: nil,
		},
		{
			name:          "invalid header format passed",
			expected:      nil,
			expectedError: fmt.Errorf("%s: %w", ErrHeaderFormat, io.EOF),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockVal := generateMockProtocolBuffer(t, tt.mockProtocolPayload)
			res, err := Unmarshal(mockVal)
			if err != nil {
				assert.Equal(t, tt.expectedError, err)
				return
			}
			assert.Equal(t, tt.expected, res)
		})
	}
}

// TODO: learn more about fuzzing and best practices
func Fuzz_Unmarshal(f *testing.F) {
	// init fuzz corpus values
	var fuzzCorpusReaders []io.Reader
	fuzzCorpusVals := []*OrioPayload{
		{
			Header: OrioHeader{Version: version, DeviceID: 10, PayloadLen: 1},
			Data: []OrioData{{
				CommandID: uint8(10),
				Len:       uint16(1),
				Data:      []byte{0xAA},
			}},
		},
		{
			Header: OrioHeader{Version: version, DeviceID: 10, PayloadLen: 0},
			Data:   nil,
		},
		{
			Header: OrioHeader{Version: version, DeviceID: 10, PayloadLen: 3},
			Data: []OrioData{
				{
					CommandID: uint8(10),
					Len:       uint16(1),
					Data:      []byte{0xAA},
				},
				{
					CommandID: uint8(10),
					Len:       uint16(1),
					Data:      []byte{0xAA},
				},
				{
					CommandID: uint8(10),
					Len:       uint16(1),
					Data:      []byte{0xAA},
				},
			},
		},
	}

	for _, v := range fuzzCorpusVals {
		fuzzCorpusReaders = append(fuzzCorpusReaders, generateMockProtocolBuffer(f, v))
	}
	r := io.MultiReader(fuzzCorpusReaders...)
	validHeaderVal, err := io.ReadAll(r)
	if err != nil {
		f.Errorf("unexpected error while generating fuzz corpus")
	}

	fuzzCorpus := [][]byte{validHeaderVal}
	fmt.Println(fuzzCorpus)
	for _, v := range fuzzCorpus {
		f.Add(v)
	}

	f.Fuzz(func(t *testing.T, b []byte) {
		_, err := Unmarshal(bytes.NewReader(b))
		if err != nil && err != ErrHeaderFormatEOF {
			t.Errorf("given test case %v;\n caused: %s", b, err)
		}
	})
}

func generateRandomMockProtocolBuffer(t testing.TB) io.Reader {
	var w bytes.Buffer
	err := binary.Write(&w, binary.BigEndian, randUint16())
	if err != nil {
		t.Fatalf("something went wrong while generating mock protocol buffer")
	}
	err = binary.Write(&w, binary.BigEndian, rand.Uint32())
	if err != nil {
		t.Fatalf("something went wrong while generating mock protocol buffer")
	}
	err = binary.Write(&w, binary.BigEndian, randUint16())
	if err != nil {
		t.Fatalf("something went wrong while generating mock protocol buffer")
	}
	return &w
}

// TODO: verify the best way in order to benchmark unmarhsalling behaivour
func Benchmark_Unmarshal(b *testing.B) {
	b.Run("pre generated random message", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		val := generateRandomMockProtocolBuffer(b)
		for n := 0; n < b.N; n++ {
			_, _ = Unmarshal(val)
		}
	})
}
