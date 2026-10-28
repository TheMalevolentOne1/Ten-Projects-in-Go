// Handles CRUD on Tasks
package tasks

/*
Handles Tasks Struct, and Task Management.
*/

import (
	"encoding/json" // For JSON handling
	"errors"        // For Error Handling
	"fmt"           // Format for console display
	"os"            // OS for arguments
	"path/filepath" // For Filepath handling
)

type Task struct {
	Id          int
	Name        string
	Description string
}

const header string = `
|--------------------|
| TASKLIST MANAGER   |
|--------------------|`

var tasks []Task

const file = "data/tasks.json"

var dir, err = os.Getwd() // GetWorkingDirectory

var filePath string = filepath.Join(dir, file) // json data file path

var filePerm = os.FileMode(0600) // tasks.json file with appropriate read and write perms for the owner.

func getNextTaskID() int {
	return len(tasks) - 1 // To Match Indexes
}

func ShowAllTasks() {
	if lenTasks := len(tasks); lenTasks == 0 {
		fmt.Println(header)
		fmt.Println("No Tasks.")
	} else {
		fmt.Println(header)
		for i := range lenTasks {
			fmt.Printf("Task Number: %d \nTask: %s \nDescription: %s \n",
				tasks[i].Id,
				tasks[i].Name,
				tasks[i].Description,
			)
		}
	}
}

/*
encodeJSONTask ensures that the Task Struct t is encoded into json (returned as []byte) and otherwise an error is returned.
Parameters:
b - bytes (the bytes of the Tasks slice)
*/
func encodeJSONTask(t Task) ([]byte, error) {
	data, err := json.Marshal(t)

	if err != nil {
		return nil, err
	} else {
		return data, nil
	}
}

/*
DecodeAllJsonTasks ensures that all JSON Tasks decoded from the JSON File into the Tasks Struct otherwise returns error.
Parameters:
b - bytes (the bytes of JSON data)
*/
func DecodeAllJSONTasks() error { // Load JSON File into Tasks Struct Memory
	if b, err := os.ReadFile(filePath); err == nil {
		err = json.Unmarshal(b, &tasks) // Decode JSON into Tasks Struct using its Reference.
		if err != nil {
			return err
		}
	}

	return nil
}

// verifyStorageAccess ensures valid permissions on the tasks.json file for Read and Write access.
func verifyStorageAccess() (bool, error) {
	info, err := os.Stat(filePath)

	if err != nil {
		return false, err
	}

	if perm := info.Mode().Perm(); perm != os.FileMode(0600) {
		//panic("Insufficient File Permissions")
		// TODO: Change using Chmod? Ask user first? (not prod so skip.)

		err := os.Chmod(filePath, filePerm) // Attempt to change file permissions.

		if err != nil {
			return false, err
		}

		return true, nil
	} else {
		return true, nil
	}
}

// doesStorageExist ensures json file exists with [ bytes at the start and end ]
// returns boolean with result if exists, with valid empty slice.
func doesStorageExist() bool {
	if hasAccess, err := verifyStorageAccess(); hasAccess == false && err != nil {
		fmt.Println(err)
		return false
	}
	if _, err := os.Stat(filePath); err == nil {
		if data, err := os.ReadFile(filePath); err == nil {
			// confirm length contains at least two and [ and ]
			if len(data) >= 2 && data[0] == '[' && data[len(data)-1] == ']' {
				return true
			}
		}
	}

	return false
}

/*
EnsureStorageExists ensures the storage exists, and if it doesn't create it.
Returnn Boolean whether creation was successful or already exists.
*/
func EnsureStorageExists() bool {
	// if storage exists but doesn't contain the 91 and 93 bytes at start and end it will overwrite.
	if !doesStorageExist() {
		fmt.Println(header)
		os.Mkdir("data", filePerm) // create data folder, with 0777 perms for simplicity

		// create a tasks.json file in data directory
		// file contents stored in an empty array in bytes
		err = os.WriteFile(filePath, []byte("[]"), filePerm)

		if err != nil {
			return false
		} else {
			return true
		}
	} else {
		return true
	}
}

// SaveTask ensures that a single Task is encoded into JSON and saved into the tasks.json file.
func AddTask(tName string, tDesc string) error {
	// Add New Task to Tasks
	newTask := Task{
		Id:          getNextTaskID(),
		Name:        tName,
		Description: tDesc,
	}

	tasks = append(tasks, newTask)

	return nil
}

/*
FindIDPos ensures that the Task ID is found and the first { position is returned.
Parameters
*/
func doesIdExists(id int) (tPos int, err error) {
	checkIDExists := false

	for i := range tasks {
		if checkIDExists {
			break
		} else {
			if tasks[i].Id == id {
				checkIDExists = true
				tPos = i
			}
		}
	}

	if !checkIDExists {
		return 0, errors.New("ID cannot be found.")
	} else {
		return tPos, nil
	}
}

// RemoveTask ensures a task is removed from the Tasks Struct Slice using the Task ID.
func RemoveTask(tID int) error {
	if !doesStorageExist() {
		return errors.New("Storage does not exist.")
	}

	tPos, err := doesIdExists(tID)

	if err != nil {
		return err
	}

	if len(tasks) <= 1 {
		tasks = tasks[:0]
	} else {
		tasks = append(tasks[:tPos], tasks[(tPos+1):]...)
	}

	fmt.Println(tasks)

	return nil
}

func SaveAllTasks() (err error) {
	encodedTasks, err := json.Marshal(&tasks)

	if err != nil {
		return err
	}

	fmt.Println(tasks)

	err2 := os.WriteFile(filePath, encodedTasks, filePerm)

	if err2 != nil {
		return err2
	} else {
		return nil
	}
}
