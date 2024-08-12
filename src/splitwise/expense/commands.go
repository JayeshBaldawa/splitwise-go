package expense

import (
	"splitwise/model"
	"strconv"
)

const (
	FloatBase  = 64
	IntBase    = 10
	IntBitSize = 64

	InvalidAmountMessage  = "Invalid amount: "
	InvalidCommandMessage = "Invalid command: "
)

// HousemateService defines the contract for housemate operations.
type HousemateService interface {
	MoveIn(housemate string) (string, error)
	MoveOut(housemate string) (string, error)
}

// TrackerService defines the contract for expense tracking operations.
type TrackerService interface {
	AddExpense(amount float64, beneficiaries []string) (string, error)
	ShowDues(housemate string) ([]string, error)
	ClearDues(from, to string, amount int64) (string, error)
}

// TerminalCmd encapsulates the command execution logic.
type TerminalCmd struct {
	HousemateService HousemateService
	TrackerService   TrackerService
}

// NewTerminalCmd creates a new instance of TerminalCmd with provided services.
func NewTerminalCmd(housemateService HousemateService, trackerService TrackerService) *TerminalCmd {
	return &TerminalCmd{
		HousemateService: housemateService,
		TrackerService:   trackerService,
	}
}

// ExecuteCommand processes the given command by invoking the appropriate service method.
func (t *TerminalCmd) ExecuteCommand(command model.Command) string {
	switch command.CommandType {
	case model.MOVE_IN:
		return t.handleMoveIn(command.Arguments[0])
	case model.MOVE_OUT:
		return t.handleMoveOut(command.Arguments[0])
	case model.SPEND:
		return t.handleSpend(command.Arguments)
	case model.CLEAR_DUES:
		return t.handleClearDues(command.Arguments)
	case model.DUES:
		return t.handleDues(command.Arguments[0])
	default:
		return InvalidCommandMessage + string(command.CommandType)
	}
}

// handleMoveIn processes the MOVE_IN command.
func (t *TerminalCmd) handleMoveIn(housemate string) string {
	result, err := t.HousemateService.MoveIn(housemate)
	return t.processResult(result, err)
}

// handleMoveOut processes the MOVE_OUT command.
func (t *TerminalCmd) handleMoveOut(housemate string) string {
	result, err := t.HousemateService.MoveOut(housemate)
	return t.processResult(result, err)
}

// handleSpend processes the SPEND command.
func (t *TerminalCmd) handleSpend(arguments []string) string {
	amount, err := strconv.ParseFloat(arguments[0], FloatBase)
	if err != nil {
		return InvalidAmountMessage + arguments[0]
	}
	beneficiaries := arguments[1:]
	result, err := t.TrackerService.AddExpense(amount, beneficiaries)
	return t.processResult(result, err)
}

// handleClearDues processes the CLEAR_DUES command.
func (t *TerminalCmd) handleClearDues(arguments []string) string {
	amount, err := strconv.ParseInt(arguments[2], IntBase, IntBitSize)
	if err != nil {
		return InvalidAmountMessage + arguments[2]
	}
	result, err := t.TrackerService.ClearDues(arguments[0], arguments[1], amount)
	return t.processResult(result, err)
}

// handleDues processes the DUES command.
func (t *TerminalCmd) handleDues(housemate string) string {
	result, err := t.TrackerService.ShowDues(housemate)
	if err != nil {
		return err.Error()
	}
	return formatDues(result)
}

// processResult formats the result or error message.
func (t *TerminalCmd) processResult(result string, err error) string {
	if err != nil {
		return err.Error()
	}
	return result
}

// formatDues formats the dues result into a newline-separated string.
func formatDues(dues []string) string {
	out := ""
	for _, d := range dues {
		out += d + "\n"
	}

	// Check if newLine exists at the end of the string
	if len(out) > 0 && out[len(out)-1] == '\n' {
		out = out[:len(out)-1]
	}
	return out
}
