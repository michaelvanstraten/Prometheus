package crypto

import (
	"crypto/sha256"
	"math/rand"
	"time"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
)

var Curve secp256k1.BitCurve = secp256k1.BitCurve{
	P 		: secp256k1.S256().P,
	N 		: secp256k1.S256().N,
	B 		: secp256k1.S256().B,
	Gx 		: secp256k1.S256().Gx,
	Gy 		: secp256k1.S256().Gy,
	BitSize : secp256k1.S256().BitSize,
}

func GenerateRandomBytes(NumberOfRandomBytes int) []byte {
	var randombytes 		= make([]byte, NumberOfRandomBytes)
	rand.Seed(time.Now().UnixNano()*420*69*69)
	rand.Read(randombytes)
	return randombytes
}

func GeneratePublicKey(Privatekey []byte) []byte {
	var x, y  = Curve.ScalarBaseMult(Privatekey)
	return Curve.Marshal(x, y)
}

func Validate(Message []byte, RandomBytes []byte, Signature []byte, Publickey []byte) bool {
	var tovalidatemessage 	= make([]byte, 0, len(Message)+len(RandomBytes)+2)
	tovalidatemessage 		= append(tovalidatemessage, Message...)
	tovalidatemessage 		= append(tovalidatemessage, RandomBytes...)
	var messagehash 		= sha256.Sum256(tovalidatemessage)
	return secp256k1.VerifySignature(Publickey, messagehash[:], Signature[:len(Signature)-1])
}

func SignMessage(Message []byte, NumberOfRandomBytes int, Privatekey []byte) ([]byte, []byte, error) {
	var randombytes 	= GenerateRandomBytes(int(NumberOfRandomBytes))
	var tosignmessage 	= make([]byte, 0, len(Message)+NumberOfRandomBytes+1)
	tosignmessage 		= append(tosignmessage, Message...)
	tosignmessage 		= append(tosignmessage, randombytes...)
	var messagehash 	= sha256.Sum256(tosignmessage)
	if signature, err 	:= secp256k1.Sign(messagehash[:], Privatekey); err == nil {
		return signature, randombytes, nil
	} else {			
		return nil, nil, err
	}
}