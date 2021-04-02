// cryptoHandler.go

package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"io/ioutil"
	"log"
	"os"

	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/md4"
	"golang.org/x/crypto/sha3"
)

// Base64: Structure that handle 'Encode' & 'Decode' method to/from 'Base64'.
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

// HashMe: This function handle a 'file' or '[]byte' as arguments.
// Returns the value Hash for the 'in'. The 'hashType' argument
// represents the targeted type: 'md5', 'sha256', 'sha512'.
func HashMe(in interface{}, hashType string) string {

	var (
		err    error
		hasher hash.Hash
		file   *os.File
	)
	/*
	   "md2", "MD2 hash, 128-bit (MD2)", md2.New
	   "md4", "MD4 hash, 128-bit (MD4)", md4.New
	   "md5", "MD5 hash, 128-bit (MD5)", md5.New
	   "sha1", "SHA-1 hash, 160-bit (SHA-1)", sha1.New
	   "sha2-224", "SHA-2 hash, 224-bit (SHA-224)", sha256.New224
	   "sha2-256", "SHA-2 hash, 256-bit (SHA-256)", sha256.New
	   "sha2-384", "SHA-2 hash, 384-bit (SHA-384)", sha512.New384
	   "sha2-512", "SHA-2 hash, 512-bit (SHA-512)", sha512.New
	   "sha2-512/224", "SHA-2 hash, 224-bit (SHA-512/224)", sha512.New512_224
	   "sha2-512/256", "SHA-2 hash, 256-bit (SHA-512/256)", sha512.New512_256
	   "sha3-224", "SHA-3 hash, 224-bit (SHA3-224)", sha3.New224
	   "sha3-256", "SHA-3 hash, 256-bit (SHA3-256)", sha3.New256
	   "sha3-384", "SHA-3 hash, 384-bit (SHA3-384)", sha3.New384
	   "sha3-512", "SHA-3 hash, 512-bit (SHA3-512)", sha3.New512
	   "sha3-512", "SHA-3 hash, 512-bit (SHA3-512)", sha3.New512
	   "shake-128", "SHA-3 hash, n-byte (SHAKE-128)", sha3.ShakeSum128, 32
	   "shake-256", "SHA-3 hash, n-byte (SHAKE-256)", sha3.ShakeSum256, 64
	   "ripemd-160", "RIPEMD hash, 160-bit (RIPEMD-160)", ripemd160.New
	*/
	switch hashType {
	case "md4":
		hasher = md4.New()
	case "md5":
		hasher = md5.New()
	case "sha1":
		hasher = sha1.New()
	case "sha256":
		hasher = sha256.New()
	case "sha384":
		hasher = sha512.New384()
	case "sha512":
		hasher = sha512.New()
	case "sha3-256":
		hasher = sha3.New256()
	case "sha3-384":
		hasher = sha3.New384()
	case "sha3-512":
		hasher = sha3.New512()
	case "blake2b256":
		hasher, _ = blake2b.New256(nil)
	case "blake2b384":
		hasher, _ = blake2b.New384(nil)
	case "blake2b512":
		hasher, _ = blake2b.New512(nil)
	default:
		return fmt.Sprintf("HashMe: hash type '%s' is not handled", hashType)
	}

	switch data := in.(type) {

	case string:
		if file, err = os.Open(data); err == nil {
			defer file.Close()

			if _, err = io.Copy(hasher, file); err == nil {
				return hex.EncodeToString(hasher.Sum(nil))
			}
		}

	case []byte:

		if _, err = hasher.Write(in.([]byte)); err == nil {
			return hex.EncodeToString(hasher.Sum(nil))
		}

	default:
		err = fmt.Errorf("argument is not []byte or string")

	}

	log.Printf("HashMe: '%s': %v\n", hashType, err)

	return fmt.Sprintf("HashMe: '%s': %v", hashType, err)
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

// Md5File: Get MD5 checksum from file, 'devMode' means break on error.
func Md5File(filename string, devMode ...bool) (outMd5 string) {
	var dMode bool
	if len(devMode) > 0 {
		dMode = devMode[0]
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
	if err != nil && dMode {
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
