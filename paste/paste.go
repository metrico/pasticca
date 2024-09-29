package paste

import (
        "crypto/aes"
        "crypto/cipher"
        "crypto/rand"
        "encoding/base64"
        "encoding/binary"
        "encoding/hex"
        "encoding/json"
        "errors"
        "io"
        "net/http"
        "regexp"
        "strings"
        "fmt"
//        "io/ioutil"
)

const clickhouseURL = "https://play.clickhouse.com/?user=paste"

// DataResponse represents the structure of the JSON response from Clickhouse
type DataResponse struct {
        Data []struct {
                Content     string `json:"content"`
                IsEncrypted uint8  `json:"is_encrypted"` // Change this to uint8
        } `json:"data"`
        Rows int `json:"rows"`
}

// Load retrieves data from Clickhouse
func Load(fingerprint, hashWithAnchor string) (string, bool, error) {
        parts := strings.SplitN(hashWithAnchor, "#", 2)
        hash := parts[0]
        var key string
        if len(parts) > 1 {
                key = parts[1]
        }

        query := fmt.Sprintf(`
                SELECT content, is_encrypted
                FROM data
                WHERE fingerprint = reinterpretAsUInt32(unhex('%s'))
                  AND hash = reinterpretAsUInt128(unhex('%s'))
                ORDER BY time DESC
                LIMIT 1
                FORMAT JSON
        `, fingerprint, hash)

        resp, err := http.Post(clickhouseURL, "application/x-www-form-urlencoded", strings.NewReader(query))
        if err != nil {
                return "", false, fmt.Errorf("HTTP request failed: %v", err)
        }
        defer resp.Body.Close()

        if resp.StatusCode != http.StatusOK {
                return "", false, fmt.Errorf("HTTP status %s", resp.Status)
        }

        body, err := io.ReadAll(resp.Body)
        if err != nil {
                return "", false, fmt.Errorf("failed to read response body: %v", err)
        }

        var response DataResponse
        err = json.Unmarshal(body, &response)
        if err != nil {
                return "", false, fmt.Errorf("failed to unmarshal JSON: %v", err)
        }

        if response.Rows < 1 {
                return "", false, fmt.Errorf("paste not found or multiple rows returned (rows: %d)", response.Rows)
        }

        content := response.Data[0].Content
        isEncrypted := response.Data[0].IsEncrypted != 0

        if isEncrypted && key != "" {
                decryptedContent, err := DecryptContent(content, key)
                if err != nil {
                        return "", true, fmt.Errorf("failed to decrypt content: %v", err)
                }
                content = decryptedContent
        }

        return content, isEncrypted, nil
}

// DecryptContent decrypts the content using the provided key
func DecryptContent(encryptedContent, keyBase64 string) (string, error) {
        key, err := base64.StdEncoding.DecodeString(keyBase64)
        if err != nil {
                return "", fmt.Errorf("failed to decode key: %v", err)
        }

        ciphertext, err := base64.StdEncoding.DecodeString(encryptedContent)
        if err != nil {
                return "", fmt.Errorf("failed to decode ciphertext: %v", err)
        }

        block, err := aes.NewCipher(key)
        if err != nil {
                return "", fmt.Errorf("failed to create cipher: %v", err)
        }

        if len(ciphertext) < aes.BlockSize {
                return "", errors.New("ciphertext too short")
        }

        iv := ciphertext[:aes.BlockSize]
        ciphertext = ciphertext[aes.BlockSize:]

        stream := cipher.NewCTR(block, iv)
        plaintext := make([]byte, len(ciphertext))
        stream.XORKeyStream(plaintext, ciphertext)

        return string(plaintext), nil
}


// Save stores data in Clickhouse
func Save(content, prevFingerprint, prevHash string, isEncrypted bool) (string, string, error) {
        text := content
        var anchor string

        if isEncrypted {
                encryptedText, key, err := aesEncrypt([]byte(text))
                if err != nil {
                        return "", "", err
                }
                text = encryptedText
                anchor = "#" + base64.StdEncoding.EncodeToString(key)
        }

        currHash := sipHash128([]byte(text))
        currFingerprint := getFingerprint(text)

        data := struct {
                FingerprintHex    string `json:"fingerprint_hex"`
                HashHex           string `json:"hash_hex"`
                PrevFingerprintHex string `json:"prev_fingerprint_hex"`
                PrevHashHex       string `json:"prev_hash_hex"`
                Content           string `json:"content"`
                IsEncrypted       bool   `json:"is_encrypted"`
        }{
                FingerprintHex:    currFingerprint,
                HashHex:           currHash,
                PrevFingerprintHex: prevFingerprint,
                PrevHashHex:       prevHash,
                Content:           text,
                IsEncrypted:       isEncrypted,
        }

        jsonData, err := json.Marshal(data)
        if err != nil {
                return "", "", err
        }

        query := "INSERT INTO data (fingerprint_hex, hash_hex, prev_fingerprint_hex, prev_hash_hex, content, is_encrypted) FORMAT JSONEachRow " + string(jsonData)

        resp, err := http.Post(clickhouseURL, "application/x-www-form-urlencoded", strings.NewReader(query))
        if err != nil {
                return "", "", err
        }
        defer resp.Body.Close()

        if resp.StatusCode != http.StatusOK {
                return "", "", errors.New("HTTP status " + resp.Status)
        }

        return currFingerprint, currHash + anchor, nil
}

// AESEncrypt encrypts content using AES-CTR
func aesEncrypt(plaintext []byte) (string, []byte, error) {
        key := make([]byte, 16)
        if _, err := io.ReadFull(rand.Reader, key); err != nil {
                return "", nil, err
        }

        block, err := aes.NewCipher(key)
        if err != nil {
                return "", nil, err
        }

        ciphertext := make([]byte, aes.BlockSize+len(plaintext))
        iv := ciphertext[:aes.BlockSize]
        if _, err := io.ReadFull(rand.Reader, iv); err != nil {
                return "", nil, err
        }

        stream := cipher.NewCTR(block, iv)
        stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

        return base64.StdEncoding.EncodeToString(ciphertext), key, nil
}

// AESDecrypt decrypts content using AES-CTR
func AESDecrypt(ciphertext string, key []byte) (string, error) {
        ciphertextBytes, err := base64.StdEncoding.DecodeString(ciphertext)
        if err != nil {
                return "", err
        }

        block, err := aes.NewCipher(key)
        if err != nil {
                return "", err
        }

        if len(ciphertextBytes) < aes.BlockSize {
                return "", errors.New("ciphertext too short")
        }

        iv := ciphertextBytes[:aes.BlockSize]
        ciphertextBytes = ciphertextBytes[aes.BlockSize:]

        stream := cipher.NewCTR(block, iv)
        stream.XORKeyStream(ciphertextBytes, ciphertextBytes)

        return string(ciphertextBytes), nil
}

// GetFingerprint generates a fingerprint from text
func getFingerprint(text string) string {
        re := regexp.MustCompile(`\p{L}{4,100}`)
        matches := re.FindAllString(text, -1)
        if len(matches) == 0 {
                return "ffffffff"
        }

        minHash := "ffffffff"
        seen := make(map[string]bool)

        for i := 0; i < len(matches)-2; i++ {
                elem := strings.Join(matches[i:i+3], ",")
                if !seen[elem] {
                        seen[elem] = true
                        hash := sipHash128([]byte(elem))[:8]
                        if hash < minHash {
                                minHash = hash
                        }
                }
        }

        return minHash
}

// SipHash128 implements the SIPHash-2-4 function
func sipHash128(message []byte) string {
        v0, v1, v2, v3 := uint64(0x736f6d6570736575), uint64(0x646f72616e646f6d), uint64(0x6c7967656e657261), uint64(0x7465646279746573)

        compress := func() {
                v0 += v1
                v1 = (v1 << 13) | (v1 >> 51)
                v1 ^= v0
                v0 = (v0 << 32) | (v0 >> 32)
                v2 += v3
                v3 = (v3 << 16) | (v3 >> 48)
                v3 ^= v2
                v0 += v3
                v3 = (v3 << 21) | (v3 >> 43)
                v3 ^= v0
                v2 += v1
                v1 = (v1 << 17) | (v1 >> 47)
                v1 ^= v2
                v2 = (v2 << 32) | (v2 >> 32)
        }

        getBlock := func(b []byte) uint64 {
                return binary.LittleEndian.Uint64(b)
        }

        messageLen := len(message)
        for len(message) >= 8 {
                m := getBlock(message)
                v3 ^= m
                compress()
                compress()
                v0 ^= m
                message = message[8:]
        }

        var lastBlock [8]byte
        copy(lastBlock[:], message)
        lastBlock[7] = byte(messageLen)

        m := getBlock(lastBlock[:])
        v3 ^= m
        compress()
        compress()
        v0 ^= m
        v2 ^= 0xff
        compress()
        compress()
        compress()
        compress()

        b := make([]byte, 16)
        binary.LittleEndian.PutUint64(b, v0^v1)
        binary.LittleEndian.PutUint64(b[8:], v2^v3)

        return hex.EncodeToString(b)
}
