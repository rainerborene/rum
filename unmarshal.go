package rum

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"reflect"
)

const (
	MARSHAL_MAJOR = 4
	MARSHAL_MINOR = 8
)

type Decoder struct {
	r *bufio.Reader
}

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{bufio.NewReader(r)}
}

func (dec *Decoder) integer() int {
	var result int
	value, _ := dec.r.ReadByte()
	c := int(value)

	if c == 0 {
		return 0
	} else if 5 < c && c < 128 {
		return c - 5
	} else if -129 < c && c < -5 {
		return c + 5
	}

	if c > 0 {
		result = 0
		for i := 0; i < c; i++ {
			n, _ := dec.r.ReadByte()
			result |= int(uint(n) << (8 * uint(i)))
		}
	} else {
		c = -c
		result = -1
		for i := 0; i < c; i++ {
			n, _ := dec.r.ReadByte()
			result &= ^(0xff << (8 * uint(i)))
			result |= int(uint(n) << (8 * uint(i)))
		}
	}

	return result
}

func (dec *Decoder) literal(v reflect.Value) {
	length := dec.integer()
	value := make([]byte, length)

	for ; length > 0; length-- {
		char, _ := dec.r.ReadByte()
		value = append(value, char)
	}

	v.SetString(string(value))
}

func (dec *Decoder) hash() {
	size := dec.integer()
	hash := make(map[string]interface{}, size)

	for i := 0; i < int(size); i++ {
		key := dec.unmarshal()
		value := dec.unmarshal()
		hash[key.(string)] = value
	}

	return hash
}

func (dec *Decoder) array(val reflect.Value) {
	size := dec.unmarshalRubyInteger()
	items := make([]interface{}, size)

	for i := 0; i < size; i++ {
		items[i] = dec.unmarshal(val)
	}

	return items
}

func (dec *Decoder) unmarshal(val reflect.Value) {
	rubyType, _ := dec.r.ReadByte()

	switch rubyType {
	case 'I', 'T', 'F', ';', 0x00, 0x06:
		dec.unmarshal(val)
	case '[':
		dec.array()
	case '{':
		hash := reflect.ValueOf(dec.unmarshalRubyHash())
		val.Set(hash)
	case '"':
		dec.literal()
	case 'i':
		dec.integer()
	case ':':
		encoding, _ := dec.r.Peek(3)
		if bytes.Equal(encoding, []byte{0x06, 'E', 'T'}) {
			dec.r.ReadByte()
			dec.r.ReadByte()
			dec.r.ReadByte()
			dec.unmarshal(val)
		}
		dec.literal()
	}
}

func (dec *Decoder) Decode(v interface{}) error {
	major, err := dec.r.ReadByte()
	minor, err := dec.r.ReadByte()

	if err != nil {
		return errors.New("unexpected end of stream")
	}

	if major != MARSHAL_MAJOR || minor > MARSHAL_MINOR {
		return errors.New("incompatible marshal file format")
	}

	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Ptr {
		return errors.New("non-pointer passed to Unmarshal")
	}

	dec.unmarshal(val.Elem())

	return nil
}

func Unmarshal(buf []byte, v interface{}) error {
	return NewDecoder(bytes.NewReader(buf)).Decode(v)
}
