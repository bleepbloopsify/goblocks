package transaction

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"io/ioutil"

	log "github.com/sirupsen/logrus"
)

// GenerateKey generates key
func GenerateKey() *ecdsa.PrivateKey {
	privkey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		log.Error(err)
		panic(err)
	}
	return privkey
}

// SaveKey saves private key to a file
func SaveKey(fname string, key *ecdsa.PrivateKey) {
	marshalled, err := x509.MarshalECPrivateKey(key)
	if err != nil {
		log.Error(err)
		panic(err)
	}
	ioutil.WriteFile(fname, marshalled, 0600)
}

// LoadKey loads key from a file
func LoadKey(fname string) *ecdsa.PrivateKey {
	bytekey, err := ioutil.ReadFile(fname)
	if err != nil {
		log.Error(err)
		panic(err)
	}
	key, err := x509.ParseECPrivateKey(bytekey)
	if err != nil {
		log.Error(err)
		panic(err)
	}
	return key
}
