package pkcs12

import (
	"encoding/asn1"
	"encoding/base64"
	"errors"
	"fmt"
	"gui/crypto/sm2"
	"gui/crypto/sm3"
	"gui/crypto/sm4"
	"gui/crypto/x509"
	"log"
	"math/big"
)

type smPdu struct {
	Version     int
	PrivContent privateKeyContent
	PubContent  publicKeyContent
}

type privateKeyContent struct {
	OID1    asn1.ObjectIdentifier // {1.2.156.10197.6.1.4.2.1} SM2_Data
	OID2    asn1.ObjectIdentifier // {1.2.156.10197.1.104} SM4_CBC
	Content asn1.RawValue
}

type publicKeyContent struct {
	OID     asn1.ObjectIdentifier // {1.2.156.10197.6.1.4.2.1} SM2_Data
	Content asn1.RawValue
}

func Decode(smData []byte, password string) (privateKey *sm2.PrivateKey, certificate *x509.Certificate, err error) {
	sm := new(smPdu)
	trailing, err := asn1.Unmarshal(smData, sm)
	if err != nil {
		return nil, nil, err
	}
	if len(trailing) != 0 {
		return nil, nil, errors.New("go-pkcs12: trailing data found")
	}

	dBytes := DecryptSm2Key(password, sm.PrivContent.Content.Bytes)

	cer, err := x509.ParseCertificate(sm.PubContent.Content.Bytes)
	if err != nil {
		log.Fatal(err)
		return nil, nil, err
	}

	//pub := cer.PublicKey.(*ecdsa.PublicKey)
	pub := cer.PublicKey.(*sm2.PublicKey)
	priv := &sm2.PrivateKey{
		PublicKey: *pub,
		D:         new(big.Int).SetBytes(dBytes),
	}

	fmt.Println(cer.Issuer)

	return priv, cer, nil
}

/**
Key Derivation function (密钥导出函数)
	将密钥扩展到所需长度的密钥
**/
func KDF(z []byte) []byte {
	ct := []byte{0, 0, 0, 1}
	sm3 := sm3.New()
	sm3.Write(z)
	sm3.Write(ct)
	h := sm3.Sum(nil)
	return h
}

/*
	解密sm2私钥
*/
func DecryptSm2Key(password string, encryptedData []byte) []byte {
	if len(encryptedData) >= 32 && len(encryptedData) <= 64 {
		encoding := make([]byte, len(encryptedData), len(encryptedData))
		if len(encryptedData) != 32 && len(encryptedData) != 48 {
			base64.StdEncoding.Decode(encoding, encryptedData)
		}
		encoding = encryptedData

		h := KDF([]byte(password))

		iv := h[:16]
		key := h[16:]

		sm4 := sm4.Init(iv, key)
		out, err := sm4.Sm4Cbc(encoding, false)
		if err != nil {
			log.Fatal(err)
		}

		d := new(big.Int)
		d.SetBytes(out)
		fmt.Println(d)
		return out
	}
	return nil
}
