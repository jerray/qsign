package qsign

import (
	"crypto/md5"
	"encoding/hex"
	"hash"
)

// Generator is function returns string. It's used to generate digest prefix or suffix.
type Generator func() string

// Filter is function receives a key-value pair, returns bool value. It's used to filter out
// invalid key-value pair in the digest.
type Filter func(key, value string) bool

func defaultFilter(key, value string) bool {
	return len(value) > 0
}

type Marshaler interface {
	MarshalQsign() string
}

// Hasher is function returns hash.Hash. By default, Qsign uses md5 hash. You can provide
// your own hash interface instead.
type Hasher func() hash.Hash

func defaultHasher() hash.Hash {
	return md5.New()
}

// Encoding is an interface for various encoding scheme.
type Encoding interface {

	// Encode encodes src using the encoding enc, writing EncodedLen(len(src)) bytes to dst.
	Encode(dst, src []byte)

	// EncodedLen returns the length in bytes of the encoding scheme of an input buffer of length n.
	EncodedLen(n int) int
}

// Encoder is a function returns Encoding interface for Qsign to encode digest.
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
