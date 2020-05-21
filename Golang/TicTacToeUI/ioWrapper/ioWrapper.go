package ioWrapper

import (
	"io/ioutil"
)

// ReadFile reads the file specified in the filepath parameter and outputs as a string its content
func ReadFile(filepath string) (string, error) {
	data, err := ioutil.ReadFile(filepath)

	if err != nil {
		return string(data), err
	}
	return string(data), nil
}
