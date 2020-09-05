package read

import (
	"io/ioutil"
)

type fileRead struct{}

func NewFileRead() *fileRead {
	return &fileRead{}
}

func (f *fileRead) Read(fileName string) ([]byte, error) {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	return data, nil
}
