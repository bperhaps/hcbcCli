package network

import (
	"github.com/golang/protobuf/proto"
	"hcbcCli/hb/strc"
	"log"
)

func SendTx(tx *strc.Transaction, port string) {

	request, err := proto.Marshal(tx)
	if err != nil {
		log.Panic(err)
	}

	conn, err := NewBroadcaster(port)
	if err != nil {
		log.Panic(err)
	}
	conn.Write(request)
	conn.Close()
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
