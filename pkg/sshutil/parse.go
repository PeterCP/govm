package sshutil

import (
	"encoding/pem"
	"fmt"

	"golang.org/x/crypto/ssh"
)

func ParseBlock(b []byte) (*pem.Block, error) {
	block, _ := pem.Decode(b)
	if block == nil {
		return nil, fmt.Errorf("invalid PEM block")
	}
	return block, nil
}

func ParseSigner(b []byte) (ssh.Signer, error) {
	return ssh.ParsePrivateKey(b)
}

func ParsePublicKey(keyBytes []byte) (ssh.PublicKey, error) {
	return ssh.ParsePublicKey(keyBytes)
}
