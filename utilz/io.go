package utilz

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/kxapp-com/kxapp-common/cryptoz"
	_ "github.com/kxapp-com/kxapp-common/cryptoz"
	"github.com/kxapp-com/kxapp-common/filez"
	"image"
	"image/png"
	"io/fs"
	"strings"
	//"github.com/kxapp-com/kxapp-common/cryptoz"
	"howett.net/plist"
	"io"
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
	if !filez.PathExists(filepath.Dir(path)) {
		os.MkdirAll(filepath.Dir(path), fs.ModePerm)
	}
	if password != "" {
		data = cryptoz.RC4Crypto(data, password)
	}
	return os.WriteFile(path, data, fs.ModePerm)
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

//
//// PathExists is used to determine whether the path folder exists
//// True if it exists, false otherwise
//func PathExists(path string) bool {
//	_, err := os.Stat(path)
//	if err == nil {
//		return true
//	}
//	if os.IsNotExist(err) {
//		return false
//	}
//	return false
//}

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

func ParseMapToObject[T any](mp map[string]any) (*T, error) {
	bt, e1 := json.Marshal(mp)
	if e1 != nil {
		return nil, e1
	}
	return ParseJsonAs[T](bt)
}
func StructToMap(obj any) (map[string]any, error) {
	objMap := make(map[string]any)
	objJson, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(objJson, &objMap)
	if err != nil {
		return nil, err
	}
	return objMap, nil
}

/*
*
读取配置文件.properties内容，配置文件可以用#!开头表示注释，用=表示键值对
增强的可以用\=表示key,value里面有=号字符
*/
func ReadPropertiesFile(filename string) (map[string]any, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return PropertiesDecode(file)
}

/*
*
读取配置文件.properties内容，配置文件可以用#!开头表示注释，用=表示键值对
增强的可以用\=表示key,value里面有=号字符
*/
func PropertiesDecode(readerPropertes io.Reader) (map[string]any, error) {
	const tp = "12e428c4-9902-42f9-8dad-d5c8dbeae091"
	var result = map[string]any{}
	reader := bufio.NewReader(readerPropertes)
	for {
		line, e := reader.ReadString('\n')
		if line != "" {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "#") || strings.HasPrefix(line, "!") {
				continue
			}
			line = strings.ReplaceAll(line, "\\=", tp) //利用uuid基本不会重复的特点可以随便替换字符串再替换回去
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				parts[0] = strings.ReplaceAll(parts[0], tp, "=")
				parts[1] = strings.ReplaceAll(parts[1], tp, "=")
				result[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
			}
		}
		if e != nil {
			break
		}
	}
	return result, nil
}

func PropertiesEncode(properties map[string]any) []byte {
	var content strings.Builder
	for key, value := range properties {
		content.WriteString(key)
		content.WriteString("=")
		content.WriteString(fmt.Sprintf("%v", value))
		content.WriteString("\n")
	}
	return []byte(content.String())
}

func ResizeImage(inputPath string, outputPath string, width int, height int) error {
	// Open the input file
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer inputFile.Close()

	// Decode the input image
	inputImage, _, err := image.Decode(inputFile)
	if err != nil {
		return err
	}
	outputImage := imaging.Resize(inputImage, width, height, imaging.Lanczos)
	if _, err := os.Stat(filepath.Dir(outputPath)); os.IsNotExist(err) {
		os.MkdirAll(filepath.Dir(outputPath), 0644)
	}
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outputFile.Close()
	err = png.Encode(outputFile, outputImage)
	if err != nil {
		return err
	}

	return nil
}
