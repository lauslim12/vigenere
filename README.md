# Vigenère

Provides a dependency-free, secure symmetric encryption implementation with Vigenère-Vernam Cipher algorithm as a library for Go. Acts as a One-Time Pad (OTP) to ensure maximum security, and supports custom alphabets.

## About

Vigenère-Vernam Cipher algorithm is theoretically a 100% secure encryption algorithm, as long as you follow these (impossible) three points:

- The key must be as long as the plaintext.
- The key must be generated with true random number generator (TRNG).
- The key must be secret and never reused.

In Computer Security, there are three pillars: Confidentiality (ensures only authorized personnel will read the data), Integrity (ensures data
is not tampered), and Availability (ensures the data can be accessed anytime). Confidentiality is assured with this cipher, but the
integrity and availability are not 100% covered. Because of this, technically, One-Time Pads are usually not feasible to implement in the
current era. Several reasons will answer why:

- Storage: Storing a long pad with an equivalent length secret requires memory. Imagine if there are 5 million people, how much will the memory cost be?
- True randomness: Computers are deterministic, and they have no idea on how to create random numbers (unlike humans).
- Authentication: One-Time Pads are stream ciphers and do not provide Message Authentication Codes (MAC). A malicious party can bit-flip the ciphertext.
- Transport: One-Time Pads are as secure as the secure channel used to deliver the keys. Internet is not secure enough.

Trying to implement One-Time Pads in the modern era raises a single question:

- If you can create a completely secure channel that is unable to be hacked by people, why don't you use that channel instead?

## Algorithm

The algorithm for encryption:

- Define a character set / alphabets to be used and a plaintext to be encrypted. In this case, let's assume that our character set is A to Z, with A is represented by 0, B is represented by 1, C is represented by 2, and so on and so forth.
- Generate a set of random numbers whose length is as long as the plaintext.
- Transform the plaintext into a numerical representative of your character set. For example, `ABC` with above character set is represented as `012`.
- Sum the plaintext numerical representative and the random numbers.
- Modulo the result by the length of your character set.
- Transform the result to the numerical representative. This is your ciphertext.

The algorithm for decryption:

- Transform the ciphertext and the secret in the numerical representative form.
- Subtract the ciphertext with the secret. If any letter goes below zero, add the length of your character set.
- Transform the result as the numerical representative. This is your plaintext.

To showcase the whole algorithm:

```bash
# Encryption:

      h       e       l       l       o  message
   7 (h)   4 (e)  11 (l)  11 (l)  14 (o) message
+ 23 (X)  12 (M)   2 (C)  10 (K)  11 (L) key
= 30      16      13      21      25     message + key
=  4 (E)  16 (Q)  13 (N)  21 (V)  25 (Z) (message + key) mod 26
      E       Q       N       V       Z  → ciphertext

# Decryption:

       E       Q       N       V       Z  ciphertext
    4 (E)  16 (Q)  13 (N)  21 (V)  25 (Z) ciphertext
−  23 (X)  12 (M)   2 (C)  10 (K)  11 (L) key
= −19       4      11      11      14     ciphertext – key
=   7 (h)   4 (e)  11 (l)  11 (l)  14 (o) ciphertext – key + length_alphabets (if less than zero)
       h       e       l       l       o  → message
```

## Features

General:

- **No dependencies.** Only needs standard Go and no dependencies are required.
- **Battle-tested.** This library conforms to the standard library.
- **Lightweight.** Small in size due to not having any dependencies.
- **Secure.** Tries its best to implement as many security considerations as possible with careful and secure coding practices.
- **Single Responsibility Principle (SRP)**. Every function in this library will do only one thing only.
- **100% tested.** As this library is small, the code coverage is still 100% for now.
- **Well documented.** Check out this `README.md` document and the technical documentation for further reading!

Specific:

- **Supports custom alphabets.** By default, Vigenère-Vernam Cipher only supports Latin alphabets A to Z, all in uppercase. This library allows you to use custom alphabets by defining the character set. You can use numbers, lowercase characters, and even symbols. Because we are using slices, these custom alphabets are implemented without reducing the overall security of the library (without inadvertently creating modulo bias, error in calculations, and others).
- **Strong random secret.** We use Golang's `crypto/rand`'s `rand.Int`, which takes the numbers from the hardware noise (`/dev/urandom` in Linux). This is by no means a perfect, true random number generator, but it can at least be cryptographically secure. The secret also has to be at least the length of the plaintext to prevent weak keys, as our goal with this library is also to make it as a One-Time Pad. Random secret in `*Vigenere` implements the `io.Reader` interface. You may replace it if needs arise.

## Documentation

Complete documentation could be seen in the official [pkg.go.dev site](https://pkg.go.dev/github.com/lauslim12/vigenere).

## Installation

I am going to assume that you are using Go version 1.18.

- Download this library.

```bash
go install github.com/lauslim12/vigenere

# or: go get -u github.com/lauslim12/vigenere
```

- Import in your source code.

```go
import "github.com/lauslim12/vigenere"
```

- Instantiate the `Vigenere` structure in your code, and define your plaintext.

```go
vigenere, err := vigenere.NewVigenere(nil) // Use default alphabets.
if err != nil {
    log.Fatal(err.Error())
}

plaintext := "PLAINTEXT"
```

- The steps are to validate your string, generate your own secret with pseudorandom number generator, encrypt your string, then convert it to the number equivalent. The steps are made like this so there would be no invalid inputs, and so the library would conform to Single Responbility Principle for each of the functions.

```go
valid, err := vigenere.ValidateString(plaintext)
if err != nil {
    log.Fatal(err.Error())
}

key, err := vigenere.GenerateSecretKey(plaintext)
if err != nil {
    log.Fatal(err.Error())
}

ciphertext, err := vigenere.Encrypt(plaintext, key)
if err != nil {
    log.Fatal(err.Error())
}

fmt.Println(ciphertext)
```

- Done! If you want, feel free to look at tests (`vigenere_test.go`, run by using `go test -cover ./... ./...`) and the source code itself.

## Notes

Several notes if you want to understand the engineering decisions that I have made during the creation of this library:

- **What happens if one step, let's say the string validation is skipped?**

The program may panic (runtime error, usually the error is slice index out of range, that happens if a text contains a rune/letter that does not exist in the character set or alphabets) if the input does not conform to the parameters. That's why string validation is important. Please use the provided `vigenere.ValidateString(str string)` to validate your string before inputting them to any of the methods.

- **Wouldn't it be better to use `rune` and ASCII characters instead of using slices for the alphabets?**

I have thought about using ASCII characters. It can be done by subtracting with `65` to get the numerical representative for the `A` character. I also have no need to create slices if I use ASCII characters. However, it comes at a cost, that is: I cannot use custom alphabets for the cipher. If I use `string` and slices, using custom alphabets is more than possible.

- **Why the time and space complexity can be so high? O(N) and O(N^2) seems a bit over the top.**

The time complexity is a bit high because of the utilization of slices. I prefer to use slices over `map[string]int64` as slices are simpler and I believe it is enough for this use-case. I don't think someone wants to encrypt a string with over 1.000.000 runes / letters as its character set. If I were to use the ASCII style as above, I would be able to reach O(1) efficiency, but we would not be able to use custom alphabets.

## Examples

Please see examples at [the example project(`example/main.go`)](./example). You may run it by using `go run example/main.go` and the results will be shown instantly.

## Contributing

This tool is open source and the contribution of this tool is highly encouraged! If you want to contribute to this project, please feel free to read the `CONTRIBUTING.md` file for the contributing guidelines.

## License

This work is licensed under MIT License. Please check the `LICENSE` file for more information.

## References

Here are the references that I have used to write this library:

- [How is One-Time Pad Perfectly Secure](https://crypto.stackexchange.com/questions/31084/how-is-the-one-time-pad-otp-perfectly-secure)
- [Is Modern Encryption Needlessly Complicated?](https://crypto.stackexchange.com/questions/596/is-modern-encryption-needlessly-complicated)
- [One-Time Pad Example](https://www.boxentriq.com/code-breaking/one-time-pad)
- [One-Time Pad Example Calculation](https://en.wikipedia.org/wiki/One-time_pad)
- [Why One-Time Pad is Useless in Practice](https://crypto.stackexchange.com/questions/15652/one-time-pad-why-is-it-useless-in-practice)
- [Why is One-Time Pad not Vulerable to Brute-Force Attacks?](https://crypto.stackexchange.com/questions/596/is-modern-encryption-needlessly-complicated)
