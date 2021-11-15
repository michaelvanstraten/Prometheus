package token

import (
	"encoding/json"
	"encoding/base64"
	"github.com/michaelvanstraten/prometheus/crypto"
)

type Token struct {
	Signature 		[]byte
	Nonce 			[]byte
	ByteClaims 		[]byte
}

func (T *Token) Valid(Publickey []byte) bool {
	if T != nil {
		return crypto.Validate(T.ByteClaims, T.Nonce, T.Signature, Publickey)
	} else {
		return false
	}
} 

func (T *Token) Raw() string {
	if tokenbytes, err := json.Marshal(T); err == nil {
		return base64.StdEncoding.EncodeToString(tokenbytes)
	} else {
		println(err.Error())
		return ""
	}
}

func (T *Token) Sign(Privatekey []byte) *Token {
	var err error
	if T.Signature, T.Nonce, err = crypto.SignMessage(T.ByteClaims, 8, Privatekey); err == nil {
		return T
	} else {
		println(err.Error())
		return &Token{}
	}
}

func Parse(RawToken string, Claims interface{}) *Token {
	var token = Token{}
	if DecodedRawToken, err := base64.StdEncoding.DecodeString(RawToken); err == nil {
		if err := json.Unmarshal(DecodedRawToken, &token); err == nil {
			if err = json.Unmarshal(token.ByteClaims, &Claims); err == nil {
				return &token
			} else {
				println(err.Error())
				return nil
			}
		} else {
			println(err.Error())
			return nil
		}
	} else {
		println(err.Error())
		return nil
	}
}

func New(TokenClaims interface{}) *Token {
	var token = Token{}
	var err error
	if token.ByteClaims, err = json.Marshal(TokenClaims); err != nil {
		println(err.Error())
		return &token
	}
	return &token
}