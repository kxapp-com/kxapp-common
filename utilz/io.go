package utilz

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"errors"
	"github.com/kxapp-com/kxapp-common/cryptoz"
	"howett.net/plist"
	"io"
	"net"
	"os"
	"path/filepath"
)

func WriteToJsonFile(path string, obj any) error {
	return WriteToJsonFileSec(path, obj, "")
}
func WriteToJsonFileSec(path string, obj any, password string) error {
	data, e2 := json.Marshal(obj)
	if e2 != nil {
		return e2
	}
	if !PathExists(filepath.Dir(path)) {
		os.MkdirAll(filepath.Dir(path), 0666)
	}
	if password != "" {
		data = cryptoz.RC4Crypto(data, password)
	}
	return os.WriteFile(path, data, 0666)
}
func ReadFromJsonFile[T any](path string) (*T, error) {
	return ReadFromJsonFileSec[T](path, "")
}
func ReadFromJsonFileSec[T any](path string, password string) (*T, error) {
	data, e := os.ReadFile(path)
	if e != nil {
		return nil, e
	}
	if password != "" {
		data = cryptoz.RC4Crypto(data, password)
	}
	var inn T
	e2 := json.Unmarshal(data, &inn)
	return &inn, e2
}
func ParseJsonAs[T any](data []byte) (*T, error) {
	if data != nil {
		obj := new(T)
		e2 := json.Unmarshal(data, obj)
		return obj, e2
	}
	return nil, errors.New("input data error")
}
func ParsePlistAs[T any](data []byte) (*T, error) {
	if data != nil {
		obj := new(T)
		_, e2 := plist.Unmarshal(data, obj)
		return obj, e2
	}
	return nil, errors.New("input data error")
}

// PathExists is used to determine whether the path folder exists
// True if it exists, false otherwise
func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func EncodeGob(obj any) []byte {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	encoder.Encode(obj)
	return buffer.Bytes()
}
func DecodeGob(data []byte, ptr any) {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	decoder.Decode(ptr)
}

func ReadFileData(path string, start, length int) ([]byte, error) {
	fd, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer fd.Close()
	data := make([]byte, length)
	n, e := fd.ReadAt(data, int64(start))
	if e == io.EOF {
		return data[0:n], nil
	}
	if n == length {
		return data, nil
	} else {
		return nil, e
	}
}

func GetOneMacAddress() string {
	is, e := net.Interfaces()
	if e != nil {
		return ""
	}
	for _, i := range is {
		mac := i.HardwareAddr.String()
		if mac != "" {
			return mac
		}
	}
	return ""
}

func ParseMapToObject[T any](mp map[string]any) (*T, error) {
	bt, e1 := json.Marshal(mp)
	if e1 != nil {
		return nil, e1
	}
	return ParseJsonAs[T](bt)
}
