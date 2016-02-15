// Copyright 2016 hit9. All rights reserved.

package reuseport

import (
	"errors"
	"fmt"
	"net"
	"os"
	"syscall"
)

var ErrProtocol = errors.New("Only tcp,tcp4,tcp6,udp,udp4,udp6 are supported")
var ErrProtocolTCP = errors.New("Only tcp,tcp4,tcp6 are supported")
var ErrProtocolUDP = errors.New("Only udp,udp4,udp6 are supported")

// FileListener fileName prefix.
const fileNamePrefix = "port."

func getSockaddr(proto, addr string) (syscall.Sockaddr, int, error) {
	switch proto {
	case "tcp", "tcp4", "tcp6":
		// TCP
		tcpAddr, err := net.ResolveTCPAddr(proto, addr)
		if err != nil {
			return nil, 0, err
		}
		switch proto {
		case "tcp", "tcp4":
			// TCP V4
			var addr4 [4]byte
			if tcpAddr.IP != nil {
				copy(addr4[:], tcpAddr.IP[12:16])
			}
			return &syscall.SockaddrInet4{Port: tcpAddr.Port, Addr: addr4}, syscall.AF_INET, nil
		case "tcp6":
			// TCP V6
			var addr6 [16]byte
			if tcpAddr.IP != nil {
				copy(addr6[:], tcpAddr.IP)
			}
			return &syscall.SockaddrInet6{Port: tcpAddr.Port, Addr: addr6}, syscall.AF_INET6, nil
		}
	case "udp", "udp4", "udp6":
		// UDP
		udpAddr, err := net.ResolveUDPAddr(proto, addr)
		if err != nil {
			return nil, 0, err
		}
		switch proto {
		case "udp", "udp4":
			// UDP V4
			var addr4 [4]byte
			if udpAddr.IP != nil {
				copy(addr4[:], udpAddr.IP[12:16])
			}
			return &syscall.SockaddrInet4{Port: udpAddr.Port, Addr: addr4}, syscall.AF_INET, nil
		case "udp6":
			// UDP V6
			var addr6 [16]byte
			if udpAddr.IP != nil {
				copy(addr6[:], udpAddr.IP)
			}
			return &syscall.SockaddrInet6{Port: udpAddr.Port, Addr: addr6}, syscall.AF_INET6, nil
		}
	}
	return nil, 0, ErrProtocol
}

func Listener(proto, addr string) (net.Listener, error) {
	switch proto {
	case "tcp", "tcp4", "tcp6":
	default:
		return nil, ErrProtocolTCP
	}
	// Get socket address.
	sockAddr, sockType, err := getSockaddr(proto, addr)
	if err != nil {
		return nil, err
	}
	// New socket.
	fd, err := syscall.Socket(sockType, syscall.SOCK_STREAM, syscall.IPPROTO_TCP)
	if err != nil {
		return nil, err
	}
	// Set socket option.
	err = syscall.SetsockoptInt(fd, syscall.SOL_SOCKET, OPT_REUSEPORT, 1)
	if err != nil {
		return nil, err
	}
	// Bind with address.
	err = syscall.Bind(fd, sockAddr)
	if err != nil {
		return nil, err
	}
	// Set backlog size to the maximum.
	err = syscall.Listen(fd, syscall.SOMAXCONN)
	if err != nil {
		return nil, err
	}
	// Create listener.
	fileName := fmt.Sprintf("%s%d", fileNamePrefix, os.Getpid())
	file := os.NewFile(uintptr(fd), fileName)
	ln, err := net.FileListener(file)
	if err != nil {
		return nil, err
	}
	err = file.Close()
	if err != nil {
		return nil, err
	}
	return ln, nil
}

func PacketConn(proto, addr string) (net.PacketConn, error) {
	switch proto {
	case "udp", "udp4", "udp6":
	default:
		return nil, ErrProtocolUDP
	}
	// Get socket address.
	sockAddr, sockType, err := getSockaddr(proto, addr)
	if err != nil {
		return nil, err
	}
	// New socket.
	fd, err := syscall.Socket(sockType, syscall.SOCK_DGRAM, syscall.IPPROTO_UDP)
	if err != nil {
		return nil, err
	}
	// Set socket option.
	err = syscall.SetsockoptInt(fd, syscall.SOL_SOCKET, OPT_REUSEPORT, 1)
	if err != nil {
		return nil, err
	}
	// Bind with address.
	err = syscall.Bind(fd, sockAddr)
	if err != nil {
		return nil, err
	}
	// Create conn.
	fileName := fmt.Sprintf("%s%d", fileNamePrefix, os.Getpid())
	file := os.NewFile(uintptr(fd), fileName)
	conn, err := net.FilePacketConn(file)
	if err != nil {
		return nil, err
	}
	err = file.Close()
	if err != nil {
		return nil, err
	}
	return conn, nil
}
