package encoderz

import (
	"bytes"
	"encoding/binary"
)

// 编码函数，将多个切片进行长度前缀编码
func EncodeLengthPrefixed(slices ...[]byte) []byte {
	buffer := new(bytes.Buffer)
	for _, data := range slices {
		length := len(data)
		// 使用 binary 包中的 PutVarint 方法将长度写入 buffer
		binary.Write(buffer, binary.LittleEndian, uint32(length))
		// 将数据体写入 buffer
		buffer.Write(data)
	}
	return buffer.Bytes()
}

// 解码函数，从长度前缀编码的数据中解析出原始切片
func DecodeLengthPrefixed(data []byte) ([][]byte, error) {
	buffer := bytes.NewReader(data)
	result := make([][]byte, 0)
	for buffer.Len() > 0 {
		// 使用 binary 包中的 ReadVarint 方法读取长度信息
		var length uint32
		err := binary.Read(buffer, binary.LittleEndian, &length)
		//length, err := binary.ReadVarint(buffer)
		if err != nil {
			return nil, err
		}
		// 根据长度读取数据体
		body := make([]byte, length)
		_, err = buffer.Read(body)
		if err != nil {
			return nil, err
		}
		result = append(result, body)
	}
	return result, nil
}
