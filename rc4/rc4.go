package rc4

// #cgo LDFLAGS: -lcrypto
// #include <openssl/rc4.h>

import "C"

import (
	"strconv"
)

type Cipher struct {
	key *_Ctype_RC4_KEY
}

type KeySizeError int

func (k KeySizeError) Error() string {
	return "rc4: invalid key size " + strconv.Itoa(int(k))
}

func NewCipher(key []byte) (*Cipher, error) {
	k := len(key)
	if k < 1 || k > 256 {
		return nil, KeySizeError(k)
	}
	var c Cipher
	c.key = &_Ctype_RC4_KEY{}
	C.RC4_set_key(c.key, C.int(k), (*_Ctype_unsignedchar)(&key[0]))

	return &c, nil
}

func (c *Cipher) XORKeyStream(dst, src []byte) {
	C.RC4(c.key, C.size_t(len(dst)), (*_Ctype_unsignedchar)(&src[0]),
		(*_Ctype_unsignedchar)(&dst[0]))
}

func (c *Cipher) Reset() {
	for i := 0; i < 256; i++ {
		c.key.data[i] = 0
	}
	c.key.x = 0
	c.key.y = 0
}
