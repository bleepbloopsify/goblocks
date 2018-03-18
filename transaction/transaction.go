package transaction

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"

	log "github.com/sirupsen/logrus"
)

// TxOut represents an unspent output into a specific address
type TxOut struct {
	address [32]byte
	amount  uint64
}

// TxOutByteSize Byte array size of txout
const TxOutByteSize = 40

func (txout *TxOut) toByteArray() []byte {
	ba := make([]byte, TxOutByteSize)
	copy(ba[:], txout.address[:])
	binary.LittleEndian.PutUint64(ba[32:], txout.amount)
	return ba
}

// UnspentTxOut is a structure that holds the details of unspent transactions
type UnspentTxOut struct {
	TxOutID    [32]byte
	txOutIndex uint64
	address    [32]byte
	amount     uint64
}

// TxIn represents an unlocking of resources and always refers to a txout for
// validity.
type TxIn struct {
	txOutID    [32]byte
	txOutIndex uint64
	signature  []byte
}

// TxInByteSize Byte array size of Txin
const TxInByteSize = 40

func (txin *TxIn) toByteArray() []byte {
	ba := make([]byte, TxInByteSize)
	copy(ba[:], txin.txOutID[:])
	binary.LittleEndian.PutUint64(ba[32:], txin.txOutIndex)
	return ba
}

// Transaction is a struct holding the details of a transaction together
type Transaction struct {
	ID     [32]byte // SHA256 hash
	txIns  []*TxIn
	txOuts []*TxOut
}

// CalculateTransactionID Generates Transaction ID hash from all of its txins and txouts.
func CalculateTransactionID(transaction *Transaction) [32]byte {
	ba := make([]byte, 40*len(transaction.txIns)+40*len(transaction.txOuts))
	bai := 0
	for _, txin := range transaction.txIns {
		copy(ba[bai:], txin.toByteArray()[:])
		bai += TxInByteSize
	}
	for _, txout := range transaction.txOuts {
		copy(ba[bai:], txout.toByteArray()[:])
		bai += TxOutByteSize
	}

	return sha256.Sum256(ba)
}

// Signature generates the signature of a transaction
func (transaction *Transaction) Signature(key *ecdsa.PrivateKey) []byte {
	id := CalculateTransactionID(transaction)
	signature, err := key.Sign(rand.Reader, id[:], crypto.SHA256.HashFunc())
	if err != nil {
		log.Error(err)
		panic(err)
	}

	return signature
}

// IsValidTxIn Checks if a txin is valid
func IsValidTxIn(txin *TxIn) bool {

	return true
}

// IsValidTransaction Checks if a transaction is valid
func IsValidTransaction(transaction *Transaction) bool {
	if CalculateTransactionID(transaction) != transaction.ID {
		log.Warn("Transaction ID does not match generated")
		return false
	}

	return true
}
