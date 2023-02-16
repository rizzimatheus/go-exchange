package util

const(
	ACTIVATE="activate"
	COMPLETED="completed"
	CANCELED="canceled"
)

// IsSupportedCurrency returns true if the currency is supported
func IsSupportedStatus(status string) bool {
	switch status {
	case ACTIVATE, COMPLETED, CANCELED:
		return true
	}
	return false
}