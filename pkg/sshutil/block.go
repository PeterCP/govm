package sshutil

import (
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"io"

	"golang.org/x/crypto/ssh"
)

func IsBlockEncrypted(block *pem.Block) bool {
	return x509.IsEncryptedPEMBlock(block)
}

func EncryptBlock(block *pem.Block, pass string) (*pem.Block, error) {
	return x509.EncryptPEMBlock(rand.Reader, block.Type,
		block.Bytes, []byte(pass), x509.PEMCipherAES256)
}

func DecryptBlock(block *pem.Block, pass string) (*pem.Block, error) {
	decrypted, err := x509.DecryptPEMBlock(block, []byte(pass))
	if err != nil {
		return nil, err
	}
	return &pem.Block{
		Type:    block.Type,
		Headers: nil,
		Bytes:   decrypted,
	}, nil
}

func BlockToSigner(block *pem.Block) (ssh.Signer, error) {
	return ParseSigner(pem.EncodeToMemory(block))
}

func WriteBlock(w io.Writer, block *pem.Block) error {
	return pem.Encode(w, block)
}

func BlockBytes(key *pem.Block) []byte {
	return pem.EncodeToMemory(key)
}
