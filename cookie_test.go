package rum

import (
	"encoding/hex"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var SecretKeyBase = os.Getenv("SECRET_KEY_BASE")

func TestKeyGenerator(t *testing.T) {
	rookie := New(SecretKeyBase)
	rookie.CookieSalt = []byte("example")

	expected_bytes, _ := hex.DecodeString("fc716f122bb1bafcc5ad060d2fe969568457879ed457e3256159d3d476112597a5fcad410a9c8f3d22364ce90aeb19b59b94fb7aeddf02925a37411d7f051828")

	assert.Equal(t, rookie.key(), expected_bytes)
}

func TestDecode(t *testing.T) {
	cookie := "WnpvU1NaQ2NSL3YvU2FGNWYxd0s4Q2xWSW1ZUm4vaTB2UENsRWQzWUFrYkZkemlGOFNxQXlxK3UyaWxBSUVEWUlkaFg2Si82dEdJb0dISWp6UXFHWDNJVlhwRW9ibEV1N2ZJbnNyVnhwVXBXc0hBVExaZVBCaVNzQVNUc3AwbG5IYmdPeE1XTDNQaytYMjJwOS9nSjZnU0h5SUMvRlVvNTZHdU1meU5lbjBUaGk4OFJmL2NmTDhOQWlJZVR1a3NWRzFPR3FmUzZvQk9HeS9RcXJDeHlVTVdxZHI2aHhUOUJZR2Q2bG1XdndqcmdhMXIway8vQVdCcEkvKzZ6Y3lyVTF0MVRYOGhDY0xFTVlzZkFBMUh4OGc9PS0tYzU3ZmpvUGdBYXVxN1JHbkVwYlkwUT09--6d6594c481b047a93bbb61f68df97e5e06c77641"

	var data interface{}

	rookie := New(SecretKeyBase)
	rookie.Decode(cookie, data)

	assert.NotEmpty(t, data)
}
