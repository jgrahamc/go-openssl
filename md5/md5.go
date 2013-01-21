package md5

// #cgo LDFLAGS: -lcrypto
// #include <openssl/md5.h>
import "C"

import (
	"hash"
	"errors"
	"unsafe"
)

// The size of an MD5 checksum in bytes.
const Size = 16

// The blocksize of MD5 in bytes.
const BlockSize = 64

type digest struct {
	context *_Ctype_MD5_CTX
}

// New returns a new hash.Hash computing the MD5 checksum.
func New() hash.Hash {
	d := new(digest)
	d.Reset()
	return d
}

func (d *digest) Reset() {
	d.context = &_Ctype_MD5_CTX{}
	C.MD5_Init(d.context)
}

func (d *digest) Size() int { return Size }

func (d *digest) BlockSize() int { return BlockSize }

func (d *digest) Write(p []byte) (nn int, err error) {
	if C.MD5_Update(d.context, unsafe.Pointer(&p[0]), (C.size_t)(len(p))) == 1 {
		return len(p), nil
	}
		
	return 0, errors.New("MD5_Update failed")
}

func (d *digest) Sum(in []byte) []byte {
	context := *d.context
	defer func() {*d.context = context}()

	md := make([]byte, Size)
	if C.MD5_Final((*_Ctype_unsignedchar)(&md[0]), d.context) == 1 {
		return append(in, md...)
	}

	return nil
}
