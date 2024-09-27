# 💊 pasticca
Store and Retrieve data from `pastila` in Go for fun a profit-loss

## Usage
```go
import ("github.com/metrico/pasticca/paste")
```
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

## Example
See `main.go`

```
Saved paste with fingerprint/hash: 913ae2b1/748ab86a806c2de1fd5753fb3ffff516
Loaded content: This is a test paste.
Is encrypted: false
```
