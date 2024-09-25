package main

import (
        "fmt"
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

        /*
        // Example: Save encrypted content
        encryptedContent := "This is a secret message."
        encryptedFingerprint, encryptedHashWithAnchor, err := paste.Save(encryptedContent, "", "", true)
        if err != nil {
                log.Fatalf("Error saving encrypted paste: %v", err)
        }
        encryptedHash := strings.Split(encryptedHashWithAnchor, "#")[0] // Remove anchor
        fmt.Printf("Saved encrypted paste with fingerprint/hash: %s/%s\n", encryptedFingerprint, encryptedHash)

        // Note: To decrypt the content, you would need to handle the key from the URL anchor
        // and pass it to the AESDecrypt function. This example doesn't show that process.
        */
}
