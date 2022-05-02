package main

import (
	"fmt"
	"log"

	"github.com/lauslim12/vigenere"
)

// Driver code.
func main() {
	// Creates a new Vigenere instance.
	vigenere, err := vigenere.NewVigenere([]string{})
	if err != nil {
		log.Fatal(err.Error())
	}

	// Define our plaintext.
	plaintext := "VIGENERECIPHER"

	// Ensures that our plaintext string is valid.
	valid, err := vigenere.ValidateString(plaintext)
	if err != nil {
		log.Fatal(err.Error())
	}

	// Generates a secret key as long as the plaintext.
	secret, err := vigenere.GenerateSecretKey(plaintext)
	if err != nil {
		panic(fmt.Sprintf("crypto/rand failed, reason: %v", err.Error()))
	}

	// Encrypts our plaintext string with our secret key.
	ciphertext, err := vigenere.Encrypt(plaintext, secret)
	if err != nil {
		log.Fatal(err.Error())
	}

	// Decrypts our ciphertext with our secret.
	decrypted, err := vigenere.Decrypt(ciphertext, secret)
	if err != nil {
		log.Fatal(err.Error())
	}

	// Print all results.
	fmt.Println("Vigenere Encryption done with the following results:")
	fmt.Println("Valid Plaintext:", valid)
	fmt.Println("Plaintext:", plaintext)
	fmt.Println("Secret Key:", secret)
	fmt.Println("Ciphertext:", ciphertext)
	fmt.Println("Decrypted:", decrypted)
}
