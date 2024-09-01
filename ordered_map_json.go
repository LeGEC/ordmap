package ordmap

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

func (m *Map[K, V]) UnmarshalJSON(p []byte) error {
	buff := jsonBuff{p: p}

	// call '.peek()' once to make sure we "eat up" all leading space
	buff.peek()
	// if there is nothing to read: nothing to do
	if buff.eof() {
		return nil
	}
	if bytes.Equal(buff.tail(), []byte("null")) {
		m.m = nil
		m.keys = nil
		return nil
	}

	const (
		decKey = iota
		decValue
	)

	var (
		key   K
		value V

		zeroK K
		zeroV V
	)

	err := buff.Expect('{')
	if err != nil {
		return fmt.Errorf("error when decoding map: %w", err)
	}

	state := decKey

decodeLoop:
	for buff.more() {
		tok := buff.peek()
		if buff.eof() {
			return fmt.Errorf("error when decoding map: %w", io.ErrUnexpectedEOF)
		}

		switch state {
		case decKey:
			if tok != '"' {
				return fmt.Errorf("error when decoding map key: expected '\"', got '%c'", tok)
			}
			// reset variable used to decode 'key' value
			key = zeroK
			err := buff.Decode(&key)
			if err != nil {
				return fmt.Errorf("error when decoding key: %w", err)
			}

			err = buff.Expect(':')
			if err != nil {
				return fmt.Errorf("error when decoding key value pair: %w", err)
			}
			state = decValue

		case decValue:
			if !isValidFirstByte(tok) {
				return fmt.Errorf("error when decoding map value: expected a value start, got '%c'", tok)
			}

			value = zeroV
			err := buff.Decode(&value)
			if err != nil {
				return fmt.Errorf("error when decoding value: %w", err)
			}
			m.Set(key, value)

			tok = buff.peek()
			if tok == '}' {
				buff.move(1)
				break decodeLoop
			}
			if tok != ',' {
				return fmt.Errorf("error when decoding map: expected ',', got '%c'", tok)
			}
			buff.move(1)
			state = decKey
		}
	}

	if buff.peek() != 0 {
		return fmt.Errorf("error when decoding map: expected JSON payload to end at '}', got extra data")
	}

	return nil
}

func (m Map[K, V]) MarshalJSON() ([]byte, error) {
	if len(m.m) == 0 {
		return []byte("{}"), nil
	}

	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	buf.WriteByte('{')
	for i, k := range m.keys {
		if i > 0 {
			buf.WriteByte(',')
		}
		v, ok := m.m[k]
		if !ok {
			return nil, errors.New("invalid OrderedMap: a key is listed in Keys but is not mapped")
		}

		err := enc.Encode(k)
		if err != nil {
			return nil, fmt.Errorf("error when encoding key: %w", err)
		}
		buf.WriteByte(':')
		err = enc.Encode(v)
		if err != nil {
			return nil, fmt.Errorf("error when encoding key: %w", err)
		}
	}
	buf.WriteByte('}')

	return buf.Bytes(), nil
}
