package wallet

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"log"
	"math/big"

	"github.com/lshtar13/blockchain/base58"
	"golang.org/x/crypto/ripemd160"
)

type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}

type WalletMimic struct {
	PrivateKey struct {
		X *big.Int
		Y *big.Int
		D *big.Int
	}

	PublicKey []byte
}

const AddressChecksumLen = 4
const version = byte(0x00)
const walletFile = "wallet_%s.dat"

func NewWallet() *Wallet {
	private, public := newKeyPair()
	wallet := Wallet{private, public}
	return &wallet
}

func newKeyPair() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()
	private, _ := ecdsa.GenerateKey(curve, rand.Reader)
	pubKey := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)

	return *private, pubKey
}

func (w *Wallet) GetAddress() string {
	pubKeyHash := HashPubkey(w.PublicKey)

	versionedPayload := append([]byte{version}, pubKeyHash...)
	checksum := calcChecksum(versionedPayload)

	fullPayload := append(versionedPayload, checksum...)
	address := base58.Base58Encode(fullPayload)

	return string(address)
}

func ValidateAddress(addr string) bool {
	decoded := base58.Base58Decode([]byte(addr))
	tgtCheckSum := getCheckSum(decoded)
	version := getVersion(decoded)
	pubKeyHash := getPubKeyHash(decoded)
	correctCheckSum := calcChecksum(append([]byte{version}, pubKeyHash...))

	return bytes.Equal(tgtCheckSum, correctCheckSum)
}

func (w *Wallet) MarshalWallet() []byte {
	var buf bytes.Buffer
	mimic := WalletMimic{}
	mimic.PrivateKey.X = w.PrivateKey.X
	mimic.PrivateKey.Y = w.PrivateKey.Y
	mimic.PrivateKey.D = w.PrivateKey.D
	mimic.PublicKey = w.PublicKey

	encoder := gob.NewEncoder(&buf)
	err := encoder.Encode(mimic)
	if err != nil {
		log.Panic(err)
	}

	return buf.Bytes()
}

func UnmarshalWallet(unmarshaled []byte) *Wallet {
	wallet := Wallet{}
	mimic := WalletMimic{}

	decoder := gob.NewDecoder(bytes.NewReader(unmarshaled))
	err := decoder.Decode(&mimic)
	if err != nil {
		log.Panic(err)
	}

	wallet.PublicKey = mimic.PublicKey
	wallet.PrivateKey.Curve = elliptic.P256()
	wallet.PrivateKey.D = mimic.PrivateKey.D
	wallet.PrivateKey.Y = mimic.PrivateKey.Y
	wallet.PrivateKey.X = mimic.PrivateKey.X

	return &wallet
}

func HashPubkey(pubKey []byte) []byte {
	publicSHA256 := sha256.Sum256(pubKey)

	RIPEMD160Hahser := ripemd160.New()
	_, err := RIPEMD160Hahser.Write(publicSHA256[:])
	if err != nil {
		log.Panic("Error : hash pubkey")
	}
	publicRIPEMD160 := RIPEMD160Hahser.Sum(nil)

	return publicRIPEMD160
}

func calcChecksum(payload []byte) []byte {
	firstSHA := sha256.Sum256(payload)
	secondSHA := sha256.Sum256(firstSHA[:])

	return secondSHA[:AddressChecksumLen]
}

func getCheckSum(payload []byte) []byte {
	return payload[len(payload)-AddressChecksumLen:]
}

func getVersion(payload []byte) byte {
	return payload[0]
}

func getPubKeyHash(payload []byte) []byte {
	return payload[1 : len(payload)-AddressChecksumLen]
}
