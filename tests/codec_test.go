package tests

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	core2 "jim/client/core"
	"jim/tcp/core"
	"testing"
)

func TestCodec(t *testing.T) {
	//var value uint16 = 50000
	//b1 := byte(value)
	//b2 := byte(value >> 8)
	//buf.WriteByte(b2)
	//buf.WriteByte(b1)
	//fmt.Println(buf.Bytes())
	//v:=uint16(b2)<<8+uint16(b1)
	//fmt.Println(v)
	//-----------------------------------------------------
	codec := core.NewJimDataFrameCodec()
	s1 := "hello world"
	pack1, err := codec.Encode(nil, []byte(s1))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("encode result:", pack1)
	s2 := "You are welcome"
	pack2, err := codec.Encode(nil, []byte(s2))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("encode result:", pack2)
	result := make([]byte, 0)
	buf := bytes.NewBuffer(result)
	buf.Write(pack1)
	buf.Write(pack2)
	//------------------------------------------------------
	bs := buf.Bytes()
	headerLen := 2
	for {
		if len(bs) == 0 {
			break
		}
		header := bs[:2]
		byteBuffer := bytes.NewBuffer(header)
		var length uint16
		_ = binary.Read(byteBuffer, binary.BigEndian, &length)
		dataLen := int(length) // max int32 can contain 210MB payload
		protocolLen := headerLen + dataLen
		bb := bs[2:protocolLen]
		fmt.Println("decode result:", string(bb))
		bs = bs[protocolLen:]
	}
}

func TestClientCodec(t *testing.T) {
	s := "hello world"
	ens, err := core2.Encode([]byte(s))
	if err != nil {
		printl(err.Error())
		return
	}
	printl(ens)

	ioreader := bufio.NewReader(bytes.NewReader(ens))
	des, err := core2.Decode(ioreader)
	if err != nil {
		printl(err.Error())
		return
	}
	printl(string(des))
}
