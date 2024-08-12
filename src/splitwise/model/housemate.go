package model

// HousemateError is a custom error type for housemate-related errors.
type HousemateError string

// Constants for error messages.
const (
	MAX_HOUSEMATES        = 3
	ZERO_DUE              = 0
	MEMBER_ALREADY_EXISTS = HousemateError("MEMBER_ALREADY_EXISTS")
	MEMBER_NOT_FOUND      = HousemateError("MEMBER_NOT_FOUND")
	HOUSEFUL              = HousemateError("HOUSEFUL")
)
