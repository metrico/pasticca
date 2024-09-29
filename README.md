# 💊 pasticca
Store and Retrieve data from `pastila` in Go for fun & profit-loss

## Usage
```go
import ("github.com/metrico/pasticca/paste")
```

## Example
See `main.go`

#### Plaintext/JSON
```go
// Example: Save some content
content := "This is a test paste."
fingerprint, hashWithAnchor, err := paste.Save(content, "", "", false)
if err != nil {
        log.Fatalf("Error saving paste: %v", err)
}
hash := strings.Split(hashWithAnchor, "#")[0] // Remove anchor if present
fmt.Printf("Saved paste with fingerprint/hash: %s/%s\n", fingerprint, hash)

// Example: Load the content we just saved
loadedContent, isEncrypted, err := paste.Load(fingerprint, hash)
if err != nil {
        log.Fatalf("Error loading paste: %v", err)
}
fmt.Printf("Loaded content: %s\nIs encrypted: %v\n", loadedContent, isEncrypted)
```
```
Saved paste with fingerprint/hash: 913ae2b1/748ab86a806c2de1fd5753fb3ffff516
Loaded content: This is a test paste.
Is encrypted: false
```

#### Encrypted
```go
// Example: Save encrypted content
encryptedContent := "This is a secret message."
encryptedFingerprint, encryptedHashWithAnchor, err := paste.Save(encryptedContent, "", "", true)
if err != nil {
        log.Fatalf("Error saving encrypted paste: %v", err)
}
fmt.Printf("Saved encrypted paste with fingerprint/hash: %s/%s\n", encryptedFingerprint, encryptedHashWithAnchor)

// Example: Load and decrypt the encrypted content
decryptedContent, isStillEncrypted, err := paste.Load(encryptedFingerprint, encryptedHashWithAnchor)
if err != nil {
log.Fatalf("Error loading encrypted paste: %v", err)
}
fmt.Printf("Loaded and decrypted content: %s\nIs encrypted: %v\n", decryptedContent, isStillEncrypted)
```
```
Saved encrypted paste with fingerprint/hash: b4765c53/94cf5b7bee267b1d41c9ada746ebe6e1#FFUgNmg29LqBLdN3LQdfzw==
Loaded and decrypted content: This is a secret message.
Is encrypted: true
```
