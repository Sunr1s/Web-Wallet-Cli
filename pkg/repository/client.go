package repository

import (
	"crypto/rsa"
	"crypto/x509"
	"database/sql"
	"encoding/json"
	"encoding/pem"
	"io/ioutil"
	"os"

	bc "github.com/Sunr1s/webclient/pkg/blockchain"
)

type Wallet struct {
	Public  string
	Private string
}

var (
	Address []string
	User    *bc.User
)

const (
	ADD_BLOCK = iota + 1
	ADD_TRNSX
	GET_BLOCK
	GET_LHASH
	GET_BLNCE

	SEPARATOR = "_SEPARATOR_"
)

type ClientAuth struct {
	db *sql.DB
}

func NewClientAuth(db *sql.DB) *ClientAuth {
	return &ClientAuth{db: db}
}

func encode(privateKey *rsa.PrivateKey, publicKey *rsa.PublicKey) (string, string) {
	x509Encoded := x509.MarshalPKCS1PrivateKey(privateKey)
	pemEncoded := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: x509Encoded})

	x509EncodedPub := x509.MarshalPKCS1PublicKey(publicKey)
	pemEncodedPub := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: x509EncodedPub})

	return string(pemEncoded), string(pemEncodedPub)
}

func decode(pemEncoded string, pemEncodedPub string) (*rsa.PrivateKey, *rsa.PublicKey) {
	block, _ := pem.Decode([]byte(pemEncoded))
	x509Encoded := block.Bytes
	privateKey, _ := x509.ParsePKCS1PrivateKey(x509Encoded)

	blockPub, _ := pem.Decode([]byte(pemEncodedPub))
	x509EncodedPub := blockPub.Bytes
	genericPublicKey, _ := x509.ParsePKCS1PublicKey(x509EncodedPub)
	publicKey := genericPublicKey

	return privateKey, publicKey
}

func NewClient(filename string) *bc.User {
	user := bc.NewUser()
	if user == nil {
		return nil
	}
	err := writeFile(filename, user.Private(), user.Public())
	if err != nil {
		return nil
	}
	return user
}

func writeFile(foldername string, priv *rsa.PrivateKey, pub *rsa.PublicKey) error {

	encPriv, encPub := encode(priv, pub)
	kdata := Wallet{
		Public:  encPub,
		Private: encPriv,
	}

	file, _ := json.MarshalIndent(kdata, "", " ")
	err := os.Mkdir(foldername, 0750)
	if err != nil && !os.IsExist(err) {
		return err
	}
	return ioutil.WriteFile(foldername+"/wallet.dat", file, 0644)
}

func (r *AuthSqlite) ReadFile(filename string) string {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return ""
	}
	return string(data)
}

func readKeys(data string, key bool) string {
	var keys Wallet

	err := json.Unmarshal([]byte(data), &keys)
	priv, pub := decode(string(keys.Private), string(keys.Public))
	if err != nil {
		return ""
	}
	if err != nil {
		return ""
	}
	if key {
		return bc.StringPrivate(priv)
	} else {
		return bc.StringPublic(pub)
	}

}
