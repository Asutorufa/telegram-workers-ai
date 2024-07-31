package ai

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type llamaStreamDecoder struct {
	*Decoder
}

func NewLlamaStreamDecoder(r io.Reader) *llamaStreamDecoder {
	dec := &llamaStreamDecoder{NewDecoder(r)}
	return dec
}

func (l *llamaStreamDecoder) Decode() (string, error) {
	e, err := l.Decoder.Decode()
	if err != nil {
		return "", err
	}

	if e.Data() == "" {
		return "", nil
	}

	x := map[string]any{}
	err = json.Unmarshal([]byte(e.Data()), &x)
	if err != nil {
		return "", err
	}

	return fmt.Sprint(x["response"]), nil
}

// Any event received by the client or sent by the server will implement this interface
type Event interface {
	// Id is an identifier that can be used to allow a client to replay
	// missed Events by returning the Last-Event-Id header.
	// Return empty string if not required.
	Id() string
	// The name of the event. Return empty string if not required.
	Event() string
	// The payload of the event.
	Data() string
}

type publication struct {
	id, event, data string
	retry           int64
}

func (s *publication) Id() string    { return s.id }
func (s *publication) Event() string { return s.event }
func (s *publication) Data() string  { return s.data }
func (s *publication) Retry() int64  { return s.retry }

// A Decoder is capable of reading Events from a stream.
type Decoder struct {
	*bufio.Reader
}

// NewDecoder returns a new Decoder instance that reads events
// with the given io.Reader.
func NewDecoder(r io.Reader) *Decoder {
	dec := &Decoder{bufio.NewReader(newNormaliser(r))}
	return dec
}

// Decode reads the next Event from a stream (and will block until one
// comes in).
// Graceful disconnects (between events) are indicated by an io.EOF error.
// Any error occuring mid-event is considered non-graceful and will
// show up as some other error (most likely io.ErrUnexpectedEOF).
func (dec *Decoder) Decode() (Event, error) {
	// peek ahead before we start a new event so we can return EOFs
	_, err := dec.Peek(1)
	if err == io.ErrUnexpectedEOF {
		err = io.EOF
	}
	if err != nil {
		return nil, err
	}
	pub := new(publication)
	inDecoding := false
	for {
		line, err := dec.ReadString('\n')
		if err != nil {
			return nil, err
		}
		if line == "\n" && inDecoding {
			// the empty line signals the end of an event
			break
		} else if line == "\n" && !inDecoding {
			// only a newline was sent, so we don't want to publish an empty event but try to read again
			continue
		}
		line = strings.TrimSuffix(line, "\n")
		if strings.HasPrefix(line, ":") {
			continue
		}
		sections := strings.SplitN(line, ":", 2)
		field, value := sections[0], ""
		if len(sections) == 2 {
			value = strings.TrimPrefix(sections[1], " ")
		}
		inDecoding = true
		switch field {
		case "event":
			pub.event = value
		case "data":
			pub.data += value + "\n"
		case "id":
			pub.id = value
		case "retry":
			pub.retry, _ = strconv.ParseInt(value, 10, 64)
		}
	}
	pub.data = strings.TrimSuffix(pub.data, "\n")
	return pub, nil
}

// A reader which normalises line endings
// "/r" and "/r/n" are converted to "/n"
type normaliser struct {
	r        io.Reader
	lastChar byte
}

func newNormaliser(r io.Reader) *normaliser {
	return &normaliser{r: r}
}

func (norm *normaliser) Read(p []byte) (n int, err error) {
	n, err = norm.r.Read(p)
	for i := 0; i < n; i++ {
		switch {
		case p[i] == '\n' && norm.lastChar == '\r':
			copy(p[i:n], p[i+1:])
			norm.lastChar = p[i]
			n--
			i--
		case p[i] == '\r':
			norm.lastChar = p[i]
			p[i] = '\n'
		default:
			norm.lastChar = p[i]
		}
	}
	return
}
