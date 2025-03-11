package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"time"
)

var MAIL_DOMAIN = os.Getenv("MAIL_DOMAIN")
var SMTP_SERVER = os.Getenv("SMTP_SERVER")
var SENDER_MAIL = os.Getenv("SENDER_MAIL")

func GetDate() string {
	t := time.Now()
	return t.Format(time.RFC1123Z)
}

func GetMessageId(domain string) string {
	uuid, err := GetUUID()
	if err != nil {
		panic(err)
		return ""
	}

	messageID := fmt.Sprintf("<%s@%s>", uuid, domain)
	return messageID
}

func GetUUID() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	timeNow := time.Now().UnixNano()
	timeBytes := make([]byte, 8)
	for i := 0; i < 8; i++ {
		timeBytes[i] = byte(timeNow >> (i * 8))
	}
	b = append(b, timeBytes...)

	uuid := hex.EncodeToString(b)
	return uuid, nil
}
