package users

import (
    "crypto/sha256"
    "math/rand"
    "time"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func Hash(s string) []byte {
    h := sha256.New()
    h.Write([]byte(s))
    
    bs := h.Sum(nil)

    return bs
}

func GenerateRandomString(length int) string {
    var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
    b := make([]byte, length)
    for i := range b {
        b[i] = charset[seededRand.Intn(len(charset))]
    }
    return string(b)
}
