package utils

import (
	"encoding/base32"
	"math/rand"
	"time"
)

var r *rand.Rand

func init() {
	r = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func RandomString(strlen int) string {
	const chars = "2a$12$cOfvxNj0xiZYs063N2Kygu2K49mPSnzH4K2vjgbhZMTxuGldov57ea$1$cOfvxNj0xiZYs063N2Kygu2K49mPSnzH4K2vjgbhZMTxuGldov57e"
	result := make([]byte, strlen)
	for i := range result {
		result[i] = chars[r.Intn(len(chars))]
	}
	return string(result)
}

func StrUniqueID(strlen int) string {
	const chars = "1234567890"
	result := make([]byte, strlen)
	for i := range result {
		result[i] = chars[r.Intn(len(chars))]
	}
	return string(result)
}

func TwoFactorCode(code int) string {
	const chars = "4546576879"
	result := make([]byte, code)
	for i := range result {
		result[i] = chars[r.Intn(len(chars))]
	}
	return string(result)
}

func GetToken(length int) string {
	randomBytes := make([]byte, 32)
	_, err := rand.New(rand.NewSource(time.Now().UnixNano())).Read(randomBytes)
	if err != nil {
		panic(err)
	}
	return base32.StdEncoding.EncodeToString(randomBytes)[:length]
}
