package util

import (
	"fmt"
	"math/rand"
	"strings"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"
var currencies = [...]string{BRL, CAD, EUR, JPY, USD}
var pairs = [...]string{USDT_BRL, USDT_CAD, USDT_EUR, USDT_JPY, USDT_USD, BTC_USDT, ETH_USDT, MATIC_USDT, SOL_USDT, ETH_BTC, MATIC_BTC, SOL_BTC, MATIC_ETH, SOL_ETH}
var status = [...]string{ACTIVE, COMPLETED, CANCELED}

// RandomInt generates a random integer between min and max
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// RandomString generates a random string of length n
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

// RandomOwner generates a random owner name
func RandomOwner() string {
	return RandomString(6)
}

// RandomEmail generates a random email
func RandomEmail() string {
	return fmt.Sprintf("%s@email.com", RandomString(6))
}

// RandomMoney generates a random amount of money
func RandomMoney() int64 {
	return RandomInt(0, 1000)
}

// RandomCurrency generates a random currency code
func RandomCurrency() string {
	n := len(currencies)
	return currencies[rand.Intn(n)]
}

// RandomPair generates a random currency code
func RandomPair() string {
	n := len(pairs)
	return pairs[rand.Intn(n)]
}

// RandomStatus generates a random status code
func RandomStatus() string {
	n := len(status)
	return status[rand.Intn(n)]
}