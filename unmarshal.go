package rookie

import (
	"bufio"
	"bytes"
	"errors"
)

var encodingSymbols = []byte{0x06, 'E', 'T'}

func unmarshalRubyInteger(inputStream *bufio.Reader) int {
	var result int
	value, _ := inputStream.ReadByte()
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
			n, _ := inputStream.ReadByte()
			result |= int(uint(n) << (8 * uint(i)))
		}
	} else {
		c = -c
		result = -1
		for i := 0; i < c; i++ {
			n, _ := inputStream.ReadByte()
			result &= ^(0xff << (8 * uint(i)))
			result |= int(uint(n) << (8 * uint(i)))
		}
	}

	return result
}

func unmarshalRubyString(inputStream *bufio.Reader) string {
	var value []byte
	length := int(unmarshalRubyInteger(inputStream))

	for ; length > 0; length-- {
		char, _ := inputStream.ReadByte()
		value = append(value, char)
	}

	return string(value)
}

func unmarshalRubyHash(inputStream *bufio.Reader) interface{} {
	size := unmarshalRubyInteger(inputStream)
	hash := make(map[string]interface{}, int(size))

	for i := 0; i < int(size); i++ {
		key := unmarshalRubyType(inputStream)
		value := unmarshalRubyType(inputStream)
		hash[key.(string)] = value
	}

	return hash
}

func unmarshalRubyArray(inputStream *bufio.Reader) []interface{} {
	size := unmarshalRubyInteger(inputStream)
	items := make([]interface{}, int(size))

	for i := 0; i < int(size); i++ {
		items[i] = unmarshalRubyType(inputStream)
	}

	return items
}

func unmarshalRubyType(inputStream *bufio.Reader) interface{} {
	rubyType, _ := inputStream.ReadByte()

	switch rubyType {
	case 'I', 'T', 'F', ';', 0x00, 0x06:
		return unmarshalRubyType(inputStream)
	case '[':
		return unmarshalRubyArray(inputStream)
	case '{':
		return unmarshalRubyHash(inputStream)
	case '"':
		return unmarshalRubyString(inputStream)
	case 'i':
		return unmarshalRubyInteger(inputStream)
	case ':':
		encoding, _ := inputStream.Peek(3)
		if bytes.Equal(encoding, encodingSymbols) {
			inputStream.ReadByte()
			inputStream.ReadByte()
			inputStream.ReadByte()
			return unmarshalRubyType(inputStream)
		}
		return unmarshalRubyString(inputStream)
	}

	return nil
}

func Unmarshal(buf []byte) (interface{}, error) {
	inputStream := bufio.NewReader(bytes.NewReader(buf))
	major, err := inputStream.ReadByte()
	minor, err := inputStream.ReadByte()

	if err != nil {
		return nil, errors.New("unexpected end of stream")
	}

	if major != 4 || minor > 8 {
		return nil, errors.New("incompatible marshal file format")
	}

	return unmarshalRubyType(inputStream), nil
}
