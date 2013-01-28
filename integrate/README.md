integrate
=========

Running 'make' in this directory will build a local version of Go
based on the latest Go in Mercurial that substitutes the md5, sha1 and
rc4 implementations from this project. Note that it does not use the
versions of md5, sha1 and rc4 from the directories above; it pulls
them directly from github.

Running 'make tester' builds a small benchmarking program.

