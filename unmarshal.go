package rookie

import (
	"bufio"
	"bytes"
	"errors"
)

var encodingSymbols = []byte{0x06, 'E', 'T'}

type decoder struct {
	r *bufio.Reader
}

func (d *decoder) unmarshalRubyInteger() int {
	var result int
	value, _ := d.r.ReadByte()
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
			n, _ := d.r.ReadByte()
			result |= int(uint(n) << (8 * uint(i)))
		}
	} else {
		c = -c
		result = -1
		for i := 0; i < c; i++ {
			n, _ := d.r.ReadByte()
			result &= ^(0xff << (8 * uint(i)))
			result |= int(uint(n) << (8 * uint(i)))
		}
	}

	return result
}

func (d *decoder) unmarshalRubyString() string {
	var value []byte
	length := int(d.unmarshalRubyInteger())

	for ; length > 0; length-- {
		char, _ := d.r.ReadByte()
		value = append(value, char)
	}

	return string(value)
}

func (d *decoder) unmarshalRubyHash() interface{} {
	size := d.unmarshalRubyInteger()
	hash := make(map[string]interface{}, int(size))

	for i := 0; i < int(size); i++ {
		key := d.unmarshalRubyType()
		value := d.unmarshalRubyType()
		hash[key.(string)] = value
	}

	return hash
}

func (d *decoder) unmarshalRubyArray() []interface{} {
	size := d.unmarshalRubyInteger()
	items := make([]interface{}, int(size))

	for i := 0; i < int(size); i++ {
		items[i] = d.unmarshalRubyType()
	}

	return items
}

func (d *decoder) unmarshalRubyType() interface{} {
	rubyType, _ := d.r.ReadByte()

	switch rubyType {
	case 'I', 'T', 'F', ';', 0x00, 0x06:
		return d.unmarshalRubyType()
	case '[':
		return d.unmarshalRubyArray()
	case '{':
		return d.unmarshalRubyHash()
	case '"':
		return d.unmarshalRubyString()
	case 'i':
		return d.unmarshalRubyInteger()
	case ':':
		encoding, _ := d.r.Peek(3)
		if bytes.Equal(encoding, encodingSymbols) {
			d.r.ReadByte()
			d.r.ReadByte()
			d.r.ReadByte()
			return d.unmarshalRubyType()
		}
		return d.unmarshalRubyString()
	}

	return nil
}

func Unmarshal(buf []byte) (interface{}, error) {
	decoder := &decoder{bufio.NewReader(bytes.NewReader(buf))}
	major, err := decoder.r.ReadByte()
	minor, err := decoder.r.ReadByte()

	if err != nil {
		return nil, errors.New("unexpected end of stream")
	}

	if major != 4 || minor > 8 {
		return nil, errors.New("incompatible marshal file format")
	}

	return decoder.unmarshalRubyType(), nil
}
