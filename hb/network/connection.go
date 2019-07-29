package network

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/golang/protobuf/proto"
	"hcbcCli/hb/strc"
	"log"
	"net"
	"sync"
)

type connection struct {
	net.Conn
	byteBuffer bytes.Buffer
	once       sync.Once
	byteChan   chan []byte
}

func SendTx(tx *strc.Transaction, conn *net.UDPConn, port string) {

	request, err := proto.Marshal(tx)
	if err != nil {
		log.Panic(err)
	}

	conn.Write(request)
}

func NewConn(conn net.Conn) *connection {
	return &connection{
		Conn:     conn,
		byteChan: make(chan []byte),
	}
}

func (conn *connection) ReadData() ([]byte, error) {

	conn.once.Do(func() {
		go func() {
			defer func() {
				conn.Close()
				close(conn.byteChan)
			}()

			for {
				data := make([]byte, 4096)
				n, err := conn.Read(data)
				if err != nil {
					break
				}

				conn.byteChan <- data[:n]
			}
		}()
	})

	var data []byte
	if conn.byteBuffer.Len() <= 0 {
		var ok bool
		data, ok = <-conn.byteChan
		if !ok {
			return nil, fmt.Errorf("eof")
		}

		conn.byteBuffer.Write(data)
	}

	lengthByte := conn.byteBuffer.Next(4)
	msgLength := binary.LittleEndian.Uint32(lengthByte)

	if int(msgLength) > conn.byteBuffer.Len() {
		for int(msgLength) > conn.byteBuffer.Len() {
			data, ok := <-conn.byteChan
			if !ok {
				break
			}
			conn.byteBuffer.Write(data)
		}
	}

	return conn.byteBuffer.Next(int(msgLength)), nil
}

//
//func sendData(addr string, request []byte) ([]byte, error) {
//	conn, err := net.Dial("tcp", addr)
//	defer func() {
//		if conn != nil {
//			conn.Close()
//		}
//	}()
//	if err != nil {
//		return nil, err
//	}
//
//	lengthBuf := make([]byte, 4)
//	binary.LittleEndian.PutUint32(lengthBuf, uint32(len(request)))
//
//	appendedRequest := append(lengthBuf, request...)
//
//	_, err = conn.Write(appendedRequest)
//	if err != nil {
//		return nil, err
//	}
//
//	response, err := readData(conn)
//	if err != nil {
//		log.Panic(err)
//	}
//	return response, nil
//
//}
//func readData(conn net.Conn) ([]byte, error) {
//	lengthBuf := make([]byte, 4)
//	_, err := conn.Read(lengthBuf)
//	if err != nil {
//		return nil, errors.New("dataread error")
//	}
//	msgLength := binary.LittleEndian.Uint32(lengthBuf)
//
//	buf := make([]byte, 4096)
//	var request bytes.Buffer
//
//	for 0 < msgLength {
//		n, err := conn.Read(buf)
//		if err != nil {
//			return nil, errors.New("dataread error")
//		}
//		if 0 < n {
//			data := buf[:n]
//			request.Write(data)
//			msgLength -= uint32(n)
//		}
//	}
//
//	return request.Bytes(), nil
//}
