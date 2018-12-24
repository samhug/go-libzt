package libzt

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"syscall"
	"time"
)

type PacketConnection struct {
	fd        int
	localIP   net.IP
	localPort uint16
}

func (c *PacketConnection) ReadFrom(b []byte) (int, net.Addr, error) {
	fmt.Println("ReadFrom")
	len, sa, err := syscall.Recvfrom(c.fd, b, 0)
	var fromAddr *net.UDPAddr

	switch sa := sa.(type) {
	case *syscall.SockaddrInet4:
		fromAddr = &net.UDPAddr{IP: sa.Addr[0:], Port: sa.Port}
	case *syscall.SockaddrInet6:
		fromAddr = &net.UDPAddr{IP: sa.Addr[0:], Port: sa.Port, Zone: string(sa.ZoneId)}
	}

	return len, toAddr(fromAddr.IP, uint16(fromAddr.Port)), err
}

func (c *PacketConnection) WriteTo(b []byte, addr net.Addr) (int, error) {
	udpAddr, ok := addr.(*net.UDPAddr)
	if !ok {
		return 0, &net.OpError{Op: "write", Net: "udp6", Source: c.LocalAddr(), Addr: addr, Err: syscall.EINVAL}
	}

	sa := toSocketAddr(udpAddr)

	err := syscall.Sendto(c.fd, b, 0, sa)

	if err != nil {
		err = &net.OpError{Op: "write", Net: "udp6", Source: c.LocalAddr(), Addr: addr, Err: err}
	}
	return len(b), err

}
func toSocketAddr(udpAddr *net.UDPAddr) *syscall.SockaddrInet6 {
	zoneId, _ := strconv.Atoi(udpAddr.Zone)
	sa := &syscall.SockaddrInet6{Port: udpAddr.Port, ZoneId: uint32(zoneId)}
	copy(sa.Addr[:], udpAddr.IP)
	return sa
}

func (c *PacketConnection) Close() error {
	val := close(c.fd)
	if val < 0 {
		return errors.New("Error closing socket")
	}
	return nil
}

func (c *PacketConnection) LocalAddr() net.Addr {
	return toAddr(c.localIP, c.localPort)
}
func toAddr(ip net.IP, port uint16) *net.UDPAddr {
	addr, _ := net.ResolveUDPAddr("udp6", fmt.Sprintf("[%s]:%d", ip.String(), port))
	return addr
}

func (c *PacketConnection) SetDeadline(t time.Time) error      { return nil }
func (c *PacketConnection) SetReadDeadline(t time.Time) error  { return nil }
func (c *PacketConnection) SetWriteDeadline(t time.Time) error { return nil }
