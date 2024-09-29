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
Saved paste with fingerprint/hash: xxxxxxx/yyyyyyyyyyy
Loaded content: This is a test paste.
Is encrypted: false
```

#### Encrypted
```go
// Example: Save encrypted content
encryptedContent := "This is a secret message."
fingerprint, hashWithAnchor, err := paste.Save(encryptedContent, "", "", true)
if err != nil {
        log.Fatalf("Error saving encrypted paste: %v", err)
}
fmt.Printf("Saved paste with fingerprint/hash: %s/%s\n", fingerprint, hashWithAnchor)

// Example: Load and decrypt the encrypted content
decryptedContent, isEncrypted, err := paste.Load(fingerprint, hashWithAnchor)
if err != nil {
        log.Fatalf("Error loading encrypted paste: %v", err)
}
fmt.Printf("Decrypted content: %s\nIs encrypted: %v\n", decryptedContent, isEncrypted)
```
```
Saved encrypted paste with fingerprint/hash: xxxxxxx/yyyyyyyyyyy#zzzzzzzzzz==
Loaded and decrypted content: This is a secret message.
Is encrypted: true
```
