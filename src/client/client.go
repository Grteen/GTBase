package client

import (
	"GtBase/pkg/constants"
	"GtBase/utils"
	"bytes"
	"errors"
	"syscall"
)

type GtBaseClient struct {
	fd         int
	sendFd     int
	readBuffer []byte
	addr       *Address
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

func (c *GtBaseClient) WriteResp(data []byte) error {
	resp := make([]byte, 0, len(data)+4)
	resp = append(resp, utils.Encodeint32ToBytesSmallEnd(int32(len(data)))...)
	resp = append(resp, data...)
	resp = append(resp, []byte(constants.CommandSep)...)

	return c.Write(resp)
}

func (c *GtBaseClient) Dial() (int, error) {
	if c.sendFd != 0 {
		return -1, nil
	}
	sendFd, err := utils.Dial(c.addr.host, c.addr.port)
	c.sendFd = sendFd
	if err != nil {
		return -1, err
	}

	return sendFd, nil
}

func (c *GtBaseClient) GenerateKey() string {
	return c.addr.GenerateKey()
}

func CreateGtBaseClient(fd int, addr *Address) *GtBaseClient {
	return &GtBaseClient{fd: fd, readBuffer: make([]byte, 0), addr: addr}
}
