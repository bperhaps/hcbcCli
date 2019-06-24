package hb

import (
	"crypto/sha256"
	"hcbcCli/hb/strc"
	"hcbcCli/hb/utils"
)

func NewTransaction(data []byte, txType int32, auth *strc.AuthHash) *strc.Transaction {

	tx := &strc.Transaction{
		TxId:     nil,
		Data:     data,
		Type:     txType,
		Auth:     auth,
	}
	transactionId := sha256.Sum256(utils.GobEncode(tx))
	tx.TxId = transactionId[:]

	return tx
}
