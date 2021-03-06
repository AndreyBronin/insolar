//
// Copyright 2019 Insolar Technologies GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package secrets

import (
	"bytes"
	"crypto"
	"encoding/json"
	"io"
	"io/ioutil"
	"path/filepath"

	"github.com/insolar/insolar/platformpolicy"
	"github.com/pkg/errors"
)

// KeyPair holds private/public keys pair.
type KeyPair struct {
	Private crypto.PrivateKey
	Public  crypto.PublicKey
}

// GenerateKeyPair generates private/public keys pair. It uses default platform policy.
func GenerateKeyPair() (*KeyPair, error) {
	ks := platformpolicy.NewKeyProcessor()
	privKey, err := ks.GeneratePrivateKey()
	if err != nil {
		return nil, errors.Wrap(err, "couldn't generate private key")
	}
	return &KeyPair{
		Private: privKey,
		Public:  ks.ExtractPublicKey(privKey),
	}, nil
}

// GetPublicKeyFromFile reads private/public keys pair from json file and return public key
func GetPublicKeyFromFile(file string) (string, error) {
	pair, err := ReadKeysFile(file)
	if err != nil {
		return "", errors.Wrap(err, "couldn't get keys")
	}
	return platformpolicy.MustPublicKeyToString(pair.Public), nil
}

// ReadKeysFile reads private/public keys pair from json file.
func ReadKeysFile(file string) (*KeyPair, error) {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, errors.Wrapf(err, " couldn't read keys file %v", file)
	}
	return ReadKeys(bytes.NewReader(b))
}

// ReadKeysFile reads and parses json from reader, returns parsed private/public keys pair.
func ReadKeys(r io.Reader) (*KeyPair, error) {
	var keys map[string]string
	err := json.NewDecoder(r).Decode(&keys)
	if err != nil {
		return nil, errors.Wrapf(err, "fail unmarshal keys data")
	}
	if keys["private_key"] == "" {
		return nil, errors.New("empty private key")
	}
	if keys["public_key"] == "" {
		return nil, errors.New("empty public key")
	}

	kp := platformpolicy.NewKeyProcessor()

	privateKey, err := kp.ImportPrivateKeyPEM([]byte(keys["private_key"]))
	if err != nil {
		return nil, errors.Wrapf(err, "fail import private key")
	}
	publicKey, err := kp.ImportPublicKeyPEM([]byte(keys["public_key"]))
	if err != nil {
		return nil, errors.Wrapf(err, "fail import private key")
	}

	return &KeyPair{
		Private: privateKey,
		Public:  publicKey,
	}, nil

}

// ReadKeysFromDir reads directory, tries to parse every file in it as json with private/public keys pair
// returns list of parsed private/public keys pairs.
func ReadKeysFromDir(dir string) ([]*KeyPair, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, errors.Wrapf(err, "can't read dir %v", dir)
	}

	nodes := make([]*KeyPair, 0, len(files))
	for _, f := range files {
		pair, err := ReadKeysFile(filepath.Join(dir, f.Name()))
		if err != nil {
			return nil, errors.Wrapf(err, "can't get keys from file %v", f.Name())
		}
		nodes = append(nodes, pair)
	}
	return nodes, nil
}
