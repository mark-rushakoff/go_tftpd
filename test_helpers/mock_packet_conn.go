package test_helpers

import (
	"bytes"
	"encoding/binary"
	"net"
	"testing"
	"time"
)

type MockPacketConn struct {
	ReadFromFunc         func([]byte) (int, net.Addr, error)
	WriteToFunc          func([]byte, net.Addr) (int, error)
	CloseFunc            func() error
	LocalAddrFunc        func() net.Addr
	SetDeadlineFunc      func(time.Time) error
	SetReadDeadlineFunc  func(time.Time) error
	SetWriteDeadlineFunc func(time.Time) error

	lastPacketOut struct {
		packet []byte
		addr   net.Addr
	}
	sentAnyPackets bool
}

func (c *MockPacketConn) ReadFrom(p []byte) (n int, addr net.Addr, err error) {
	return c.ReadFromFunc(p)
}

func (c *MockPacketConn) WriteTo(b []byte, addr net.Addr) (n int, err error) {
	c.lastPacketOut.packet = b
	c.lastPacketOut.addr = addr
	c.sentAnyPackets = true

	return c.WriteToFunc(b, addr)
}

func (c *MockPacketConn) Close() error {
	return c.CloseFunc()
}
func (c *MockPacketConn) LocalAddr() net.Addr {
	return c.LocalAddrFunc()
}

func (c *MockPacketConn) SetDeadline(t time.Time) error {
	return c.SetDeadlineFunc(t)
}
func (c *MockPacketConn) SetReadDeadline(t time.Time) error {
	return c.SetReadDeadlineFunc(t)
}
func (c *MockPacketConn) SetWriteDeadline(t time.Time) error {
	return c.SetWriteDeadlineFunc(t)
}

func (c *MockPacketConn) LastPacketOut() (packet []byte, addr net.Addr, ok bool) {
	packet = c.lastPacketOut.packet
	addr = c.lastPacketOut.addr
	ok = c.sentAnyPackets
	return
}

func NewMockPacketConnWithBytes(t *testing.T, addr net.Addr, data []interface{}) *MockPacketConn {
	wasCalledOnce := false
	return &MockPacketConn{
		ReadFromFunc: func(b []byte) (int, net.Addr, error) {
			if wasCalledOnce {
				t.Fatalf("Called fake reader more than once")
			}
			wasCalledOnce = true

			buf := new(bytes.Buffer)
			for _, v := range data {
				str, isString := v.(string)
				if isString {
					_, err := buf.WriteString(str)
					if err != nil {
						t.Errorf("Failed to write string to buffer")
					}
				} else {
					err := binary.Write(buf, binary.BigEndian, v)
					if err != nil {
						t.Errorf("Failed to write data to buffer")
					}
				}
			}
			n := copy(b, buf.Bytes())
			return n, addr, nil
		},
	}
}
