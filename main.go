package main

import (
        "fmt"
        "time"
        "log"
        "strings"
        "github.com/metrico/pasticca/paste"
)

func main() {
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

               // Example: Save encrypted content
        encryptedContent := "This is a secret message."
        encryptedFingerprint, encryptedHashWithAnchor, err := paste.Save(encryptedContent, "", "", true)
        if err != nil {
                log.Fatalf("Error saving encrypted paste: %v", err)
        }
        fmt.Printf("Saved encrypted paste with fingerprint/hash: %s/%s\n", encryptedFingerprint, encryptedHashWithAnchor)
        time.Sleep(2 * time.Second)

        // Example: Load and decrypt the encrypted content
        decryptedContent, isStillEncrypted, err := paste.Load(encryptedFingerprint, encryptedHashWithAnchor)
        if err != nil {
                log.Fatalf("Error loading encrypted paste: %v", err)
        }
        fmt.Printf("Loaded and decrypted content: %s\nIs still encrypted: %v\n", decryptedContent, isStillEncrypted)

        // Example: Try to load encrypted content without the key
        encryptedHashWithoutAnchor := strings.Split(encryptedHashWithAnchor, "#")[0]
        encryptedContentWithoutKey, isEncryptedWithoutKey, err := paste.Load(encryptedFingerprint, encryptedHashWithoutAnchor)
        if err != nil {
                log.Fatalf("Error loading encrypted paste without key: %v", err)
        }
        fmt.Printf("Loaded encrypted content without key: %s\nIs encrypted: %v\n", encryptedContentWithoutKey, isEncryptedWithoutKey)
}
