package hb

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"hcbcCli/hb/network"
	"hcbcCli/hb/strc"
	"hcbcCli/hb/utils"
	"log"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/golang/protobuf/proto"
)

type Hb struct {
	DeviceId   []byte
	Port       string
	PuttedData map[string]interface{}
	Mutex      *sync.Mutex
	debug      bool
	conn	   *net.UDPConn
}

func NewHb() *Hb {
	return &Hb{
		DeviceId:   nil,
		Port:       "0",
		PuttedData: make(map[string]interface{}),
		Mutex:      &sync.Mutex{},
		debug:      false,
		conn:	    nil,
	}
}

func (hb *Hb) SetDebug(flag bool) {
	hb.debug = flag
}

func (hb *Hb) SetDeviceId(deviceId []byte) {
	hb.DeviceId = deviceId
}

func (hb *Hb) Close() {
	hb.conn.Close()
	GetAuth().Save()
}

func (hb *Hb) SetPort(port string) {
	if hb.conn != nil { hb.conn.Close() }
	hb.Port = port
	conn, err := network.NewBroadcaster(port); if err != nil { panic(err) }
	hb.conn = conn
}

func (hb *Hb) SendData() {
	if len(hb.PuttedData) == 0 {
		fmt.Println("inputted data is not exsist")
	}

	jsonData, err := json.Marshal(hb.PuttedData)
	if err != nil {
		log.Panic(err)
	}
	
	st := time.Now()
	authhash := GetAuth().GetAuthHash(hb.DeviceId, hb.sendUpdateAuthhash)
	fmt.Println("ah gen", time.Since(st))
	transaction := NewTransaction(jsonData, network.DataTx, authhash)

	hb.PuttedData = make(map[string]interface{})

	//st := time.Now()
	network.SendTx(transaction, hb.conn, hb.Port)
	//fmt.Println("sendTx Time : ", time.Since(st))
}

func (hb *Hb) Regist(deviceId []byte, authhash [][]byte, idx int32, perm int32) {
	regist := &strc.Regist{
		DeviceId: deviceId,
		Authhash: authhash,
		Perm:     perm,
	}
	d, err := proto.Marshal(regist)
	if err != nil {
		log.Panic(err)
	}
	tx := NewTransaction(d, network.RegistTx, GetAuth().GetAuthHash(hb.DeviceId, hb.sendUpdateAuthhash))
	network.SendTx(tx, hb.conn, hb.Port)
}

func (hb *Hb) PutData(key, value string) {
	hb.PuttedData[key] = value
}

func wattingData(result chan map[string]interface{}) {
	udpAddr, err := net.ResolveUDPAddr("udp", ":2831")
	if err != nil {
		log.Panic(err)
	}
	ln, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		log.Panic(err)
	}
	defer ln.Close()

	dataLog := make(map[string][]byte)

	for {
		buffer := make([]byte, 8192)
		n, addr, err := ln.ReadFromUDP(buffer)
		if err != nil {
			log.Panic(err)
		}
		if n == 0 {
			continue
		}
		i := strings.Index(addr.String(), ":")
		ip := addr.String()[:i]

		data := make(map[string]interface{})
		json.Unmarshal(buffer[:n], &data)

		var cnt int
		if _, ok := dataLog[ip]; !ok {
			for _, data := range dataLog {
				if bytes.Compare(data, buffer[:n]) == 0 {
					cnt++
				}
			}
			dataLog[ip] = buffer[:n]
		}

		if int(data["quorum"].(float64)) >= 2 && cnt == int(data["quorum"].(float64)) {
			delete(data, "quorum")
			result <- data
			break
		}
	}
}

func (hb *Hb) GetData(key string) (map[string]interface{}, error) {
	return hb.GetDataFromOthers(string(hb.DeviceId), key)
}

func (hb *Hb) GetDataFromOthers(deviceId, key string) (map[string]interface{}, error) {

	result := make(chan map[string]interface{})

	data := make(map[string]interface{})
	data["deviceId"] = deviceId
	data["key"] = key

	mdata, _ := json.Marshal(data)
	auth := GetAuth().GetAuthHash(hb.DeviceId, hb.sendUpdateAuthhash)
	tx := NewTransaction(mdata, network.RequestData, auth)

	go wattingData(result)
	network.SendTx(tx, hb.conn, hb.Port)

	r := <-result

	return r, nil
}

func (hb *Hb) GetRootHash() []byte {
	toggle := GetAuth().Toggle
	return GetAuth().GetRootHash(toggle)
}

func (hb *Hb) PrintEnv() {
	deviceid := hb.DeviceId
	fmt.Println("device id : ", string(deviceid))

	fmt.Println("authhash : ")

	authhashs := []interface{}{}
	for i := 0; i < 2; i++ {
		authhash := GetAuth().GetRootHash(i)
		base64Authhahs := base64.StdEncoding.EncodeToString(authhash)
		authhashs = append(authhashs, base64Authhahs)
	}

	doc, _ := json.MarshalIndent(authhashs, "", "    ")
	fmt.Println(string(doc))
}

func (hb *Hb) sendUpdateAuthhash(authhash []byte, idx, toggle int) {
	hb.Mutex.Lock()

	data := make(map[string]interface{})
	data["authhash"] = GetAuth().GetRootHash(toggle)
	data["toggle"] = toggle

	byteData, err := json.Marshal(data)
	if err != nil {
		log.Panic(err)
	}

	deviceId := hb.DeviceId
	auth := &strc.AuthHash{
		DevicdId: deviceId,
		Toggle:   int32(toggle),
		AuthHash: authhash,
		Idx:      int32(idx),
	}

	transaction := &strc.Transaction{
		TxId: nil,
		Data: byteData,
		Type: network.UpdateAuthHashTx,
		Auth: auth,
	}

	transactionId := sha256.Sum256(utils.GobEncode(transaction))
	transaction.TxId = transactionId[:]

	go network.SendTx(transaction, hb.conn, hb.Port)

	hb.Mutex.Unlock()
}
