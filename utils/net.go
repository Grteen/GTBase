package utils

import (
	"net"
	"syscall"
)

func BindAndListen(port int) (int32, error) {
	listenSock, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if err != nil {
		return -1, err
	}

	addr := syscall.SockaddrInet4{Port: port}
	copy(addr.Addr[:], net.ParseIP("127.0.0.1").To4())

	err = syscall.Bind(listenSock, &addr)
	if err != nil {
		return -1, err
	}

	syscall.Listen(listenSock, 0)
	return int32(listenSock), nil
}

func Accept(listenFd int) (int, error) {
	nfd, _, err := syscall.Accept(listenFd)
	if err != nil {
		return -1, err
	}

	return nfd, nil
}

func LocalDial(port int) (int, error) {
	sockfd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if err != nil {
		return -1, err
	}

	serverAddr := syscall.SockaddrInet4{
		Port: port,
		Addr: [4]byte{127, 0, 0, 1},
	}

	errc := syscall.Connect(sockfd, &serverAddr)
	if errc != nil {
		return -1, errc
	}

	return sockfd, nil
}

func Dial(host string, port int) (int, error) {
	sockfd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if err != nil {
		return -1, err
	}

	ipAddr, errr := net.ResolveIPAddr("ip", host)
	if errr != nil {
		return -1, errr
	}

	ip := ipAddr.IP.To4()

	serverAddr := syscall.SockaddrInet4{
		Port: port,
		Addr: [4]byte(ip),
	}

	errc := syscall.Connect(sockfd, &serverAddr)
	if errc != nil {
		return -1, errc
	}

	return sockfd, nil
}

func ReadFd(fd int) ([]byte, error) {
	buf := make([]byte, 1024)
	n, err := syscall.Read(fd, buf)
	if err != nil {
		return nil, err
	}

	return buf[:n], nil
}

func CloseFd(fd int) {
	syscall.Close(fd)
}
