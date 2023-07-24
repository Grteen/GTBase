package server

import (
	"syscall"
)

type GtBaseClient struct {
	fd int
}

func (c *GtBaseClient) GetFd() int {
	return c.fd
}

func (c *GtBaseClient) Read() ([]byte, error) {
	result := make([]byte, 1024)
	n, err := syscall.Read(c.fd, result)
	if err != nil {
		return nil, err
	}

	return result[:n], err
}

func (c *GtBaseClient) Write(data []byte) error {
	sendData := data
	for len(sendData) != 0 {
		n, err := syscall.Write(c.fd, sendData)
		if err != nil {
			return err
		}
		sendData = sendData[n:]
	}

	return nil
}

func CreateGtBaseClient(fd int) *GtBaseClient {
	return &GtBaseClient{fd: fd}
}
