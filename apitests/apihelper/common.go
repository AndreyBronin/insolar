package apihelper

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
)

func NewMemberSignature() (MemberSignature, error) {
	var err error
	privateKey := new(ecdsa.PrivateKey)
	privateKey, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return MemberSignature{}, err
	}
	var publicKey ecdsa.PublicKey
	publicKey = privateKey.PublicKey
	// Convert the public key into PEM format:
	x509PublicKey, err := x509.MarshalPKIXPublicKey(&publicKey)
	if err != nil {
		return MemberSignature{}, err
	}
	pemPublicKey := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: x509PublicKey})

	return MemberSignature{
		PublicKey:     publicKey,
		PrivateKey:    privateKey,
		X509PublicKey: x509PublicKey,
		PemPublicKey:  pemPublicKey,
	}, err
}

//func loadAdminMemberKeys() (string, string) {
//	text, err := ioutil.ReadFile("~/go/src/github.com/insolar/insolar/.artifacts/launchnet/configs/migration_admin_member_keys.json")
//	if err != nil {
//		errors.Wrapf(err, "[ loadMemberKeys ] could't load member keys")
//	}
//	var data map[string]string
//	err = json.Unmarshal(text, &data)
//	if err != nil {
//		 errors.Wrapf(err, "[ loadMemberKeys ] could't unmarshal member keys")
//	}
//	if data["private_key"] == "" || data["public_key"] == "" {
//		errors.New("[ loadMemberKeys ] could't find any keys")
//	}
//	privateKey := data["private_key"]
//	publicKey := data["public_key"]
//
//	return privateKey, publicKey
//}