# Makefile that gets the latest version of Go using hg and then
# substitutes the MD5, SHA1 and RC4 implementations for OpenSSL Doing
# this requires a trick to rename the native Go implementation of MD5
# to nativemd5 and use it to build cgo otherwise the build hangs with
# a circular dependency (md5 needs cgo and cgo needs md5).

ROOT := 
GO := $(ROOT)go/
GOOPENSSL := $(ROOT)go-openssl/

.PHONY: all
all: .build-go

.build-go: .fixup-go
.build-go:
	@cd $(GO)src && ./all.bash
	@touch $@

GOSRC := $(GO)src/
GOCRYPTO := $(GOSRC)pkg/crypto/

copy = cp $(GOOPENSSL)$(if $2,$2,$1)/$1.go $(GOCRYPTO)$1/$1.go

MD5 := $(GOCRYPTO)md5/
SHA1 := $(GOCRYPTO)sha1/
RC4 := $(GOCRYPTO)rc4/
GOMD5 := $(GOCRYPTO)gomd5/

# This target modifies the Go installation in the following ways:
#
# 1. Copies the native Go md5 implementation and calls it crypto/gomd5
#
# 2. Alters cgo so that it uses crypto/gomd5 and fixes gomd5 so that
# it doesn't reference crypto/md5
#
# 3. Alters the deps_test.go to recognize that the OpenSSL-linked
#    versions of md5, sha1 and rc4 are dependent on the C include.
#
# 4. Modifies the api/except.txt list to cause the change in signature
#    on the md5, sha1 and rc4 packages to be ignored when testing API
#    compatibility (this is necessary because using cgo changes the
#    signature even though they are compatible with the existing
#    packages)
#
# 5. Deletes the md5block.go and sha1block.go implementations since
#    they are not needed.
#
# 6. Copies in the OpenSSL md5, sha1 and rc4.
#
# 7. Modifies rc4.go to use either C.ulong or C.size_t depending on
#    how the RC4 function is defined by OpenSSL
#
# 8. Does a bunch of work if the operating system is Darwin (i.e Mac
#    OS X) because Apple (in their infinite wisdom) have decided to
#    deprecate OpenSSL.

define NEW_EXCEPTIONS
pkg crypto/md5, const BlockSize ideal-int
pkg crypto/md5, const Size ideal-int
pkg crypto/md5, func New() hash.Hash
pkg crypto/sha1, const BlockSize ideal-int
pkg crypto/sha1, const Size ideal-int
pkg crypto/sha1, func New() hash.Hash
pkg crypto/rc4, func NewCipher([]byte) (*Cipher, error)
pkg crypto/rc4, method (*Cipher) Reset()
pkg crypto/rc4, method (*Cipher) XORKeyStream([]byte, []byte)
pkg crypto/rc4, method (KeySizeError) Error() string
pkg crypto/rc4, type Cipher struct
pkg crypto/rc4, type KeySizeError int
pkg crypto/rc4, func NewCipher([]uint8) (*Cipher, error)
pkg crypto/rc4, method (*Cipher) XORKeyStream([]uint8, []uint8)
endef
export NEW_EXCEPTIONS

IS_DARWIN :=
RC4SIZE_T :=

UNAME := $(shell uname)

ifeq ($(UNAME),Darwin)
IS_DARWIN := T
else
RC4SIZE_T := $(shell grep size_t /usr/include/openssl/rc4.h)
endif

.fixup-go: .get-go .get-go-openssl
.fixup-go: .movemd5 .fixapi .copyopenssl
.movemd5:
	@cp -r $(MD5) $(GOMD5)
	@sed -i~ -e 's/md5/gomd5/' $(addprefix $(GOMD5),md5.go md5block.go md5_test.go) $(GOSRC)cmd/cgo/main.go
	@touch $@

DEPS := $(GOSRC)pkg/go/build/deps_test.go

.fixapi:
	@sed -i~ -e 's/"crypto\/md5":\( *\){"L3"},/"crypto\/md5":\1{"L3", "C"},/' $(DEPS)
	@sed -i~ -e 's/"crypto\/sha1":\( *\){"L3"},/"crypto\/sha1":\1{"L3", "C"},/' $(DEPS)
	@sed -i~ -e 's/"crypto\/rc4":\( *\){"L3"},/"crypto\/rc4":\1{"L3", "C"},/' $(DEPS)
	@echo "$$NEW_EXCEPTIONS" >> $(GO)api/except.txt
	@touch $@

.copyopenssl:
	@rm -f $(addprefix $(GOCRYPTO),md5/md5block.go sha1/sha1block.go)
	@rm -f $(addprefix $(GOCRYPTO)rc4/,rc4_amd64.s rc4_asm.go rc4_ref.go)
	@$(call copy,md5)
	@$(call copy,sha1)
	@$(if $(IS_DARWIN),$(call copy,rc4_darwin,rc4),$(call copy,rc4))
	@$(if $(IS_DARWIN),sed -i~ -e 's/C.size_t/_Ctype_CC_LONG/' $(MD5)md5.go $(SHA1)sha1.go)
	@$(if $(IS_DARWIN),sed -i~ -e 's/_Ctype_MD5_CTX/_Ctype_CC_MD5_CTX/' $(MD5)md5.go)
	@$(if $(IS_DARWIN),sed -i~ -e 's/_Ctype_SHA_CTX/_Ctype_CC_SHA1_CTX/' $(SHA1)sha1.go)
	@$(if $(RC4SIZE_T),,sed -i~ -e 's/size_t/ulong/' $(RC4)rc4.go)
	@touch $@

ifneq ($(REV),)
override REV := --rev $(REV)
endif

.get-go: ; @rm -rf $(GO) && hg clone $(REV) https://code.google.com/p/go $(GO) && touch $@
.get-go-openssl: ; @rm -rf $(GOOPENSSL) && git clone git://github.com/jgrahamc/go-openssl.git $(GOOPENSSL) && touch $@

.PHONY: clean
clean: ; @rm -f .get-go .get-go-openssl .fixup-go .build-go .copyopenssl .fixapi .movemd5

.PHONY: tester
tester: .build-go
tester: ; @$(GO)bin/go build tester.go
