package transport_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vysina/linebotgo/transport"
)

func TestThriftReaderWriter(t *testing.T) {
	buf := &bytes.Buffer{}
	w := transport.NewThriftWriter(buf)
	err := w.WriteString("hello")
	assert.NoError(t, err)

	r := transport.NewThriftReader(bytes.NewReader(buf.Bytes()))
	got, err := r.ReadString()
	assert.NoError(t, err)
	assert.Equal(t, "hello", got)
}
