package headers 

import (
	"testing"
	"io"


	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)


type chunkReader struct {
	data            string
	numBytesPerRead int
	pos             int
}

// Read reads up to len(p) or numBytesPerRead bytes from the string per call
// its useful for simulating reading a variable number of bytes per chunk from a network connection
func (cr *chunkReader) Read(p []byte) (n int, err error) {
	if cr.pos >= len(cr.data) {
		return 0, io.EOF
	}
	endIndex := cr.pos + cr.numBytesPerRead
	if endIndex > len(cr.data) {
		endIndex = len(cr.data)
	}
	n = copy(p, cr.data[cr.pos:endIndex])
	cr.pos += n
	if n > cr.numBytesPerRead {
		n = cr.numBytesPerRead
		cr.pos -= n - cr.numBytesPerRead
	}
	return n, nil
}

func TestParse(t *testing.T) {

// Test: Valid single header
headers := NewHeaders()
data := []byte("Host: localhost:42069\r\n\r\n")
n, done, err := headers.Parse(data)
require.NoError(t, err)
require.NotNil(t, headers)
assert.Equal(t, "localhost:42069", headers["host"])
assert.Equal(t, 23, n)
assert.False(t, done)


// Test: Valid single header with extra whitespace
headers = NewHeaders()
data = []byte("       Host: localhost:42069\r\n\r\n")
n, done, err = headers.Parse(data)
require.NoError(t, err)
require.NotNil(t, headers)
assert.Equal(t, "localhost:42069", headers["host"])
assert.Equal(t, 30, n)
assert.False(t, done)


// Test: Valid done
headers = NewHeaders()
data = []byte("\r\n\r\n")
n, done, err = headers.Parse(data)
require.NoError(t, err)
require.NotNil(t, headers)
// assert.Equal(t, "localhost:42069", headers["Host"])
assert.Equal(t, 2, n)
assert.True(t, done)


// Test: Invalid spacing header
headers = NewHeaders()
data = []byte("       Host : localhost:42069       \r\n\r\n")
n, done, err = headers.Parse(data)
require.Error(t, err)
assert.Equal(t, 0, n)
assert.False(t, done)

//Test: lower and upper case headers, should only produce one entry in map
headers = NewHeaders()
data = []byte("Accept: Contents\r\n\r\n")
n, done, err = headers.Parse(data)
require.NoError(t, err)
assert.NotEqual(t, "Contents",headers["Accept"])
assert.Equal(t, "Contents",headers["accept"])
assert.Equal(t, 18, n)
assert.False(t, done)


//Test: Invalid Characters in Header Key
headers = map[string]string{"Accept": "Contents"}
data = []byte("@ccept: contents\r\n\r\n")
n, done, err = headers.Parse(data)
require.Error(t, err)
assert.Equal(t, 0, n)
assert.False(t, done)


//Test: add additional entries to a map
headers = map[string]string{"content-type":"text/html"}
data = []byte("Content-Type: json\r\n\r\n")
n, done, err = headers.Parse(data)
require.NoError(t, err)
assert.Equal(t, "text/html, json",headers["content-type"])
assert.Equal(t, 20, n)
assert.False(t, done)




}