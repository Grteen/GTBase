package cmd

import "GtBase/src/object"

func WriteAndRead(req []byte, c *GtBaseClient) object.Object {
	errw := c.writeToGtBase(c.fd, req)
	if errw != nil {
		return object.CreateGtString(errw.Error())
	}

	result, errr := c.readFromGtBase(c.fd)
	if errr != nil {
		return object.CreateGtString(errr.Error())
	}

	return object.CreateGtString(string(result))
}
