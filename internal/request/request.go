package request

import (
	"fmt"
	"io"
	"strings"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	requestBytes, err := io.ReadAll(reader)
	if err != nil {
		return &Request{}, err
	}

	requestString := string(requestBytes)

	requestArr := strings.Split(requestString, "\n")

	reqLine, err := parseRequestLine(requestArr[0])
	if err != nil {
		return &Request{}, err
	}

	return &Request{
		RequestLine: reqLine,
	}, err
}

func parseRequestLine(requestLine string) (RequestLine, error) {
	reqLineInfo := strings.Split(requestLine, " ")

	if len(reqLineInfo) < 3 {
		err := fmt.Errorf("error: not enough arguments for request")
		return RequestLine{}, err
	}

	for i := range reqLineInfo {
		reqLineInfo[i] = strings.Trim(reqLineInfo[i], "\r")
	}

	method := reqLineInfo[0]
	reqTarget := reqLineInfo[1]
	httpVersion := reqLineInfo[2]

	if strings.Compare(method, strings.ToUpper(method)) != 0 {
		err := fmt.Errorf("error: incorrect request method syntax")
		return RequestLine{}, err
	}

	verNum := strings.ReplaceAll(httpVersion, "HTTP/", "")
	if verNum != "1.1" {
		err := fmt.Errorf("error: invalid http version")
		return RequestLine{}, err
	}

	newRequestLine := RequestLine{
		HttpVersion:   verNum,
		RequestTarget: reqTarget,
		Method:        method,
	}

	return newRequestLine, nil
}
