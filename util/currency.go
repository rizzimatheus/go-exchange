package util

import "strings"

// Constants for all supported currencies
const (
	BRL = "BRL"
	CAD = "CAD"
	EUR = "EUR"
	JPY = "JPY"
	USD = "USD"

	BTC   = "BTC"
	ETH   = "ETH"
	MATIC = "MATIC"
	SOL   = "SOL"
	USDT  = "USDT"
)

// Constants for all supported pairs
const (
	USDT_BRL = "USDT/BRL"
	USDT_CAD = "USDT/CAD"
	USDT_EUR = "USDT/EUR"
	USDT_JPY = "USDT/JPY"
	USDT_USD = "USDT/USD"

	BTC_USDT = "BTC/USDT"
	ETH_USDT = "ETH/USDT"
	MATIC_USDT = "MATIC/USDT"
	SOL_USDT = "SOL/USDT"
	
	ETH_BTC = "ETH/BTC"
	MATIC_BTC = "MATIC/BTC"
	SOL_BTC = "SOL/BTC"
	
	MATIC_ETH = "MATIC/ETH"
	SOL_ETH = "SOL/ETH"
)

// IsSupportedCurrency returns true if the currency is supported
func IsSupportedCurrency(currency string) bool {
	switch currency {
	case BRL, CAD, EUR, JPY, USD, BTC, ETH, MATIC, SOL, USDT:
		return true
	}
	return false
}

// IsSupportedPair returns true if the pair is supported
func IsSupportedPair(pair string) bool {
	switch pair {
	case USDT_BRL, USDT_CAD, USDT_EUR, USDT_JPY, USDT_USD, BTC_USDT, ETH_USDT, MATIC_USDT, SOL_USDT, ETH_BTC, MATIC_BTC, SOL_BTC, MATIC_ETH, SOL_ETH:
		return true
	}
	return false
}

// CurrenciesFromPair returns both currencies from a given pair
func CurrenciesFromPair(pair string) (string, string) {
	currencies := strings.Split(pair, "/")
	return currencies[0], currencies[1]
}
