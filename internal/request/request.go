package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
)

type Request struct {
	RequestLine RequestLine

	state requestState
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type requestState int

const (
	requestStateInitialized requestState = iota
	requestStateDone
)

const crlf = "\r\n"
const bufferSize = 8

func RequestFromReader(reader io.Reader) (*Request, error) {
	buffer := make([]byte, bufferSize, bufferSize)
	readToIndex := 0

	req := &Request{
		state: requestStateInitialized,
	}

	for req.state != requestStateDone {
		if readToIndex >= len(buffer) {
			newBuffer := make([]byte, len(buffer)*2)
			copy(newBuffer, buffer)
			buffer = newBuffer
		}

		numBytesRead, err := reader.Read(buffer[readToIndex:])
		if err != nil {
			if errors.Is(err, io.EOF) {
				req.state = requestStateDone
				break
			}
			return nil, err
		}

		readToIndex += numBytesRead

		numBytesParsed, err := req.parse(buffer[:readToIndex])
		if err != nil {
			return nil, err
		}

		copy(buffer, buffer[numBytesParsed:])
		readToIndex -= numBytesParsed
	}
	return req, nil
}

func parseRequestLine(data []byte) (*RequestLine, int, error) {
	index := bytes.Index(data, []byte(crlf))
	if index == -1 {
		return nil, 0, nil
	}

	requestLineText := string(data[:index])
	requestLine, err := requestLineFromString(requestLineText)
	if err != nil {
		return nil, 0, err
	}
	return requestLine, index + 2, nil
}

func requestLineFromString(str string) (*RequestLine, error) {
	parts := strings.Split(str, " ")
	if len(parts) != 3 {
		return nil, fmt.Errorf("error: poorly formatted request-line - %s", str)
	}

	method := parts[0]
	for _, char := range method {
		if char < 'A' || char > 'Z' {
			return nil, fmt.Errorf("error: invalid method - %s", method)
		}
	}

	requestTarget := parts[1]

	versionParts := strings.Split(parts[2], "/")
	if len(versionParts) != 2 {
		return nil, fmt.Errorf("error: malformed start-line - %s", str)
	}

	httpPart := versionParts[0]
	if httpPart != "HTTP" {
		return nil, fmt.Errorf("error: unrecognized HTTP-version - %s", httpPart)
	}
	version := versionParts[1]
	if version != "1.1" {
		return nil, fmt.Errorf("error: unrecognized HTTP-version - %s", version)
	}

	return &RequestLine{
		Method:        method,
		RequestTarget: requestTarget,
		HttpVersion:   version,
	}, nil
}

func (r *Request) parse(data []byte) (int, error) {
	switch r.state {
	case requestStateInitialized:
		requestLine, n, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		}
		if n == 0 {
			return 0, nil
		}
		r.RequestLine = *requestLine
		r.state = requestStateDone
		return n, nil
	case requestStateDone:
		return 0, fmt.Errorf("error: trying to read data in a done state")
	default:
		return 0, fmt.Errorf("unknown state")
	}
}
