package waarp

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"encoding/hex"
	"fmt"
	"io/ioutil"
)

type DesKey struct {
	b cipher.Block
}

func (d *DesKey) Encrypt(s string) string {
	return fmt.Sprintf("%x", encrypt(s, d.b))
}

func NewDesKey(keyFile string) (*DesKey, error) {
	b, err := makeBlock(keyFile)
	if err != nil {
		return nil, err
	}
	return &DesKey{b}, nil
}

func NewDesKeyFromBytes(key []byte) (*DesKey, error) {
	b, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	return &DesKey{b}, nil
}

func EncryptPassword(keyPath string, pwd string) (string, error) {
	des, err := NewDesKey(keyPath)
	if err != nil {
		return "", err
	}
	return des.Encrypt(pwd), nil
}

// func testDecrypt(b cipher.Block) {
//  fmt.Println("\nTesting decryption")
//  ciphertext := "d7eb3b3d0ce0ef0cc8105214a8edf429"

//  fmt.Printf("%s\n", Decrypt(ciphertext, b))
// }

func decrypt(h string, b cipher.Block) []byte {
	ciphertext, _ := hex.DecodeString(h)

	mode := newECBDecrypter(b)
	mode.CryptBlocks(ciphertext, ciphertext)
	return unPkcs5Padding(ciphertext)
}

// func testEncrypt(b cipher.Block) {
//  fmt.Println("\nTesting encryption")
//  plaintext := "password"

//  fmt.Printf("%x\n", Encrypt(plaintext, b))
// }

func encrypt(p string, b cipher.Block) []byte {
	plaintext := pkcs5Padding([]byte(p))
	ciphertext := make([]byte, len(plaintext))

	mode := newECBEncrypter(b)
	mode.CryptBlocks(ciphertext, plaintext)
	return ciphertext
}

// func testEncryptDecrypt(plaintext string, b cipher.Block) {
//  fmt.Println("\nTesting encryption/decryption")
//  fmt.Printf("original:  %s\n", plaintext)

//  encrypted := Encrypt(plaintext, b)
//  encryptedHex := hex.EncodeToString(encrypted)
//  fmt.Printf("encrypted (hex): %s\n", encryptedHex)

//  decrypted := Decrypt(encryptedHex, b)

//  fmt.Printf("decrypted: %s\n", decrypted)
// }

func makeBlock(keyPath string) (cipher.Block, error) {
	key, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}

	b, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	return b, nil
}

////
//
// Libraries found elsewhere
//
////

// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Electronic Code Book (ECB) mode.

// ECB provides confidentiality by assigning a fixed ciphertext block to each
// plaintext block.

// See NIST SP 800-38A, pp 08-09

// package cipher

type ecb struct {
	b         cipher.Block
	blockSize int
}

func newECB(b cipher.Block) *ecb {
	return &ecb{
		b:         b,
		blockSize: b.BlockSize(),
	}
}

type ecbEncrypter ecb

// NewECBEncrypter returns a BlockMode which encrypts in electronic code book
// mode, using the given Block.
func newECBEncrypter(b cipher.Block) cipher.BlockMode {
	return (*ecbEncrypter)(newECB(b))
}

func (x *ecbEncrypter) BlockSize() int { return x.blockSize }

func (x *ecbEncrypter) CryptBlocks(dst, src []byte) {
	if len(src)%x.blockSize != 0 {
		panic("crypto/cipher: input not full blocks")
	}
	if len(dst) < len(src) {
		panic("crypto/cipher: output smaller than input")
	}
	for len(src) > 0 {
		x.b.Encrypt(dst, src[:x.blockSize])
		src = src[x.blockSize:]
		dst = dst[x.blockSize:]
	}
}

type ecbDecrypter ecb

// NewECBDecrypter returns a BlockMode which decrypts in electronic code book
// mode, using the given Block.
func newECBDecrypter(b cipher.Block) cipher.BlockMode {
	return (*ecbDecrypter)(newECB(b))
}

func (x *ecbDecrypter) BlockSize() int { return x.blockSize }

func (x *ecbDecrypter) CryptBlocks(dst, src []byte) {
	if len(src)%x.blockSize != 0 {
		panic("crypto/cipher: input not full blocks")
	}
	if len(dst) < len(src) {
		panic("crypto/cipher: output smaller than input")
	}
	for len(src) > 0 {
		x.b.Decrypt(dst, src[:x.blockSize])
		src = src[x.blockSize:]
		dst = dst[x.blockSize:]
	}
}

func pkcs5Padding(data []byte) []byte {
	var newData []byte

	mod := len(data) % 8
	lenAdd := 8 - mod
	complement := bytes.Repeat([]byte{byte(lenAdd)}, lenAdd)
	newData = append(data, complement...)
	return newData
}

func unPkcs5Padding(data []byte) []byte {
	dataLen := len(data)

	endIndex := int(data[dataLen-1])

	if 8 >= endIndex {

		if 1 < endIndex {
			for i := dataLen - endIndex; i < dataLen; i++ {
				if data[dataLen-1] != data[i] {
					fmt.Println("不同的字节码，尾部字节码:", data[dataLen-1], "  下标：", i, "  字节码：", data[i])
				}
			}
		}

		return data[:dataLen-endIndex]
	}

	return nil
}
