// Package vigenere provides a secure symmetric encryption implementation with Vigenère-Vernam Cipher algorithm. The Vigenère Cipher
// in this library also acts as a One-Time Pad (OTP) for additional security. Vigenère Cipher is done by defining a character set or
// alphabets to be used. Then, we use a truly random number (in this library, pseudorandom number generator is used). Sum it
// with the numerical equivalent of our plaintext. Numerical equivalent in this library refers to the index in a slice in which
// the alphabet is stored. After that, the sum of the plaintext and the random number is modulo'ed by the length of our alphabets.
//
// The plaintext and secret itself must conform to the alphabets defined beforehand, or the encryption process will ignore the alphabet (may cause
// errors). In this library, there is a function named `ValidateString` which will validate whether the plaintext inserted into the function
// is correct or not. There should also not be any duplicates in the character set. The default character set / alphabets is A to Z (uppercase only).
//
// This library attempts its best to conform to SRP (Single Responsibility Principle). Basically, this library expects you to understand
// the process (what to do) when performing encryption or decryption, but in most cases, all you need to do is (example for encryption):
//
// - Instantiate a new `Vigenere` struct.
// - Receive a plaintext from the user, verify if the plaintext conforms to the character set / alphabets.
// - Ensure that the secret given by the user is equal in length with the plaintext, if not, we can generate our own secrets.
// - Encrypt the data.
// - Done, and you'll get the ciphertext.
// - You are free to use the output in any form as you see fit: UTF-8, Base32, Base64, Base16 (Hex), Buffer/Bytes, or anything.
//
// To check the implementation as a console application, see example at `example/main.go`.
package vigenere

import (
	"crypto/rand"
	"errors"
	"io"
	"math/big"
	"strings"
)

// errDuplicateAlphabets is an error upon using `NewVigenere`, but with duplicate alphabets.
var errDuplicateAlphabets = errors.New("NewVigenere: Character set contains duplicate characters")

// errInvalidString is an error that will be thrown on `ValidateString` method.
var errInvalidString = errors.New("ValidateString: string does not conform to the defined vigenere's alphabets")

// errDecryptLengthNotEqual is an error that will be thrown if `Decrypt` secret key's length is less than the ciphertext.
var errDecryptLengthNotEqual = errors.New("Decrypt: ciphertext and secret are not of equal length")

// errEncryptLengthNotEqual is an error that will be thrown if `Encrypt` secret key's length is less than the plaintext.
var errEncryptLengthNotEqual = errors.New("Encrypt: plaintext and secret are not of equal length")

// Vigenere provides a struct to hold the required alphabets, its length, and its source of randomness.
type Vigenere struct {
	Alphabets    []string  // List of alphabets or characters to be used for the algorithm. The numerical equivalent of the alphabet is the slice index. Do not input duplicate characters. Data type is not `rune` for interoperability (Go internally uses UTF-8 encoding).
	Length       int64     // Length of `Alphabets` slice. Automatically created in order to not waste time and space in various methods in which this attribute is used.
	RandomSource io.Reader // Source of the random number generator. Automatically set to `rand.Reader` in order to use `crypto/rand` module for pseudorandom number generation.
}

// NewVigenere creates a new instance of `Vigenere`, along with its methods. If you desire
// to use the default alphabets, pass `nil` or `{}string[]` as the argument. This method will
// throw an error if you pass any duplicate alphabets (passing 'A' and 'A' for example).
func NewVigenere(alphabets []string) (*Vigenere, error) {
	set := make(map[string]bool, 0)

	if len(alphabets) == 0 {
		alphabets = GenerateDefaultAlphabets()
	}

	// Ensures the character set does not contain any duplicates.
	for _, alphabet := range alphabets {
		if ok := set[alphabet]; ok {
			return nil, errDuplicateAlphabets
		}

		set[alphabet] = true
	}

	return &Vigenere{
		Alphabets:    alphabets,
		Length:       int64(len(alphabets)),
		RandomSource: rand.Reader,
	}, nil
}

// ConvertToNumber converts a string into their numeric equivalent. The numeric equivalent
// is found on `Alphabets` slice.
func (v *Vigenere) ConvertToNumeric(str string) []int64 {
	numeric := make([]int64, 0)

	// This iteration costs O(N^2) or quadratic time (exhaustive search).
	for _, r := range str {
		for i, char := range v.Alphabets {
			if char == string(r) {
				numeric = append(numeric, int64(i))
			}
		}
	}

	return numeric
}

// ConvertToString converts an int64 slice to their alphabet equivalent, and
// transforms them to a single string, ready to be used. Essentially, the slice index
// is the numerical equivalent of an alphabet.
func (v *Vigenere) ConvertToString(numbers []int64) string {
	str := make([]string, 0)

	for _, num := range numbers {
		str = append(str, v.Alphabets[num])
	}

	return strings.Join(str, "")
}

// Decrypt decrypts a ciphertext with a secret key. Make sure that the secret key is equal in length with the
// ciphertext. Returns the plaintext.
func (v *Vigenere) Decrypt(ciphertext, secret string) (string, error) {
	// Prepare a slice to hold the plaintext numerical representative.
	plaintext := make([]int64, 0)

	// Ensures that the length are at least equal.
	if len(secret) < len(ciphertext) {
		return "", errDecryptLengthNotEqual
	}

	// Change the ciphertext and the secrets to their numeric equivalent.
	numericCiphertext := v.ConvertToNumeric(ciphertext)
	numericSecret := v.ConvertToNumeric(secret)

	// Subtract the numeric equivalent with the equivalent numeric secret. If the result
	// is lower than zero, increment by the length of the `Alphabets` slice. We have no need
	// for modulo as doing it would only return the same value (no effect).
	for i, char := range numericCiphertext {
		equivalent := char - numericSecret[i]
		if equivalent < 0 {
			equivalent = equivalent + v.Length
		}

		plaintext = append(plaintext, equivalent)
	}

	// Return plaintext.
	return v.ConvertToString(plaintext), nil
}

// Encrypt encrypts a plaintext with a secret key. Make sure that the secret key is equal in length with the
// plaintext. Returns the ciphertext.
func (v *Vigenere) Encrypt(plaintext, secret string) (string, error) {
	// Prepare a slice to hold the ciphertext numerical representative.
	ciphertext := make([]int64, 0)

	// Ensures that the length are at least equal.
	if len(secret) < len(plaintext) {
		return "", errEncryptLengthNotEqual
	}

	// Change the plaintext and secret to their numeric equivalent.
	numericPlaintext := v.ConvertToNumeric(plaintext)
	numericSecret := v.ConvertToNumeric(secret)

	// Sum the numeric equivalent with the numeric equivalent of the secret, then modulo
	// it with the length of the `Alphabets` slice.
	for i, char := range numericPlaintext {
		equivalent := (numericSecret[i] + char) % v.Length
		ciphertext = append(ciphertext, equivalent)
	}

	// Return ciphertext.
	return v.ConvertToString(ciphertext), nil
}

// GenerateDefaultAlphabets generates alphabets from 'A' to 'Z', or characters with code
// 65 to 90 in ASCII format. The `i` variable is transformed from `rune` to `string` for
// compatibility.
func GenerateDefaultAlphabets() []string {
	alphabets := make([]string, 0)

	for i := 65; i <= 90; i++ {
		alphabets = append(alphabets, string(rune(i)))
	}

	return alphabets
}

// GenerateRandomNumber securely generates a random number (index) from zero to the max length of
// the `Alphabets` slice. Function `rand.Int` is exclusive (from 0 to n - 1) where `n` can be of any length.
func (v *Vigenere) GenerateRandomNumber() (int64, error) {
	number, err := rand.Int(v.RandomSource, big.NewInt(v.Length))
	if err != nil {
		return 0, err
	}

	return number.Int64(), nil
}

// GenerateSecretKey generates a secure secret that is equal in length with the plaintext. This is done
// to ensure the security of the encryption. Complexity is O(N).
func (v *Vigenere) GenerateSecretKey(plaintext string) (string, error) {
	secret := make([]int64, 0)

	for range plaintext {
		number, err := v.GenerateRandomNumber()
		if err != nil {
			return "", err
		}

		secret = append(secret, number)
	}

	return v.ConvertToString(secret), nil
}

// ValidateString validates a string whether it conforms to the required alphabets or not. The process
// transforms the `Alphabets` slice into a hash map / object / map, and then checks the availability of
// each letters.
func (v *Vigenere) ValidateString(str string) (bool, error) {
	set := make(map[string]bool)

	// Building this map/object/hash table costs O(N) time.
	for _, alphabet := range v.Alphabets {
		set[alphabet] = true
	}

	// Checks whether a letter from `str` exists in the new set. Costs O(1) for retrieval, and
	// O(N) to iterate the whole `str` variable.
	for _, r := range str {
		if ok := set[string(r)]; !ok {
			return false, errInvalidString
		}
	}

	return true, nil
}
