package model

type CommandType string

// Command types that represent various actions.
const (
	MOVE_IN    CommandType = "MOVE_IN"
	MOVE_OUT   CommandType = "MOVE_OUT"
	SPEND      CommandType = "SPEND"
	DUES       CommandType = "DUES"
	CLEAR_DUES CommandType = "CLEAR_DUE"
)

// Command represents an action with a specific CommandType and associated arguments.
type Command struct {
	CommandType CommandType
	Arguments   []string
}

// CommandError defines a type for errors related to command execution.
type CommandError string

// Error messages related to command execution.
const (
	FAILURE CommandError = "FAILURE"
)

type CommandSuccess string

// Success messages related to command execution.
const (
	SUCCESS CommandSuccess = "SUCCESS"
)
