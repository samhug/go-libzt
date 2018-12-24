package libzt

import (
	"errors"
	"fmt"
	"net"
	"syscall"
	"time"
)

type Connection struct {
	fd        int
	localIP   net.IP
	localPort uint16

	remoteIp   net.IP
	remotePort uint16
}

func (c *Connection) Read(b []byte) (n int, err error) {
	return syscall.Read(c.fd, b)
}

func (c *Connection) Write(b []byte) (n int, err error) {
	return syscall.Write(c.fd, b)
}

func (c *Connection) Close() error {
	val := close(c.fd)
	if val < 0 {
		return errors.New("Error closing socket")
	}
	return nil
}

func (c *Connection) LocalAddr() net.Addr {
	addr, _ := net.ResolveTCPAddr("tcp6", fmt.Sprintf("[%s]:%d", c.localIP.String(), c.localPort))
	return addr
}

func (c *Connection) RemoteAddr() net.Addr {
	addr, _ := net.ResolveTCPAddr("tcp6", fmt.Sprintf("[%s]:%d", c.remoteIp.String(), c.remotePort))
	return addr
}

func (c *Connection) SetDeadline(time.Time) error      { return errors.New("Not yet supported") }
func (c *Connection) SetReadDeadline(time.Time) error  { return errors.New("Not yet supported") }
func (c *Connection) SetWriteDeadline(time.Time) error { return errors.New("Not yet supported") }
