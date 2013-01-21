package sha1

// #cgo LDFLAGS: -lcrypto
// #include <openssl/sha.h>
import "C"

import (
	"hash"
	"errors"
	"unsafe"
)

// The size of an SHA1 checksum in bytes.
const Size = 20

// The blocksize of SHA1 in bytes.
const BlockSize = 64

type digest struct {
	context *_Ctype_SHA_CTX
}

// New returns a new hash.Hash computing the SHA1 checksum.
func New() hash.Hash {
	d := new(digest)
	d.Reset()
	return d
}

func (d *digest) Reset() {
	d.context = &_Ctype_SHA_CTX{}
	C.SHA1_Init(d.context)
}

func (d *digest) Size() int { return Size }

func (d *digest) BlockSize() int { return BlockSize }

func (d *digest) Write(p []byte) (nn int, err error) {
	if C.SHA1_Update(d.context, unsafe.Pointer(&p[0]), (C.size_t)(len(p))) == 1 {
		return len(p), nil
	}
		
	return 0, errors.New("SHA1_Update failed")
}

// Note that this is different from the native Go package sha1.Sum()
// because it resets the hash and so an sha1.Write() after calling this
// will be starting with a reset hash.
func (d *digest) Sum(in []byte) []byte {
	md := make([]byte, Size)
	if C.SHA1_Final((*_Ctype_unsignedchar)(&md[0]), d.context) == 1 {
		return append(in, md...)
	}

	return nil
}
