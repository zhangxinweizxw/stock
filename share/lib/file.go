package lib

import (
    "errors"
	"io/ioutil"
	"os"
	"path/filepath"
)

func ReadFile(path string) (string, error) {
	file, err := os.OpenFile(path, os.O_RDONLY, os.ModeType)
	if err != nil {
		return "", err
	}
	defer file.Close()

	b, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

// open
func OpenFile(fileName string) (*os.File, error) {
	_, err := os.Stat(fileName)

	var file *os.File
	if os.IsNotExist(err) {
		_, err = os.Create(fileName)
	} else {
		_, err = os.OpenFile(fileName, os.O_TRUNC, 0666)
	}
	if err != nil {
		return nil, err
	}

	file, err = os.Create(fileName)
	return file, err
}

// wrrite file
func WriteFile(fileName string, data string) error {
	if CheckDir(filepath.Dir(fileName)) == false {
		return errors.New("Create Dir Error")
	}

	file, err := OpenFile(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(data)
	return err
}

// delete file
func DeleteFile(fileName string) error {
	if IsFileExist(fileName) == false {
		return nil
	}
	return os.Remove(fileName)
}
