package core

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/panjf2000/gnet"
)

type JimDataFrameCodec struct {
}

func NewJimDataFrameCodec() *JimDataFrameCodec {
	return &JimDataFrameCodec{}
}

func (cc *JimDataFrameCodec) Encode(c gnet.Conn, buf []byte) ([]byte, error) {
	// 数据长度控制在
	dataLen := uint16(len(buf))
	if dataLen > 65530 || dataLen == 0 {
		return nil, errors.New("Package data length is wrong.")
	}
	// 创建一个最终buff
	result := make([]byte, 0)
	buffer := bytes.NewBuffer(result)
	// 写入长度段数据
	if err := binary.Write(buffer, binary.BigEndian, dataLen); err != nil {
		s := fmt.Sprintf("Pack datalength error , %v", err)
		return nil, errors.New(s)
	}
	// 写入原始数据
	if err := binary.Write(buffer, binary.BigEndian, buf); err != nil {
		s := fmt.Sprintf("Pack data error , %v", err)
		return nil, errors.New(s)
	}

	return buffer.Bytes(), nil
}

func (cc *JimDataFrameCodec) Decode(c gnet.Conn) ([]byte, error) {
	// 长度段长度为2个byte
	headerLen := 2 // uint16
	// 取出长度段数据
	if size, header := c.ReadN(headerLen); size == headerLen {
		byteBuffer := bytes.NewBuffer(header)
		var length uint16
		_ = binary.Read(byteBuffer, binary.BigEndian, &length)
		// 解析原始数据长度
		dataLen := int(length) // max int32 can contain 210MB payload
		protocolLen := headerLen + dataLen
		// 根据长度获取原始数据
		if dataSize, data := c.ReadN(protocolLen); dataSize == protocolLen {
			c.ShiftN(protocolLen)
			return data[headerLen:], nil
		}
		return nil, errors.New("not enough payload data")

	}
	return nil, errors.New("not enough header data")
}
