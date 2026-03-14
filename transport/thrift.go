package transport

import (
	"context"
	"io"

	"github.com/apache/thrift/lib/go/thrift"
)

// ThriftWriter wraps Thrift binary protocol for writing.
type ThriftWriter struct {
	proto thrift.TProtocol
}

func NewThriftWriter(w io.Writer) *ThriftWriter {
	trans := thrift.NewStreamTransportW(w)
	proto := thrift.NewTBinaryProtocolTransport(trans)
	return &ThriftWriter{proto: proto}
}

func (w *ThriftWriter) WriteString(s string) error {
	if err := w.proto.WriteString(context.Background(), s); err != nil {
		return err
	}
	return w.proto.Flush(context.Background())
}

// ThriftReader wraps Thrift binary protocol for reading.
type ThriftReader struct {
	proto thrift.TProtocol
}

func NewThriftReader(r io.Reader) *ThriftReader {
	trans := thrift.NewStreamTransportR(r)
	proto := thrift.NewTBinaryProtocolTransport(trans)
	return &ThriftReader{proto: proto}
}

func (r *ThriftReader) ReadString() (string, error) {
	return r.proto.ReadString(context.Background())
}
