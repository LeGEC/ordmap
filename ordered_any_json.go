package ordmap

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

func (x Any) MarshalJSON() ([]byte, error) {
	if x.v == nil {
		return []byte("null"), nil
	}
	return json.Marshal(x.v)
}

func (x *Any) UnmarshalJSON(p []byte) error {
	buff := jsonBuff{p: p}

	x.v = nil

	tok := buff.peek()
	cc := fmt.Sprintf("%c", tok)
	_ = cc
	if buff.eof() {
		return fmt.Errorf("error when decoding any: %w", io.ErrUnexpectedEOF)
	}
	if bytes.Equal(buff.tail(), []byte("null")) {
		return nil
	}
	if !isValidFirstByte(tok) {
		return fmt.Errorf("error decoding value: expected a value start, got '%c'", tok)
	}

	var v any
	var err error
	switch tok {
	case '{':
		v, err = jsonUnmarshalAnyObject(&buff)
	case '[':
		v, err = jsonUnmarshalAnyArray(&buff)
	default:
		err = buff.Decode(&v)
	}

	if err != nil {
		return err
	}
	buff.eatSpace()
	if !buff.eof() {
		return fmt.Errorf("error decoding value: trailing data after value")
	}
	x.v = v
	return nil
}

func jsonUnmarshalAnyObject(buff *jsonBuff) (*Map[string, any], error) {
	var key string

	err := buff.Expect('{')
	if err != nil { // note: will not happen, we reach here because we come from .UnmarshalJSON
		return nil, err
	}

	var m Map[string, any]

	// check for empty object case:
	if buff.more() && buff.peek() == '}' {
		buff.move(1)
		return &m, nil
	}

	const (
		stateDecodeKey = iota
		stateDecodevalue
		stateObjectCompleted
	)

	state := stateDecodeKey

decodeLoop:
	for buff.more() {
		tok := buff.peek()

		cc := fmt.Sprintf("%c", tok)
		_ = cc

		switch state {
		case stateDecodeKey:
			switch {
			case buff.eof():
				return nil, fmt.Errorf("error when decoding object key: %w", io.ErrUnexpectedEOF)
			case tok != '"':
				return nil, fmt.Errorf("error when decoding object key: expected '\"', got '%c'", tok)
			}

			key = ""
			err = buff.Decode(&key)
			if err != nil {
				return nil, fmt.Errorf("error when decoding object key: %w", err)
			}

			err = buff.Expect(':')
			if err != nil {
				return nil, fmt.Errorf("error when decoding object kv pair: %w", err)
			}

			state = stateDecodevalue

		case stateDecodevalue:
			switch {
			case buff.eof():
				return nil, fmt.Errorf("error when decoding object value: %w", io.ErrUnexpectedEOF)
			case !isValidFirstByte(tok):
				return nil, fmt.Errorf("error when decoding object value: expected a value start, got '%c'", tok)
			}

			var v any
			var err error
			switch tok {
			case '{':
				v, err = jsonUnmarshalAnyObject(buff)
			case '[':
				v, err = jsonUnmarshalAnyArray(buff)
			default:
				err = buff.Decode(&v)
			}

			if err != nil {
				return nil, err
			}
			m.Set(key, v)

			tok = buff.peek()
			if tok == '}' {
				state = stateObjectCompleted
				buff.move(1)
				break decodeLoop
			}
			if tok != ',' {
				return nil, fmt.Errorf("error when decoding object: expected ',', got '%c'", tok)
			}
			buff.move(1)
			state = stateDecodeKey
		}
	}

	if state != stateObjectCompleted {
		return nil, fmt.Errorf("error when decoding object: %w", io.ErrUnexpectedEOF)
	}

	return &m, nil
}

func jsonUnmarshalAnyArray(buff *jsonBuff) ([]any, error) {
	var res = make([]any, 0)

	err := buff.Expect('[')
	if err != nil { // note: will not happen, we reach here because we come from .UnmarshalJSON
		return nil, err
	}

	firstValue := true
	arrayCompleted := false

decodeLoop:
	for buff.more() {
		tok := buff.peek()
		if tok == ']' {
			arrayCompleted = true
			buff.move(1)
			break decodeLoop
		}

		if !firstValue {
			err = buff.Expect(',')
			if err != nil {
				return nil, fmt.Errorf("error when decoding array: %w", err)
			}
			tok = buff.peek()
		}

		cc := fmt.Sprintf("%c", tok)
		_ = cc

		switch {
		case buff.eof(): // we are waiting for a value or for array end, EOF is unexpected
			return nil, fmt.Errorf("error when decoding array: %w", io.ErrUnexpectedEOF)
		case !isValidFirstByte(tok):
			return nil, fmt.Errorf("error when decoding array: expected a value start, got '%c'", tok)
		}

		var v any
		var err error
		if tok == '{' {
			v, err = jsonUnmarshalAnyObject(buff)
		} else if tok == '[' {
			v, err = jsonUnmarshalAnyArray(buff)
		} else {
			err = buff.Decode(&v)
		}

		if err != nil {
			return nil, err
		}

		firstValue = false
		res = append(res, v)
	}

	if !arrayCompleted {
		return nil, fmt.Errorf("error when decoding array: %w", io.ErrUnexpectedEOF)
	}

	return res, nil
}
