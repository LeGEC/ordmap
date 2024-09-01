package ordmap

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

// isSpace
//
// check if a byte represents one of the valid space characters in JSON
func isSpace(b byte) bool {
	return b == ' ' || b == '\t' || b == '\n' || b == '\r'
}

// isValidFirstByte
//
// returns true iff 'b' is valid to be the leading byte of a new json value.
func isValidFirstByte(b byte) bool {
	switch b {
	case 'n', 't', 'f', // litteral: null, true or false
		'-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', // number
		'"', '[', '{': // string, array or object
		return true

	default:
		return false
	}
}

// jsonBuff
// utility wrapper on top of a []byte, used in `UmarshalJSON(b []byte)` to parse JSON values.
type jsonBuff struct {
	p []byte
	i int
}

func (b *jsonBuff) eatSpace() {
	for b.i < len(b.p) && isSpace(b.p[b.i]) {
		b.i++
	}
}

func (b *jsonBuff) peek() byte {
	b.eatSpace()
	if b.i >= len(b.p) {
		return 0
	}
	return b.p[b.i]
}

func (b *jsonBuff) eof() bool {
	return b.i >= len(b.p)
}

func (b *jsonBuff) more() bool {
	b.eatSpace()
	return b.i < len(b.p)
}

func (b *jsonBuff) tail() []byte {
	if b.i >= len(b.p) {
		return nil
	}
	return b.p[b.i:]
}

func (b *jsonBuff) move(x int) error {
	if b.i+x < 0 {
		b.i = 0
		return errors.New("negative move before beginning of buffer")
	}
	if b.i+x > len(b.p) {
		b.i = len(b.p)
		return errors.New("move past end of buffer")
	}
	b.i += x
	return nil
}

func (b *jsonBuff) Decode(v any) error {
	dec := json.NewDecoder(bytes.NewReader(b.tail()))
	err := dec.Decode(v)
	if err != nil {
		return err
	}
	err = b.move(int(dec.InputOffset()))
	if err != nil {
		return err
	}

	return nil
}

func (b *jsonBuff) Expect(expectedTok byte) error {
	got := b.peek()
	if b.eof() {
		return io.ErrUnexpectedEOF
	}
	if got != expectedTok {
		return fmt.Errorf("expected '%c', but got '%c'", expectedTok, got)
	}
	b.move(1)
	return nil
}
