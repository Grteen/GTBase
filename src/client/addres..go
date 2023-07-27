package client

import (
	"net"
	"strconv"
	"syscall"
)

type Address struct {
	host string
	port int
}

func (a *Address) GenerateKey() string {
	p := strconv.Itoa(a.port)
	return a.host + ":" + p
}

func CreateAddress(host string, port int) *Address {
	return &Address{host: host, port: port}
}

func GetAddressByFd(fd int) (*Address, error) {
	addr, err := syscall.Getsockname(int(fd))
	if err != nil {
		return nil, err
	}
	sockaddrIn := addr.(*syscall.SockaddrInet4)
	port := sockaddrIn.Port
	ip := net.IP(sockaddrIn.Addr[:])

	return CreateAddress(ip.String(), port), nil
}
