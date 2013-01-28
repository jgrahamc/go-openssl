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

copy = cp $(GOOPENSSL)$1/$1.go $(GOCRYPTO)$1

MD5 := $(GOCRYPTO)md5/
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
endef
export NEW_EXCEPTIONS

.fixup-go: .get-go .get-go-openssl
.fixup-go:
	@cp -r $(MD5) $(GOMD5)
	@sed -i~ -e 's/md5/gomd5/' $(addprefix $(GOMD5),md5.go md5block.go md5_test.go) $(GOSRC)cmd/cgo/main.go
	@sed -i~ -e 's/"crypto\/\(md5\|sha1\|rc4\)":\( *\){"L3"},/"crypto\/\1":\2{"L3", "C"},/' $(GOSRC)pkg/go/build/deps_test.go
	@echo "$$NEW_EXCEPTIONS" >> $(GO)api/except.txt
	@rm $(addprefix $(GOCRYPTO),md5/md5block.go sha1/sha1block.go)
	@$(call copy,md5)
	@$(call copy,sha1)
	@$(call copy,rc4)
	@touch $@

.get-go: ; @rm -rf $(GO) && hg clone https://code.google.com/p/go $(GO) && touch $@
.get-go-openssl: ; @rm -rf $(GOOPENSSL) && git clone git://github.com/jgrahamc/go-openssl.git $(GOOPENSSL) && touch $@

.PHONY: clean
clean: ; @rm -f .get-go .get-go-openssl .fixup-go .build-go

.PHONY: tester
tester: .build-go
tester: ; @$(GO)bin/go build tester.go
