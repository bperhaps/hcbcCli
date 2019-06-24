package hb

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"hcbcCli/hb/strc"
	"hcbcCli/hb/utils"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/boltdb/bolt"
)

type Authenticate struct {
	AuthHash   [2][][]byte
	Toggle     int
	Idx        int
	MaxIdx     int
}

const authDB = "db/authenticate.db"
const authBucket = "authenticate"
const authhashBucket = "authhashs"

var mutex = &sync.Mutex{}

func GetAuth() *Authenticate {
	if utils.IsFileExist(authDB) { return newAuth(1000) }

	var deviceByte []byte

	db, err := bolt.Open(authDB, 0600, nil)
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()

	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(authBucket))
		deviceByte = b.Get([]byte(authBucket))
		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	authenticate := DeSerialize(deviceByte)

	return authenticate
}

func newAuth(maxIdx int) *Authenticate {

	authhash := &Authenticate{
		[2][][]byte{},
		0,
		maxIdx - 1,
		maxIdx,
	}

	authhash.SetAuthhash(0)
	authhash.SetAuthhash(1)

	authhash.Save()

	return authhash

}

func (authenticate *Authenticate) GetRootHash(toggle int) []byte {
	return authenticate.AuthHash[toggle][authenticate.MaxIdx]
}

func (authenticate *Authenticate) GetAuthHash(deviceId []byte, sendUpdateAuthhash func(authhash []byte, idx, toggle int)) *strc.AuthHash {
	mutex.Lock()

	authhash := authenticate.AuthHash[authenticate.Toggle][authenticate.Idx]
	idx := authenticate.MaxIdx - authenticate.Idx
	toggle := authenticate.Toggle
	authenticate.Idx--

	auth := &strc.AuthHash{
		DevicdId: deviceId,
		Toggle:   int32(toggle),
		AuthHash: authhash,
		Idx:      int32(idx),
	}

	if authenticate.Idx <= 1 {
		authhash := authenticate.AuthHash[authenticate.Toggle][authenticate.Idx]
		idx := authenticate.MaxIdx - authenticate.Idx

		authenticate.SetAuthhash(authenticate.Toggle)
		go sendUpdateAuthhash(authhash, idx, authenticate.Toggle)
		authenticate.Toggle = (authenticate.Toggle + 1) % 2
		authenticate.Idx = authenticate.MaxIdx - 1
	}
	mutex.Unlock()

	fmt.Println("toggle : ", authenticate.Toggle, " idx :", authenticate.Idx)

	authenticate.Save()

	return auth

	//일정 수준이 되면 SetAuthhash를 호출하는 부분이 필요함.
}

func (authenticate *Authenticate) SetAuthhash(toggle int) {
	seed := sha256.Sum256([]byte(strconv.FormatInt(time.Now().UnixNano(), 10)))
	time.Sleep(1)

	if len(authenticate.AuthHash[toggle]) == 0 {
		authenticate.AuthHash[toggle] = append(authenticate.AuthHash[toggle], []byte{})
	}
	authenticate.AuthHash[toggle][0] = seed[:]

	for i := 1; i <= authenticate.MaxIdx; i++ {
		authenticate.AuthHash[toggle] = append(authenticate.AuthHash[toggle], []byte{})
		prevKey := sha256.Sum256(authenticate.AuthHash[toggle][i-1])
		authenticate.AuthHash[toggle][i] = prevKey[:]
	}
}

func DeSerialize(authenticateByte []byte) *Authenticate {
	var payload Authenticate

	dec := gob.NewDecoder(bytes.NewReader(authenticateByte))
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	return &payload
}

func (auth Authenticate) Serialize() []byte {
	return utils.GobEncode(auth)
}

func (auth *Authenticate) Save() {

	db, err := bolt.Open(authDB, 0600, nil)
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		a, err := tx.CreateBucketIfNotExists([]byte(authBucket))
		if err != nil {
			log.Panic(err)
		}

		err = a.Put([]byte(authBucket), auth.Serialize())
		if err != nil {
			log.Panic(err)
		}

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

}
