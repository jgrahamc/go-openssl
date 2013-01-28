// Program that take expects to find a file called 'test.txt' which it
// will MD5 and SHA1 1,000 times and will encrypt 1,000 times with
// RC4. Outputs timing information.

package main

import (
	"hash"
	"crypto/md5"
	"crypto/sha1"
	"crypto/rc4"
	"fmt"
	"io/ioutil"
	"time"
)

const testFile = "test.txt"
const iterations = 1000

func timed(name string, size float64, f func()) {
	start := time.Now()
	f()
	elapsed := time.Since(start).Seconds()
	speed := size / elapsed
	fmt.Printf("%s %.2fs %.fMB/s\n", name, elapsed, speed / 1000000)
}

func h(name string, hasher hash.Hash, t []byte) {
	timed(name, float64(iterations * len(t)), func() {
		for i := 0; i < iterations; i++ {
			hasher.Write(t)
		}
		hasher.Sum(nil)
	})
}

func e(t []byte) {
	if c, err := rc4.NewCipher([]byte("8235672ru9203849234932")); err != nil {
		fmt.Printf("Unable to create RC4 cipher: %s\n", err)
	} else {
		timed("rc4", float64(iterations * len(t)), func() {
			dst := make([]byte, len(t))
			for i := 0; i < iterations; i++ {
				c.XORKeyStream(dst, t)
			}
		})
	}
}

func main() {
	if t, err := ioutil.ReadFile(testFile); err != nil {
		fmt.Printf("Unable to read test file %s: %s\n", testFile, err)
	} else {
		h("md5", md5.New(), t)
		h("sha1", sha1.New(), t)
		e(t)
	}
}
