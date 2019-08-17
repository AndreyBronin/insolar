package apihelper

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/asn1"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"github.com/insolar/insolar/api/requester"
	"github.com/insolar/insolar/apitests/apihelper/apilogger"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"testing"
)

type errorStruct struct {
	Error struct {
		Data struct {
			RequestReference string `json:"requestReference"`
			TraceID          string `json:"traceID"`
		} `json:"data"`
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

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

func Sign(payload interface{}, privateKey *ecdsa.PrivateKey) (string, string, map[string]string) {
	var err error
	// get hash of byte slice of the payload encoded with the same way as openapi-generator does in the generated client.
	// this is done to avoid setting incorrect body value into request by generated code.
	// if you use custom code to create insolar-api client, use 'json.Marshal(payload)' and get hash value of it s result.
	bodyBuf := &bytes.Buffer{}
	err = json.NewEncoder(bodyBuf).Encode(payload)
	if err != nil {
		log.Fatalln(err)
	}
	request, err := http.NewRequest("ignore", "ignore", bodyBuf)
	memberCreateRequest := reflect.TypeOf(payload)
	rawBody, err := requester.UnmarshalRequest(request, &memberCreateRequest)
	if err != nil {
		apilogger.Fatal(err)
	}
	hash := sha256.Sum256(rawBody)

	// Sign the hash with the private key:
	r, s, err := ecdsa.Sign(rand.Reader, privateKey, hash[:])
	if err != nil {
		apilogger.Fatal(err)
	}

	// See if the signature is valid:
	valid := ecdsa.Verify(&privateKey.PublicKey, hash[:], r, s)
	if !valid {
		apilogger.Fatal("signature not verified")
	}

	// Convert the signature into ASN.1 format:
	sig := ecdsaSignature{
		R: r,
		S: s,
	}
	signature, _ := asn1.Marshal(sig)

	// Convert both hash and signature into a Base64 string:
	hash64 := base64.StdEncoding.EncodeToString(hash[:])
	signature64 := base64.StdEncoding.EncodeToString(signature)

	var Digest = "SHA-256=" + hash64
	var Signature = "keyId=\"member-pub-key\", algorithm=\"ecdsa\", headers=\"digest\", signature=" + signature64
	return Digest, Signature, map[string]string{"Digest": Digest, "Signature": Signature}
}

func LoadAdminMemberKeys() (string, string) {
	text, err := ioutil.ReadFile("~/go/src/github.com/insolar/insolar/.artifacts/launchnet/configs/migration_admin_member_keys.json")
	if err != nil {
		errors.Wrapf(err, "[ loadMemberKeys ] could't load member keys")
	}
	var data map[string]string
	err = json.Unmarshal(text, &data)
	if err != nil {
		errors.Wrapf(err, "[ loadMemberKeys ] could't unmarshal member keys")
	}
	if data["private_key"] == "" || data["public_key"] == "" {
		errors.New("[ loadMemberKeys ] could't find any keys")
	}
	privateKey := data["private_key"]
	publicKey := data["public_key"]
	apilogger.Printf("pk: %v/n privk: %v", publicKey, privateKey)

	return privateKey, publicKey
}

func CheckResponseHasNoError(t *testing.T, response interface{}) {
	j, err := json.Marshal(response)
	require.Nil(t, err)
	var errorBody errorStruct
	err = json.Unmarshal(j, &errorBody)
	require.Nil(t, err, "error while unmarshaling")
	if errorBody.Error.Message != "" || errorBody.Error.Code != 0 {
		require.Emptyf(t, errorBody.Error.Message, "error in response: %v", errorBody.Error.Message)
	}
}
