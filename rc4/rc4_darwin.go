package rc4

// #include <CommonCrypto/CommonCryptor.h>
import "C"

import (
	"strconv"
	"unsafe"
)

const reserve = 256

type Cipher struct {
	memory [reserve]byte
	state  C.CCCryptorRef
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
	var used C.size_t
	C.CCCryptorCreateFromData(C.kCCEncrypt, C.kCCAlgorithmRC4, 0,
		unsafe.Pointer(&key[0]), C.size_t(k), nil, unsafe.Pointer(&c.memory), reserve,
		&c.state, &used)

	return &c, nil
}

func (c *Cipher) XORKeyStream(dst, src []byte) {
	var moved C.size_t
	C.CCCryptorUpdate(c.state, unsafe.Pointer(&src[0]), C.size_t(len(src)),
		unsafe.Pointer(&dst[0]), C.size_t(len(dst)), &moved)
}

func (c *Cipher) Reset() {
     C.CCCryptorReset(c.state, nil)
}
