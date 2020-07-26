package core

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

func Encode(in []byte) (out []byte, err error) {
	dataLen := uint16(len(in))
	if dataLen > 65530 || dataLen == 0 {
		err = errors.New("Package data length is wrong.")
		return
	}
	// 创建一个最终buff
	result := make([]byte, 0)
	buffer := bytes.NewBuffer(result)
	// 写入长度段数据
	if err = binary.Write(buffer, binary.BigEndian, dataLen); err != nil {
		errors.New(fmt.Sprintf("Pack datalength error , %v", err))
		return
	}
	// 写入原始数据
	if err = binary.Write(buffer, binary.BigEndian, in); err != nil {
		errors.New(fmt.Sprintf("Pack data error , %v", err))
		return
	}
	out = buffer.Bytes()
	return
}

func Decode(reader *bufio.Reader) (out []byte, err error) {
	header := make([]byte, 2)
	n, err := reader.Read(header)
	if err != nil {
		return
	}
	if n != 2 {
		err = errors.New("package header size error")
		return
	}
	byteBuffer := bytes.NewBuffer(header)
	var length uint16
	err = binary.Read(byteBuffer, binary.BigEndian, &length)
	if err != nil {
		return
	}
	dataLen := int(length)
	out = make([]byte, dataLen)
	n, err = reader.Read(out)
	if err != nil {
		return
	}
	if n != dataLen {
		err = errors.New("package body size error")
		return
	}
	return
}
