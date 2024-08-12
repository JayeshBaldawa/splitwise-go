package expense

import (
	"splitwise/global"
	"splitwise/model"
	"strings"
	"testing"
)

func TestExecuteCommand(t *testing.T) {
	globalStorage := global.NewGlobalMapStorage()
	terminalCmd := NewTerminalCmd(NewHousemateServiceImpl(globalStorage), NewTrackerServiceImpl(globalStorage))

	tests := []struct {
		name     string
		testPlan []struct {
			command string
			output  string
		}
	}{
		{
			name: "Test Plan 1",
			testPlan: []struct {
				command string
				output  string
			}{
				{"MOVE_IN ANDY", "SUCCESS"},
				{"MOVE_IN WOODY", "SUCCESS"},
				{"MOVE_IN BO", "SUCCESS"},
				{"SPEND 6000 WOODY ANDY BO", "SUCCESS"},
				{"SPEND 6000 ANDY BO", "SUCCESS"},
				{"DUES ANDY", "BO 0\nWOODY 0"},
				{"DUES BO", "WOODY 4000\nANDY 1000"},
				{"CLEAR_DUE BO ANDY 1000", "0"},
				{"CLEAR_DUE BO WOODY 4000", "0"},
				{"MOVE_OUT ANDY", "SUCCESS"},
				{"MOVE_OUT WOODY", "SUCCESS"},
			},
		},
		{
			name: "Test Plan 2",
			testPlan: []struct {
				command string
				output  string
			}{
				{"MOVE_IN ANDY", "SUCCESS"},
				{"MOVE_IN WOODY", "SUCCESS"},
				{"MOVE_IN BO", "SUCCESS"},
				{"MOVE_IN REX", "HOUSEFUL"},
				{"SPEND 3000 ANDY WOODY BO", "SUCCESS"},
				{"SPEND 300 WOODY BO", "SUCCESS"},
				{"SPEND 300 WOODY REX", "MEMBER_NOT_FOUND"},
				{"DUES BO", "ANDY 1150\nWOODY 0"},
				{"DUES WOODY", "ANDY 850\nBO 0"},
				{"CLEAR_DUE BO ANDY 500", "650"},
				{"CLEAR_DUE BO ANDY 2500", "INCORRECT_PAYMENT"},
				{"MOVE_OUT ANDY", "FAILURE"},
				{"MOVE_OUT WOODY", "FAILURE"},
				{"MOVE_OUT BO", "FAILURE"},
				{"CLEAR_DUE BO ANDY 650", "0"},
				{"MOVE_OUT BO", "SUCCESS"},
			},
		},
	}

	for _, tt := range tests {
		// Reset the global storage for each test plan
		globalStorage.Reset()
		t.Run(tt.name, func(t *testing.T) {
			for _, test := range tt.testPlan {
				t.Run(test.command, func(t *testing.T) {
					args := strings.Fields(test.command)
					cmdModel := model.Command{
						CommandType: model.CommandType(args[0]),
						Arguments:   args[1:],
					}
					result := terminalCmd.ExecuteCommand(cmdModel)
					if result != test.output {
						t.Errorf("Expected output: %s, but got: %s for command: %s", test.output, result, test.command)
					}
				})
			}
		})
	}
}
