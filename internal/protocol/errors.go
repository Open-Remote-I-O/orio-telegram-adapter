package protocol

import (
	"fmt"
	"io"
)

var (
	ErrHeaderFormat    = "invalid protocol header format sent"
	ErrHeaderFormatEOF = fmt.Errorf("%s: %w", ErrHeaderFormat, io.EOF)

	ErrDataFormat    = "invalid protocol data format sent"
	ErrDataFormatEOF = fmt.Errorf("%s: %w", ErrDataFormat, io.EOF)

	ErrDataLen    = "length provided in the data params is not equal to actual data sent"
	ErrDataLenEOF = fmt.Errorf("%s: %w", ErrDataLen, io.EOF)
)
