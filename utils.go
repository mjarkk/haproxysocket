package haproxysocket

import (
	"bytes"
	"errors"
	"io"
	"net"
	"strings"
)

// q executes a query
func (h *HaproxyInstace) q(query string) (string, error) {
	c, err := net.Dial(h.Network, h.Address)
	if err != nil {
		panic(err)
	}
	defer c.Close()

	c.Write([]byte(query + "\n"))

	var buf bytes.Buffer
	io.Copy(&buf, c)

	return strings.TrimSpace(buf.String()), nil
}

func (h *HaproxyInstace) qMap(query string, alternateSplit ...string) ([]map[string]string, error) {
	out, err := h.q(query)
	if err != nil {
		return []map[string]string{}, err
	}
	return csvToArrMap(out, alternateSplit...)
}

func csvToArrMap(in string, alternateSplit ...string) ([]map[string]string, error) {
	toReturn := []map[string]string{}

	splitter := ","
	switch len(alternateSplit) {
	case 0:
	case 1:
		splitter = alternateSplit[0]
	default:
		return toReturn, errors.New("alternateSplit can't be more than 1")
	}

	lines := strings.Split(in, "\n")
	if len(lines) == 0 {
		return toReturn, errors.New("No output")
	}

	title := lines[0]
	lines = lines[1:]

	parts := strings.Split(title, splitter)
	parts[0] = strings.TrimLeft(parts[0], "# ")

	for _, line := range lines {
		lineParts := strings.Split(line, splitter)
		toAdd := map[string]string{}
		for i, part := range lineParts {
			if i >= len(parts) {
				break
			}
			partName := parts[i]
			if partName == "" {
				continue
			}
			toAdd[partName] = part
		}
		toReturn = append(toReturn, toAdd)
	}

	return toReturn, nil
}
