package file

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
)

func ExistFile(fileName string) bool {
	_, err := os.Stat(fileName)
	return err == nil || os.IsExist(err)
}

func CopyFile(src, dst string) (w int64, err error) {
	srcFile, err := os.Open(src)
	if err != nil {
		return -1, err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)

	if err != nil {
		return -1, err
	}
	defer dstFile.Close()

	return io.Copy(dstFile, srcFile)
}

func ReadFileToString(fileName string) string {
	file, err := os.Open(fileName)
	if err != nil {
		panic("open file get err:" + err.Error())
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		panic("read file get err:" + err.Error())
	}
	return string(data)
}

func UnmarshalFromFile(fileName string, v interface{}) error {
	file, err := os.Open(fileName)
	if err != nil {
		panic("open file get error:" + err.Error())
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		panic("read file get error:" + err.Error())
	}
	return json.Unmarshal(data, v)
}
