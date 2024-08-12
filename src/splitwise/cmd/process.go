package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"splitwise/expense"
	"splitwise/global"
	"splitwise/model"
)

var terminalCmd *expense.TerminalCmd

func init() {
	globalStorage := global.NewGlobalMapStorage()
	housemateService := expense.NewHousemateServiceImpl(globalStorage)
	trackerService := expense.NewTrackerServiceImpl(globalStorage)
	terminalCmd = expense.NewTerminalCmd(housemateService, trackerService)
}

// ProcessFile reads commands from the specified file and processes each line.
func ProcessFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("error opening the input file: %w", err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if err := processLine(scanner.Text()); err != nil {
			return fmt.Errorf("error processing line: %w", err)
		}
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading the input file: %w", err)
	}
	return nil
}

// processLine parses and executes a command from a line of input.
func processLine(line string) error {
	args := strings.Fields(line)
	if len(args) == 0 {
		return nil
	}
	commandType := model.CommandType(args[0])
	commandModel := model.Command{
		CommandType: commandType,
		Arguments:   args[1:],
	}
	result := terminalCmd.ExecuteCommand(commandModel)
	fmt.Println(result)
	return nil
}
