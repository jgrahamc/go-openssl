go-openssl
==========

Go wrapper for some OpenSSL functions in libcrypto

Currently provides a wrapper for MD5 and SHA1 hashes and the RC4
encryption algorithm. These are considerably faster than the native Go
versions.

My benchmarks have OpenSSL MD5 as 1.5x faster and SHA1 as 5x
faster on AMD64. RC4 is 7x faster using OpenSSL.

