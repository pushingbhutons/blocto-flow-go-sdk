/*
 * Flow Go SDK
 *
 * Copyright 2019-2020 Dapper Labs, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package crypto

import (
	"fmt"

	"github.com/portto/blocto-flow-go-sdk/crypto/internal/crypto/hash"
)

// Signer interface
type signer interface {
	// generatePrKey generates a private key
	generatePrivateKey([]byte) (PrivateKey, error)
	// decodePrKey loads a private key from a byte array
	decodePrivateKey([]byte) (PrivateKey, error)
	// decodePubKey loads a public key from a byte array
	decodePublicKey([]byte) (PublicKey, error)
}

// newSigner initializes a new signer from the given signing algorithm.
func newSigner(algo SigningAlgorithm) (signer, error) {
	switch algo {
	case ECDSAP256:
		return newECDSAP256(), nil
	case ECDSASecp256k1:
		return newECDSASecp256k1(), nil
	default:
		return nil, fmt.Errorf("the signature scheme %s is not supported.", algo)
	}
}

// GeneratePrivateKey generates a private key of the algorithm using the entropy of the given seed
func GeneratePrivateKey(algo SigningAlgorithm, seed []byte) (PrivateKey, error) {
	signer, err := newSigner(algo)
	if err != nil {
		return nil, fmt.Errorf("key generation failed: %w", err)
	}
	return signer.generatePrivateKey(seed)
}

// DecodePrivateKey decodes an array of bytes into a private key of the given algorithm
func DecodePrivateKey(algo SigningAlgorithm, data []byte) (PrivateKey, error) {
	signer, err := newSigner(algo)
	if err != nil {
		return nil, fmt.Errorf("decode private key failed: %w", err)
	}
	return signer.decodePrivateKey(data)
}

// DecodePublicKey decodes an array of bytes into a public key of the given algorithm
func DecodePublicKey(algo SigningAlgorithm, data []byte) (PublicKey, error) {
	signer, err := newSigner(algo)
	if err != nil {
		return nil, fmt.Errorf("decode public key failed: %w", err)
	}
	return signer.decodePublicKey(data)
}

// Signature type tools

// Bytes returns a byte array of the signature data
func (s Signature) Bytes() []byte {
	return s[:]
}

// String returns a String representation of the signature data
func (s Signature) String() string {
	return fmt.Sprintf("%#x", s.Bytes())
}

// Key Pair

// PrivateKey is an unspecified signature scheme private key
type PrivateKey interface {
	// Algorithm returns the signing algorithm related to the private key.
	Algorithm() SigningAlgorithm
	// Size return the key size in bytes.
	Size() int
	// String return a hex representation of the key
	String() string
	// Sign generates a signature using the provided hasher.
	Sign([]byte, hash.Hasher) (Signature, error)
	// PublicKey returns the public key.
	PublicKey() PublicKey
	// Encode returns a bytes representation of the private key
	Encode() []byte
	// Equals returns true if the given PrivateKeys are equal. Keys are considered unequal if their algorithms are
	// unequal or if their encoded representations are unequal. If the encoding of either key fails, they are considered
	// unequal as well.
	Equals(PrivateKey) bool
}

// PublicKey is an unspecified signature scheme public key.
type PublicKey interface {
	// Algorithm returns the signing algorithm related to the public key.
	Algorithm() SigningAlgorithm
	// Size() return the key size in bytes.
	Size() int
	// String return a hex representation of the key
	String() string
	// Verify verifies a signature of an input message using the provided hasher.
	Verify(Signature, []byte, hash.Hasher) (bool, error)
	// Encode returns a bytes representation of the public key.
	Encode() []byte
	// Equals returns true if the given PublicKeys are equal. Keys are considered unequal if their algorithms are
	// unequal or if their encoded representations are unequal. If the encoding of either key fails, they are considered
	// unequal as well.
	Equals(PublicKey) bool
}
