package helper

import (
	"crypto/rand"
	"fmt"
	"regexp"
	"strings"
	"time"
)

func Slugify(s string) string {
	s = strings.ToLower(s)

	re := regexp.MustCompile(`[^a-z0-9]+`)
	s = re.ReplaceAllString(s, "-")

	s = strings.Trim(s, "-")

	return s
}

func GenerateInvoiceCode() string {
	currentTime := time.Now().Format("20060102")
	
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const length = 6
	randomBytes := make([]byte, length)
	rand.Read(randomBytes)
	for i := 0; i < length; i++ {
		randomBytes[i] = charset[int(randomBytes[i])%len(charset)]
	}
	randomStr := string(randomBytes)
	
	return fmt.Sprintf("INV-%s-%s", currentTime, randomStr)
}