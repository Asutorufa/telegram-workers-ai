package ai

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"
)

type llamaStreamDecoder struct {
	r *bufio.Scanner
}

func NewLlamaStreamDecoder(r io.Reader) *llamaStreamDecoder {
	dec := &llamaStreamDecoder{bufio.NewScanner(r)}
	return dec
}

func (l *llamaStreamDecoder) Decode() (string, error) {
	for l.r.Scan() {
		text := l.r.Text()
		if text == "" {
			continue
		}

		sections := strings.SplitN(text, ":", 2)
		field, value := sections[0], ""
		if len(sections) == 2 {
			value = strings.TrimPrefix(sections[1], " ")
		}
		switch field {
		case "event":
		case "data":
			if value == "[DONE]" {
				return "", io.EOF
			}

			x := map[string]any{}
			err := json.Unmarshal([]byte(value), &x)
			if err != nil {
				log.Println("json unmarshal", "err", err, "value", value)
				continue
			}

			if x["response"] == nil {
				continue
			}

			return fmt.Sprint(x["response"]), nil
		case "id":
		case "retry":
		}
	}

	return "", io.EOF
}
