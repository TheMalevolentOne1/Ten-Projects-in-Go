package main

/*
Main Starting Module

Main - Command-Line Arguments ONLY.
Tasks - Handles the Struct, Slices, File Handling and JSON Un/Marshaling.
*/

import (
	"fmt"            // Format Text To Display in Terminal
	"os"             // File Handling
	"strconv"        // Verify and Convert String Arg to ID Int - Variable Type
	"strings"        // String Manipulation Functions
	"tasklist/tasks" // Tasks Exported Functions
)

func toLower(x string) string {
	var newStr string = ""
	for i := 0; i < len(x); i++ {
		newStr += strings.ToLower(string(x[i]))
	}
	return newStr
}

const addCmd = "add"
const removeCmd = "remove"
const allCmd = "all"
const test = "test"

// Program Starting Point
func main() {
	// Ensure Arguments Present
	if len(os.Args) == 1 {
		fmt.Println("No Arguments Provided")
		return
	}

	if !tasks.EnsureStorageExists() {
		fmt.Println("Storage could not be created.")
		return
	}

	err := tasks.DecodeAllJSONTasks() // Decode Existing JSON Data and Populate Tasks Struct.

	if err != nil {
		fmt.Println(err)
		return
	}

	// Ensure Subcommand Exists
	subcommand := toLower(os.Args[1]) // all, add, remove

	switch subcommand {
	case test:
		fmt.Println(byte(int(1)))
		fmt.Println(byte(1))
		return

	case allCmd:
		tasks.ShowAllTasks()
		return

	case addCmd:
		if len(os.Args) < 3 {
			fmt.Println("Missing Task Name")
			return
		}

		taskName := os.Args[2]
		taskDesc := ""

		if len(os.Args) >= 4 {
			taskDesc = os.Args[3]
		}

		if err := tasks.AddTask(taskName, taskDesc); err != nil {
			fmt.Println(err)
			fmt.Println("Task Not Saved.")
			return
		}

		fmt.Println("Task Added.")

	case removeCmd:
		if len(os.Args) < 3 {
			fmt.Println("Missing Task ID")
			return
		}

		var taskId string = os.Args[2]

		if id, err := strconv.Atoi(taskId); err != nil { // uses string convert - ascii to integer to parse integer from string e.g. "3" -> 3
			fmt.Println("Task ID provided is not a number!")
			return
		} else {
			err := tasks.RemoveTask(id)

			if err != nil {
				fmt.Println(err)
			}
		}

	case "help":
		fmt.Println(`
		HELP MENU
		Key:
		<> Required [] Optional

		Commands:
		tasklist all
		tasklist add <TASK_NAME> [DESCRIPTION] if Desc Empty then Empty Body.
		tasklist `)
	default:
		fmt.Println("Argument Not Recognised")
	}

	tasks.SaveAllTasks()
}
