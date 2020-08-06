// cryptoHandler.go

package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"io"
	"io/ioutil"
	"log"
	"os"
)

type Base64 struct {
}

func (b *Base64) Encode(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

func (b *Base64) Decode(encoded string) (decoded []byte) {
	var err error
	if decoded, err = base64.StdEncoding.DecodeString(encoded); err != nil {
		log.Fatalln(err)
	}
	return
}

// Md5String: Get MD5 checksum from string.
func Md5String(inString string) (outMd5 string) {
	var err error
	hasher := md5.New()
	if _, err = hasher.Write([]byte(inString)); err == nil {
		outMd5 = hex.EncodeToString(hasher.Sum(nil))
	} else {
		log.Fatalln(err.Error())
	}
	return
}

// Md5File: Get MD5 checksum from file.
func Md5File(filename string, noError ...bool) (outMd5 string) {
	var devMode bool
	if len(noError) > 0 {
		devMode = noError[0]
	}
	var err error
	var file *os.File
	if file, err = os.Open(filename); err == nil {
		defer file.Close()
		hasher := md5.New()
		if _, err = io.Copy(hasher, file); err == nil {
			outMd5 = hex.EncodeToString(hasher.Sum(nil))
		}
	}
	if err != nil && devMode {
		log.Fatalln(err.Error())
	}
	return
}

/*
	Frankenstein algorithm, code based on a web article by Nic Raboy:
	https://www.thepolyglotdeveloper.com/2018/02/encrypt-decrypt-data-golang-application-crypto-packages/
	thanks to him.
*/

// AES256CipherStruct: Structure that contain information for encrypting/
// decrypting file or data. Result of string encryption is contained in Data
type AES256CipherStruct struct {
	Data        []byte // In/Out data store.
	InFilename  string
	OutFilename string

	passphrase string
}

func (c *AES256CipherStruct) Encrypt(passphrase string) error {
	c.passphrase = passphrase
	return c.encrypt()
}
func (c *AES256CipherStruct) Uncrypt(passphrase string) error {
	c.passphrase = passphrase
	return c.decrypt()
}
func (c *AES256CipherStruct) EncryptFile(passphrase string) (err error) {
	c.passphrase = passphrase
	if c.Data, err = ioutil.ReadFile(c.InFilename); err == nil {
		err = c.encryptFile()
	}
	return
}
func (c *AES256CipherStruct) UncryptFile(passphrase string) (err error) {
	c.passphrase = passphrase
	if err = c.decryptFile(); err == nil {
		err = ioutil.WriteFile(c.OutFilename, c.Data, os.ModePerm)
	}
	return
}

// createHash:
func (c *AES256CipherStruct) createHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}

// encrypt:
func (c *AES256CipherStruct) encrypt() (err error) {
	var gcm cipher.AEAD
	var block cipher.Block
	if len(c.Data) > 0 {
		if block, err = aes.NewCipher([]byte(c.createHash(c.passphrase))); err == nil {
			if gcm, err = cipher.NewGCM(block); err == nil {
				nonce := make([]byte, gcm.NonceSize())
				if _, err = io.ReadFull(rand.Reader, nonce); err == nil {
					c.Data = gcm.Seal(nonce, nonce, c.Data, nil)
				}
			}
		}
	}
	return
}

// decrypt:
func (c *AES256CipherStruct) decrypt() (err error) {
	var gcm cipher.AEAD
	var block cipher.Block
	if len(c.Data) > 0 {
		key := []byte(c.createHash(c.passphrase))
		if block, err = aes.NewCipher(key); err == nil {
			if gcm, err = cipher.NewGCM(block); err == nil {
				nonceSize := gcm.NonceSize()
				nonce, ciphertext := c.Data[:nonceSize], c.Data[nonceSize:]
				c.Data, err = gcm.Open(nil, nonce, ciphertext, nil)
			}
		}
	}
	return
}

// encryptFile:
func (c *AES256CipherStruct) encryptFile() (err error) {
	var f *os.File
	if f, err = os.Create(c.OutFilename); err == nil {
		defer f.Close()
		if err = c.encrypt(); err == nil {
			if _, err = f.Write(c.Data); err == nil {
				err = f.Sync()
			}
		}
	}
	return
}

// decryptFile:
func (c *AES256CipherStruct) decryptFile() (err error) {
	if c.Data, err = ioutil.ReadFile(c.InFilename); err == nil {
		err = c.decrypt()
	}
	return
}

/* TESTING AES256CipherStruct

func main() {
	var err error
	filename := "--filename--"
	pp := "1235"

	c := new(gcr.AES256CipherStruct)
	c.InFilename = filename
	c.OutFilename = filename + "-c"
	err = c.EncryptFile(pp)
	if err != nil {
		panic(err)
	}
	c.InFilename = filename + "-c"
	c.OutFilename = filename + "-d"
	err = c.UncryptFile(pp)
	if err != nil {
		panic(err)
	}
	c.Data = []byte(filename)
	err = c.Encrypt(pp)
	if err != nil {
		panic(err)
	}
	err = c.Uncrypt(pp)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(c.Data))
}

*/
