package qsign

import (
	"bytes"
)

// Qsign is the signer which signs structs.
type Qsign struct {
	prefixGenerator Generator
	suffixGenerator Generator
	encoder         Encoder
	filter          Filter
	hasher          Hasher
	delimiter       string
	connector       string
}

// Options is optional attributes for building NewSign function to build *Qsign.
//
// PrefixGenerator and SuffixGenerator is two functions which you can use to generating
// prefix string prepending to digest and suffix string appending to digest string.
//
// Filter is a function used to get rid of some keys or values. For example you want
// a field which its value is empty being ignored. And this is the default filter.
//
// Encoder is a function which returns Encoding interface. By default it returns hex
// encoding. If you want to use base64 encoding, you can give a function which returns
// base64.StdEncoding.
//
// Hasher is a function which returns hash.Hash. By default it returns the Hash from
// crypto/md5.
type Options struct {
	PrefixGenerator Generator
	SuffixGenerator Generator
	Encoder         Encoder
	Filter          Filter
	Hasher          Hasher
}

// NewQsign returns a new *Qsign computing signature.
func NewQsign(options *Options) *Qsign {
	filter := options.Filter
	if filter == nil {
		filter = defaultFilter
	}

	hasher := options.Hasher
	if hasher == nil {
		hasher = defaultHasher
	}

	encoder := options.Encoder
	if encoder == nil {
		encoder = defaultEncoder
	}

	q := &Qsign{
		prefixGenerator: options.PrefixGenerator,
		suffixGenerator: options.SuffixGenerator,
		encoder:         encoder,
		filter:          filter,
		hasher:          hasher,
		delimiter:       "&",
		connector:       "=",
	}

	return q
}

// Sign returns signature bytes for interface v. It calculate the digest of input struct first. And
// then gets checksum of the digest using hasher. Finally encodes the checksum and returns.
func (q *Qsign) Sign(v interface{}) ([]byte, error) {
	digest, err := q.Digest(v)
	if err != nil {
		return nil, err
	}

	h := q.hasher()
	h.Write(digest)

	e := q.encoder()
	dst := make([]byte, e.EncodedLen(h.Size()))
	e.Encode(dst, h.Sum(nil))

	return dst, nil
}

// Digest generates digest bytes for interface v. By default, it parses struct v, gets all the
// keys and values, and connects them like an HTTP query string.
//
// Key's value is struct field name if there is no tags like "qsign", "json", "yaml" or "xml". If
// any key has a tag mentioned before, it will get value from the tag for that key. Tag name "qsign"
// has the highest priority. Field tag valuewith "-" will be ignored.
//
// All the values expect for Array, Slice and Struct type will be parsed to string. There is an
// exception here, if the struct has a String method (`func String() string`), it will be parsed.
func (q *Qsign) Digest(v interface{}) ([]byte, error) {
	buf := new(bytes.Buffer)

	if q.prefixGenerator != nil {
		if _, err := buf.WriteString(q.prefixGenerator()); err != nil {
			return buf.Bytes(), err
		}
	}

	vs := getStructValues(v)
	l := len(vs) - 1
	for i, f := range vs {
		if !q.filter(f.name, f.value) {
			continue
		}

		buf.WriteString(f.name)
		buf.WriteString(q.connector)
		buf.WriteString(f.value)

		if i != l {
			buf.WriteString(q.delimiter)
		}
	}

	if q.suffixGenerator != nil {
		if _, err := buf.WriteString(q.suffixGenerator()); err != nil {
			return buf.Bytes(), err
		}
	}

	return buf.Bytes(), nil
}

// SetDelimiter changes the default delimiter.
func (q *Qsign) SetDelimiter(s string) {
	q.delimiter = s
}

// SetConnector changes the default connector.
func (q *Qsign) SetConnector(s string) {
	q.connector = s
}
