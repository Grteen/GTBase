package client

import (
	"GtBase/pkg/constants"
	"bytes"
	"errors"
	"syscall"
)

type GtBaseClient struct {
	fd         int
	readBuffer []byte
}

func (c *GtBaseClient) GetFd() int {
	return c.fd
}

func (c *GtBaseClient) Read() ([]byte, error) {
	err := c.read()
	if err != nil {
		return nil, err
	}

	idx := bytes.Index(c.readBuffer, []byte(constants.CommandSep))
	if idx != -1 {
		result := c.readBuffer[:idx]
		c.readBuffer = c.readBuffer[idx+len(constants.CommandSep):]
		return result, nil
	}

	return nil, nil
}

func (c *GtBaseClient) read() error {
	result := make([]byte, 1024)
	n, err := syscall.Read(c.fd, result)
	if err != nil {
		return err
	}

	if n == 0 {
		return errors.New(constants.ClientExitError)
	}

	c.readBuffer = append(c.readBuffer, result[:n]...)
	return nil
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
	return &GtBaseClient{fd: fd, readBuffer: make([]byte, 0)}
}
