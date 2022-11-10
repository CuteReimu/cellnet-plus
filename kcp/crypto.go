package kcp

import (
	"crypto/sha1"

	"github.com/xtaci/kcp-go/v5"
	"golang.org/x/crypto/pbkdf2"
)

var blockCrypto kcp.BlockCrypt

func init() {
	const (
		defaultPass = "2bc2af69dbfb7dcb0985a66d76b226b9"
		defaultSalt = "b36ca0200a07b6d9d15784a78515059e"
	)
	key := pbkdf2.Key([]byte(defaultPass), []byte(defaultSalt), 1024, 32, sha1.New)
	var err error
	if blockCrypto, err = kcp.NewAESBlockCrypt(key); err != nil {
		panic(err)
	}
}
