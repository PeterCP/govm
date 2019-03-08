package sshutil

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"

	"golang.org/x/crypto/ssh"
)

// GenerateKeyPair generates a private and public RSA key pair, encrypted
// with `pass` (if given), and writes them to `priv` and `pub` respectively
func GenerateKeyPair(bitSize int, pass string) (*pem.Block, ssh.PublicKey, error) {
	// Generate private key
	rsaPrivKey, err := GenerateRSAPrivateKey(bitSize)
	if err != nil {
		return nil, nil, err
	}
	block := GenerateRSABlock(rsaPrivKey)
	if pass != "" {
		block, err = EncryptBlock(block, pass)
	}

	// Generate public key
	pubKey, err := ssh.NewPublicKey(&rsaPrivKey.PublicKey)
	if err != nil {
		return nil, nil, err
	}

	return block, pubKey, nil
}

func GenerateRSAPrivateKey(bitSize int) (*rsa.PrivateKey, error) {
	key, err := rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil {
		return nil, err
	}
	err = key.Validate()
	if err != nil {
		return nil, err
	}

	return key, nil
}

func GenerateRSABlock(key *rsa.PrivateKey) *pem.Block {
	return &pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   x509.MarshalPKCS1PrivateKey(key),
	}
}

func SSHPublicKey(key *rsa.PrivateKey) (ssh.PublicKey, error) {
	return ssh.NewPublicKey(&key.PublicKey)
}
