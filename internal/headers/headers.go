package headers

import (
	"fmt"
	"strings"
)

type Headers map[string]string

const crlf = "\r\n"

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	byteString := string(data)

	byteString = strings.Trim(byteString, " ")

	index := strings.Index(byteString, crlf)

	switch index {
	case -1:
		return 0, false, nil
	case 0:
		return 2, true, nil
	}

	line := byteString[:index]

	colIdx := strings.Index(line, ":")

	if colIdx == -1 {
		return 0, false, fmt.Errorf("error: no colon present")
	} else if colIdx == 0 {
		return 0, false, fmt.Errorf("error: empty key")
	} else if line[colIdx-1] == ' ' || line[colIdx-1] == '\t' {
		return 0, false, fmt.Errorf("error: key formatting issue")
	}

	key := line[:colIdx]

	if strings.Contains(key, " ") || strings.Contains(key, "\t") {
		return 0, false, fmt.Errorf("error: key formatting issue")
	}

	value := strings.Trim(line[colIdx+1:], " ")

	h[key] = value

	return index + len(crlf), false, nil
}

func NewHeaders() Headers {
	return make(map[string]string)
}
