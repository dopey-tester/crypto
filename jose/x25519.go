package jose

import (
	"crypto"
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"github.com/pkg/errors"
	"go.step.sm/crypto/x25519"
)

const x25519ThumbprintTemplate = `{"crv":"X25519","kty":"OKP","x":"%s"}`

func x25519Thumbprint(key x25519.PublicKey, hash crypto.Hash) ([]byte, error) {
	if len(key) != 32 {
		return nil, errors.New("invalid elliptic key")
	}
	h := hash.New()
	h.Write([]byte(fmt.Sprintf(x25519ThumbprintTemplate, base64.RawURLEncoding.EncodeToString(key))))
	return h.Sum(nil), nil
}

// X25519Signer implements the jose.OpaqueSigner using an X25519 key and XEdDSA
// as a signing algorithm.
type X25519Signer x25519.PrivateKey

// Public returns the public key of the current signing key.
func (s X25519Signer) Public() *JSONWebKey {
	return &JSONWebKey{
		Key: x25519.PrivateKey(s).Public(),
	}
}

// Algs returns a list of supported signing algorithms, in this case only
// XEdDSA.
func (s X25519Signer) Algs() []SignatureAlgorithm {
	return []SignatureAlgorithm{
		XEdDSA,
	}
}

// SignPayload signs a payload with the current signing key using the given
// algorithm, it will fails if it's not XEdDSA.
func (s X25519Signer) SignPayload(payload []byte, alg SignatureAlgorithm) ([]byte, error) {
	if alg != XEdDSA {
		return nil, errors.Errorf("x25519 key does not support the signature algorithm %s", alg)
	}
	return x25519.PrivateKey(s).Sign(rand.Reader, payload, crypto.Hash(0))
}

// X25519Verifier implements the jose.OpaqueVerifier interface using an X25519
// key and XEdDSA as a signing algorithm.
type X25519Verifier x25519.PublicKey

func (v X25519Verifier) VerifyPayload(payload []byte, signature []byte, alg SignatureAlgorithm) error {
	if alg != XEdDSA {
		return errors.Errorf("x25519 key does not support the signature algorithm %s", alg)
	}
	if !x25519.Verify(x25519.PublicKey(v), payload, signature) {
		return errors.New("failed to verify XEdDSA signature")
	}
	return nil
}
