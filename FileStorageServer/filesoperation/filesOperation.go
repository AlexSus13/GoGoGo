package filesoperation

import (
	"github.com/pkg/errors"

	"io/ioutil"
	"os"
)

func Get(path string, FN string) ([]byte, error) {

	file, err := os.Open(path + FN)
	if err != nil {
		return nil, errors.Wrap(err, "Open File, func Get") //error handling
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, errors.Wrap(err, "Read File, func Get") //error handling
	}

	return data, nil
}

func SaveNew(path string, FN string, data []byte)  error {

	file, err := os.Create(path + FN)
	if err != nil {
		return errors.Wrap(err, "Create File, func SaveNew") //error handling
	}
	defer file.Close()

	file.Write(data)

	return nil
}

func DeleteOld(path string, FN string) error {

	err := os.Remove(path + FN)
	if err != nil {
		return errors.Wrap(err, "Remove File, func DeleteOld") //error handling
	}

	return nil
}
