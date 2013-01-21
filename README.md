go-openssl
==========

Go wrapper for some OpenSSL functions in libcrypto

Currently provides a wrapper for MD5 and SHA1 hashes. These are
considerably faster than the native Go versions.

My benchmarks have OpenSSL MD5 as 1.5x faster and SHA1 as 5x
faster on AMD64.

The interface is almost the same as crypto/md5 and crypto/sha1 
except that Sum() does the equivalent of a Reset() when called.
