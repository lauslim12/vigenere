package vigenere

import (
	"errors"
	"log"
	"testing"
	"testing/iotest"
)

func TestDecrypt(t *testing.T) {
	vigenere, err := NewVigenere(nil)
	if err != nil {
		log.Fatal(err.Error())
	}

	vigenereCustom, err := NewVigenere([]string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O"})
	if err != nil {
		log.Fatal(err.Error())
	}

	tests := []struct {
		name           string
		ciphertext     string
		secret         string
		expectedOutput string
		expectedError  bool
		vigenere       *Vigenere
	}{
		{
			name:           "test_success_decrypt_1",
			ciphertext:     "VMKQTRKAUSGHQ",
			secret:         "TERSFMBAKSPOQ",
			expectedOutput: "CITYOFJAKARTA",
			expectedError:  false,
			vigenere:       vigenere,
		},
		{
			name:           "test_success_decrypt_2",
			ciphertext:     "WPFEC",
			secret:         "PLUTO",
			expectedOutput: "HELLO",
			expectedError:  false,
			vigenere:       vigenere,
		},
		{
			name:           "test_success_decrypt_3",
			ciphertext:     "CCMB",
			secret:         "ADNF",
			expectedOutput: "COOL",
			expectedError:  false,
			vigenere:       vigenereCustom,
		},
		{
			name:           "test_failure_decrypt_1",
			ciphertext:     "WPFTY",
			secret:         "A",
			expectedOutput: "",
			expectedError:  true,
			vigenere:       vigenere,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := tc.vigenere.Decrypt(tc.ciphertext, tc.secret)
			if err != nil && !tc.expectedError {
				t.Errorf("Encrypt method should not result in an error. Got: %v.", err.Error())
			}

			if tc.expectedOutput != result {
				t.Errorf("Expected and actual output are different! Expected: %v. Got: %v.", tc.expectedOutput, result)
			}
		})
	}
}

func TestEncrypt(t *testing.T) {
	vigenere, err := NewVigenere(nil)
	if err != nil {
		log.Fatal(err.Error())
	}

	vigenereCustom, err := NewVigenere([]string{"A", "B", "C", "D", "E", "F", "G", "H", "I"})
	if err != nil {
		log.Fatal(err.Error())
	}

	tests := []struct {
		name           string
		plaintext      string
		secret         string
		expectedOutput string
		expectedError  bool
		vigenere       *Vigenere
	}{
		{
			name:           "test_success_encrypt_1",
			plaintext:      "VIGENERECIPHER",
			secret:         "SDTLHTHMTMQWER",
			expectedOutput: "NLZPUXYQVUFDII",
			expectedError:  false,
			vigenere:       vigenere,
		},
		{
			name:           "test_success_encrypt_2",
			plaintext:      "VIGENEREVERNAMCIPHER",
			secret:         "LOPWRTSKVMSKLFKRYJTQ",
			expectedOutput: "GWVAEXJOQQJXLRMZNQXH",
			expectedError:  false,
			vigenere:       vigenere,
		},
		{
			name:           "test_success_encrypt_3",
			plaintext:      "HI",
			secret:         "EF",
			expectedOutput: "CE",
			expectedError:  false,
			vigenere:       vigenereCustom,
		},
		{
			name:           "test_failure_encrypt",
			plaintext:      "VIGENERECIPHER",
			secret:         "A",
			expectedOutput: "",
			expectedError:  true,
			vigenere:       vigenere,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := tc.vigenere.Encrypt(tc.plaintext, tc.secret)
			if err != nil && !tc.expectedError {
				t.Errorf("Encrypt method should not result in an error. Got: %v.", err.Error())
			}

			if tc.expectedOutput != result {
				t.Errorf("Expected and actual output are different! Expected: %v. Got: %v.", tc.expectedOutput, result)
			}
		})
	}
}

func TestGenerateSecret(t *testing.T) {
	vigenere, err := NewVigenere(nil)
	if err != nil {
		log.Fatal(err.Error())
	}

	tests := []struct {
		name                 string
		plaintext            string
		expectedOutputLength int
		expectedError        bool
		vigenere             *Vigenere
	}{
		{
			name:                 "test_generate_secret_success",
			plaintext:            "HELLO",
			expectedOutputLength: 5,
			expectedError:        false,
			vigenere:             vigenere,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			secret, err := tc.vigenere.GenerateSecretKey(tc.plaintext)
			if err != nil {
				log.Fatal(err.Error())
			}

			if len(secret) != tc.expectedOutputLength {
				t.Errorf("Expected and actual output are different! Expected: %v. Got: %v.", tc.expectedOutputLength, secret)
			}
		})
	}
}

func TestNewVigenere(t *testing.T) {
	t.Run("test_vigenere_creation", func(t *testing.T) {
		_, err := NewVigenere(nil)
		if err != nil {
			t.Error("Call of 'NewVigenere' should not return an error.")
		}

		// Yes, this was inversed. If the error is not nil.
		_, err = NewVigenere([]string{"A", "A", "B"})
		if err == nil {
			t.Error("Call of 'NewVigenere' should return an error with duplicate characters as the alphabets.")
		}
	})
}

func TestRandomNumberGenerationFailure(t *testing.T) {
	vigenere := &Vigenere{
		Alphabets:    GenerateDefaultAlphabets(),
		Length:       int64(len(GenerateDefaultAlphabets())),
		RandomSource: iotest.ErrReader(errors.New("Mocked crypto/rand error!")),
	}

	t.Run("test_random_number_generation_failure", func(t *testing.T) {
		key, err := vigenere.GenerateSecretKey("SHIBAMIYUKI")
		if err == nil {
			t.Errorf("Test RNG failure should return an error! Received: %v.", key)
		}
	})
}

func TestValidateString(t *testing.T) {
	vigenere, err := NewVigenere(nil)
	if err != nil {
		log.Fatal(err.Error())
	}

	tests := []struct {
		name           string
		text           string
		expectedResult bool
		vigenere       *Vigenere
	}{
		{
			name:           "test_success_1",
			text:           "HELLO",
			expectedResult: true,
			vigenere:       vigenere,
		},
		{
			name:           "test_success_2",
			text:           "ABC",
			expectedResult: true,
			vigenere:       vigenere,
		},
		{
			name:           "test_failure_1",
			text:           "123",
			expectedResult: false,
			vigenere:       vigenere,
		},
		{
			name:           "test_failure_2",
			text:           ".,;'[]",
			expectedResult: false,
			vigenere:       vigenere,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			valid, _ := tc.vigenere.ValidateString(tc.text) // Skip the error as we know it will definitely be correct.

			if valid != tc.expectedResult {
				t.Errorf("Expected and actual output are different! Expected: %v. Got: %v.", tc.expectedResult, valid)
			}
		})
	}
}
