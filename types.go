package qsign

import (
	"crypto/md5"
	"encoding/hex"
	"hash"
)

type Generator func() string

type Filter func(key, value string) bool

func defaultFilter(key, value string) bool {
	return len(value) > 0
}

type Hasher func() hash.Hash

func defaultHasher() hash.Hash {
	return md5.New()
}

type Encoding interface {
	Encode(dst, src []byte)
	EncodedLen(n int) int
}

type Encoder func() Encoding

type hexEncoding struct{}

func (h *hexEncoding) Encode(dst, src []byte) {
	hex.Encode(dst, src)
}

func (h *hexEncoding) EncodedLen(n int) int {
	return hex.EncodedLen(n)
}

var defaultHexEncoding = &hexEncoding{}

func defaultEncoder() Encoding {
	return defaultHexEncoding
}
