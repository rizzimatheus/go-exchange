package util

const(
	ACTIVE="active"
	COMPLETED="completed"
	CANCELED="canceled"
)

// IsSupportedStatus returns true if the status is supported
func IsSupportedStatus(status string) bool {
	switch status {
	case ACTIVE, COMPLETED, CANCELED:
		return true
	}
	return false
}