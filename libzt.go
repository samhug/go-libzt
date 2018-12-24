package libzt

/*
#cgo linux LDFLAGS: -L ${SRCDIR} -lzt_linux -lstdc++ -lm -std=c++11

#include "libzt.h"
#include <netdb.h>
*/
import "C"
import (
	"encoding/binary"
	"errors"
	"net"
	"strconv"
	"syscall"
	"unsafe"
)

type ZT struct {
	id       uint64
	homePath string
}

func Init(id string, homePath string) *ZT {
	_id, err := parseNWID(id)
	if err != nil {
		panic(err)
	}

	zt := &ZT{id: _id, homePath: homePath}
	C.zts_startjoin(C.CString(homePath), C.uint64_t(_id))
	return zt
}

func (zt *ZT) GetNumAssignedAddresses() int {
	return (int)(C.zts_get_num_assigned_addresses(C.uint64_t(zt.id)))
}

func (zt *ZT) GetAddressAtIndex(index int) (net.IP, error) {
	rsa := &syscall.RawSockaddrAny{}
	rsa_size := syscall.SizeofSockaddrAny

	ret := (int)(C.zts_get_address_at_index(C.uint64_t(zt.id), C.int(index), (*C.struct_sockaddr)(unsafe.Pointer(rsa)), (*C.uint)(unsafe.Pointer(&rsa_size))))
	if ret != 0 {
		return nil, errors.New("No address found")
	}

	switch rsa.Addr.Family {
	case syscall.AF_INET:
		pp := (*syscall.RawSockaddrInet4)(unsafe.Pointer(rsa))
		return net.IPv4(pp.Addr[0], pp.Addr[1], pp.Addr[2], pp.Addr[3]), nil

	case syscall.AF_INET6:
		pp := (*syscall.RawSockaddrInet6)(unsafe.Pointer(rsa))
		return net.IP{
			pp.Addr[0], pp.Addr[1], pp.Addr[2], pp.Addr[3],
			pp.Addr[4], pp.Addr[5], pp.Addr[6], pp.Addr[7],
			pp.Addr[8], pp.Addr[9], pp.Addr[10], pp.Addr[11],
			pp.Addr[12], pp.Addr[13], pp.Addr[14], pp.Addr[15],
		}, nil
	default:
		return nil, syscall.EAFNOSUPPORT
	}
}

func (zt *ZT) GetAddress(family int) (net.IP, error) {
	rsa := &syscall.RawSockaddrAny{}

	ret := (int)(C.zts_get_address(C.uint64_t(zt.id), (*C.struct_sockaddr_storage)(unsafe.Pointer(rsa)), C.int(family)))
	if ret != 0 {
		return nil, errors.New("No address found")
	}

	switch family {
	case syscall.AF_INET:
		pp := (*syscall.RawSockaddrInet4)(unsafe.Pointer(rsa))
		return net.IPv4(pp.Addr[0], pp.Addr[1], pp.Addr[2], pp.Addr[3]), nil

	case syscall.AF_INET6:
		pp := (*syscall.RawSockaddrInet6)(unsafe.Pointer(rsa))
		return net.IP{
			pp.Addr[0], pp.Addr[1], pp.Addr[2], pp.Addr[3],
			pp.Addr[4], pp.Addr[5], pp.Addr[6], pp.Addr[7],
			pp.Addr[8], pp.Addr[9], pp.Addr[10], pp.Addr[11],
			pp.Addr[12], pp.Addr[13], pp.Addr[14], pp.Addr[15],
		}, nil
	default:
		return nil, syscall.EAFNOSUPPORT
	}
}

func (zt *ZT) GetIPv4Address() (net.IP, error) {
	return zt.GetAddress(syscall.AF_INET)
}

func (zt *ZT) GetIPv6Address() (net.IP, error) {
	return zt.GetAddress(syscall.AF_INET6)
}

func (zt *ZT) Listen6(port uint16) (net.Listener, error) {
	fd := socket(syscall.AF_INET6, syscall.SOCK_STREAM, 0)
	if fd < 0 {
		return nil, errors.New("Error in opening socket")
	}

	serverSocket := syscall.RawSockaddrInet6{Flowinfo: 0, Family: syscall.AF_INET6, Port: htonl(port)}
	retVal := bind6(fd, serverSocket)
	if retVal < 0 {
		return nil, errors.New("ERROR on binding")
	}

	retVal = listen(fd, 1)
	if retVal < 0 {
		return nil, errors.New("ERROR listening")
	}

	ip, err := zt.GetIPv6Address()
	if err != nil {
		return nil, err
	}

	return &TCP6Listener{fd: fd, localIP: ip}, nil
}

func (zt *ZT) Listen6UDP(port uint16) (net.PacketConn, error) {
	fd := socket(syscall.AF_INET6, syscall.SOCK_DGRAM, syscall.IPPROTO_UDP)
	if fd < 0 {
		return nil, errors.New("Error in opening socket")
	}

	serverSocket := syscall.RawSockaddrInet6{Flowinfo: 0, Family: syscall.AF_INET6, Port: htonl(port)}
	retVal := bind6(fd, serverSocket)
	if retVal < 0 {
		return nil, errors.New("ERROR on binding")
	}

	ip, err := zt.GetIPv6Address()
	if err != nil {
		return nil, err
	}

	return &PacketConnection{fd: fd, localIP: ip, localPort: port}, nil
}

func (zt *ZT) Dial6UDP(ip string, port uint16) (net.Conn, error) {
	clientSocket := syscall.RawSockaddrInet6{Flowinfo: 0, Family: syscall.AF_INET6, Port: htonl(port), Addr: parseIPV6(ip)}

	fd := socket(30, 2, 17)
	if fd < 0 {
		return nil, errors.New("Error in opening socket")
	}

	localIP, err := zt.GetIPv6Address()
	if err != nil {
		return nil, err
	}

	conn := &Connection{
		fd:         fd,
		localIP:    localIP,
		localPort:  clientSocket.Port,
		remoteIp:   net.ParseIP(ip),
		remotePort: port,
	}
	return conn, nil
}

func (zt *ZT) Connect6(ip string, port uint16) (net.Conn, error) {
	clientSocket := syscall.RawSockaddrInet6{Flowinfo: 0, Family: syscall.AF_INET6, Port: htonl(port), Addr: parseIPV6(ip)}

	fd := socket(syscall.AF_INET6, syscall.SOCK_STREAM, 0)
	if fd < 0 {
		return nil, errors.New("Error in opening socket")
	}

	retVal := (int)(C.zts_connect(cint(fd), (*C.struct_sockaddr)(unsafe.Pointer(&clientSocket)), syscall.SizeofSockaddrInet6))
	if retVal < 0 {
		return nil, errors.New("Unable to connect")
	}

	localIP, err := zt.GetIPv6Address()
	if err != nil {
		return nil, err
	}

	conn := &Connection{
		fd:         fd,
		localIP:    localIP,
		localPort:  clientSocket.Port,
		remoteIp:   net.ParseIP(ip),
		remotePort: port,
	}
	return conn, nil
}

func htonl(number uint16) uint16 {
	bytes := make([]byte, 2)
	binary.BigEndian.PutUint16(bytes, number)
	return *(*uint16)(unsafe.Pointer(&bytes[0]))
}

func close(fd int) int {
	return (int)(C.zts_shutdown(cint(fd), 3))
}

func socket(family int, socketType int, protocol int) int {
	return (int)(C.zts_socket(cint(family), cint(socketType), cint(protocol)))
}

func listen(fd int, backlog int) int {
	return (int)(C.zts_listen(cint(fd), cint(backlog)))
}

func bind6(fd int, sockerAddr syscall.RawSockaddrInet6) int {
	return (int)(C.zts_bind(cint(fd), (*C.struct_sockaddr)(unsafe.Pointer(&sockerAddr)), syscall.SizeofSockaddrInet6))
}

func accept6(fd int) (int, syscall.RawSockaddrInet6) {
	socketAddr := syscall.RawSockaddrInet6{}
	socketLength := syscall.SizeofSockaddrInet6
	newFd := (int)(C.zts_accept(cint(fd), (*C.struct_sockaddr)(unsafe.Pointer(&socketAddr)), (*C.socklen_t)(unsafe.Pointer(&socketLength))))
	return newFd, socketAddr
}

func cint(value int) C.int {
	return (C.int)(value)
}

func parseIPV6(ipString string) [16]byte {
	ip := net.ParseIP(ipString)
	var arr [16]byte
	copy(arr[:], ip)
	return arr
}

func parseNWID(nwid string) (uint64, error) {
	if len(nwid) != 16 {
		return 0, errors.New("invalid nwid, must be 16 hex digits")
	}
	return strconv.ParseUint(nwid, 16, 64)
}
