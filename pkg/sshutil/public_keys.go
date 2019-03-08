package sshutil

import (
	"io"

	"golang.org/x/crypto/ssh"
)

func PublicKeyBytes(key ssh.PublicKey) []byte {
	return ssh.MarshalAuthorizedKey(key)
}

func WritePublicKey(w io.Writer, key ssh.PublicKey) error {
	_, err := w.Write(ssh.MarshalAuthorizedKey(key))
	return err
}
