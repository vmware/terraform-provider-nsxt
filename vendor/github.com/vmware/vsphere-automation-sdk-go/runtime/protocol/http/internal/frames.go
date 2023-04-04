/* Copyright © 2021 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package internal

import (
	"bufio"
	"errors"
	"io"
	"strconv"
)

const (
	MaxHexDigits int = 8
)

type VapiFrameReader struct {
	reader  *bufio.Reader
	isEmpty bool
}

// NewVapiFrameReader creates a new frame reader.
func NewVapiFrameReader(reader io.Reader) *VapiFrameReader {
	bufferedReader, ok := reader.(*bufio.Reader)
	if !ok {
		bufferedReader = bufio.NewReader(reader)
	}
	return &VapiFrameReader{reader: bufferedReader, isEmpty: true}
}

func (cr *VapiFrameReader) readLength() ([]byte, error) {
	data, err := cr.reader.ReadBytes('\n')
	if err != nil {
		if err == io.EOF &&
			(cr.isEmpty || len(data) > 0) {
			return data, io.ErrUnexpectedEOF
		}
		return data, err
	}
	cr.isEmpty = false
	length := len(data)
	if data[length-2] != '\r' {
		return data, errors.New("missing \\r at the end of frame's length segment")
	}

	if length > MaxHexDigits {
		return data, errors.New("provided length is above the max supported frame size")
	}

	return data[0 : length-2], nil
}

func (cr *VapiFrameReader) readData(frameLength int) ([]byte, error) {
	delimiterLength := 2
	data := make([]byte, frameLength+delimiterLength)
	numBytes, err := io.ReadFull(cr.reader, data)
	if err != nil {
		return data[0:numBytes], err
	}

	if data[frameLength] != '\r' || data[frameLength+1] != '\n' {
		return data[0:numBytes], errors.New("valid delimiter is not provided at the end of frame data")
	}

	return data[0:frameLength], nil
}

// ReadFrame processes the length and data from one VAPI frame.
// Returns a byte array with the read data or an error.
// The method will read data till it reaches EOF, then it will
// return []byte{}, io.EOF.
func (cr *VapiFrameReader) ReadFrame() ([]byte, error) {

	var intFrameLength int64
	readLength, err := cr.readLength()
	if err != nil {
		return readLength, err
	}

	intFrameLength, err = strconv.ParseInt(string(readLength), 16, 64)
	if err != nil {
		return readLength, err
	}

	data, err := cr.readData(int(intFrameLength))
	if err != nil {
		return data, err
	}
	return data, nil
}
