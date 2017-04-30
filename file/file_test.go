package file

import (
	"testing"

	"github.com/bmizerany/assert"
)

func TestFileExist(t *testing.T) {

}

var data = `{
	"name":"zhangsan",
	"age":12
	}`

type student struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func TestFileApis(t *testing.T) {
	fileName := ".tmp"
	if ExistFile(fileName) {
		t.Error("exist file: " + fileName)
	}

	err := CreateFileAndWriteString(fileName, data)
	if err != nil {
		t.Errorf("create file and write file get err:%s", err.Error())
	}

	readData := ReadFileToString(fileName)
	if readData != data {
		t.Errorf("err: write data %s, read %s", data, readData)
	}

	s := &student{}
	err = UnmarshalFromFile(fileName, s)
	if err != nil {
		t.Errorf("unmarshal err:%s", err.Error())
	}

	assert.Equalf(t, 12, s.Age, "unmarshal student age:%d", s.Age)

	err = RemoveFile(fileName)
	if err != nil {
		t.Errorf("remove file ger err:%s", err.Error())
	}
}
