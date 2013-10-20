package rookie

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var SecretKeyBase = os.Getenv("SECRET_KEY_BASE")

func TestKeyGenerator(t *testing.T) {
	rookie := New(SecretKeyBase)
	rookie.CookieSalt = []byte("example")

	expected_bytes := []byte{
		252, 113, 111, 18, 43, 177, 186, 252, 197, 173, 6, 13, 47, 233, 105, 86,
		132, 87, 135, 158, 212, 87, 227, 37, 97, 89, 211, 212, 118, 17, 37, 151,
		165, 252, 173, 65, 10, 156, 143, 61, 34, 54, 76, 233, 10, 235, 25, 181,
		155, 148, 251, 122, 237, 223, 2, 146, 90, 55, 65, 29, 127, 5, 24, 40,
	}

	encrypted_key := rookie.generateKey()

	assert.Equal(t, encrypted_key, expected_bytes)
}

func TestDecode(t *testing.T) {
	cookie := "WnpvU1NaQ2NSL3YvU2FGNWYxd0s4Q2xWSW1ZUm4vaTB2UENsRWQzWUFrYkZkemlGOFNxQXlxK3UyaWxBSUVEWUlkaFg2Si82dEdJb0dISWp6UXFHWDNJVlhwRW9ibEV1N2ZJbnNyVnhwVXBXc0hBVExaZVBCaVNzQVNUc3AwbG5IYmdPeE1XTDNQaytYMjJwOS9nSjZnU0h5SUMvRlVvNTZHdU1meU5lbjBUaGk4OFJmL2NmTDhOQWlJZVR1a3NWRzFPR3FmUzZvQk9HeS9RcXJDeHlVTVdxZHI2aHhUOUJZR2Q2bG1XdndqcmdhMXIway8vQVdCcEkvKzZ6Y3lyVTF0MVRYOGhDY0xFTVlzZkFBMUh4OGc9PS0tYzU3ZmpvUGdBYXVxN1JHbkVwYlkwUT09--6d6594c481b047a93bbb61f68df97e5e06c77641"

	rookie := New(SecretKeyBase)
	data, _ := rookie.Decode(cookie)
	assert.Contains(t, string(data), "session_id")
}
