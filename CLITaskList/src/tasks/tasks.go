// Handles CRUD on Tasks
package tasks

/*
Handles Tasks Struct, and Task Management.
*/

import (
	"encoding/json" // For JSON handling
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

var Tasks []Task

const file = "data/tasks.json"

var dir, err = os.Getwd() // GetWorkingDirectory

var filePath string = filepath.Join(dir, file) // json data file path

var filePerm = os.FileMode(0600) // data.json file with appropriate read and write perms for the owner.

func GetNextTaskID() int {
	return len(Tasks) + 1
}

func ShowAllTasks() {
	if lenTasks := len(Tasks); lenTasks == 0 {
		fmt.Println(header)
		fmt.Println("No Tasks.")
	} else {
		fmt.Println(header)
		for i := range lenTasks {
			fmt.Printf("Task Number: %d \nTask: %s \nDescription: %s \n",
				Tasks[i].Id,
				Tasks[i].Name,
				Tasks[i].Description,
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
		err = json.Unmarshal(b, &Tasks) // Decode JSON into Tasks Struct using its Reference.
		if err != nil {
			return err
		}
	}

	return nil
}

// verifyStorageAccess ensures valid permissions on the data.json file for Read and Write access.
func verifyStorageAccess() (bool, error) {
	info, err := os.Stat(filePath)

	if err != nil {
		return false, err
	}

	if perm := info.Mode().Perm(); perm != os.FileMode(0600) {
		//panic("Insufficient File Permissions")
		// TODO: Change using Chmod? Ask user first (not prod so skip.)

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

func SaveTask(t Task) error {
	var taskEncoded []byte

	if tJSON, err := encodeJSONTask(t); err == nil {
		taskEncoded = tJSON
	} else {
		return err
	}

	if data, err := os.ReadFile(filePath); err == nil {
		// Remove Last Element of Slice Source: https://stackoverflow.com/questions/26172196/how-to-remove-the-last-element-from-a-slice
		data = data[:len(data)-1]     // Remove ]
		if data[len(data)-1] == '}' { // Ensure Previous Char is closing '}' bracket
			data = append(data, []byte(", ")...) // add ,
		}
		data = append(data, taskEncoded...)           // append JSON Task
		data = append(data, ']')                      // ] end the array.
		err := os.WriteFile(filePath, data, filePerm) // Replace old file bytes with new.
		if err != nil {
			return err
		}
	} else {
		return err
	}

	return nil
}

// TODO: ID numbers could overlap as based on size.
// - Some ID Correction function should be implemented to ensure the numbers are accurate;
// Lazier Method: Scratching that overwriting the IDs to be manually accurate 1,2,3,3,5 to 1,2,3,4,5.
// Problem: Loses integrity of the original ID... The Wrong ID could be deleted in sequential updates.
// No undo or error mitigation handling...
func fixIDNumbers() {

}

// FindID ensures that the Task ID is found and the first { position is returned.
func FindIDTask(id int) (int, error) {
	if err != nil {
		return 0, err
	}

	for i := 0; i < len(Tasks); i++ {
		if Tasks[i].Id == id {
			return i, nil
		} else {
			continue
		}
	}

	return id, nil
}

// RemoveTask ensures a task is removed from the JSON file using the first curly bracket then deleting all subsequent bytes adding a comma if necessary.
func RemoveTask(tPos int) error {
	if !doesStorageExist() {
		return nil
	} // Verify Array in File.

	data, err := os.ReadFile(filePath)

	if err != nil {
		return err
	}

	openCurlyCount := 0
	closedCurlyCount := 0

	jsonTaskPoss := []int{}

	for i := 0; i < len(data); i++ {
		if i == '{' {
			openCurlyCount++
			jsonTaskPoss = append(jsonTaskPoss, i)
		} else if i == '}' {
			if i == len(data)-1 && closedCurlyCount != openCurlyCount {
				fmt.Println("Data is missing valid '}' for '{' starts on Line: ") // TODO
				openCurlyCount--
				continue
			}

			closedCurlyCount++
		}
	}

	fmt.Println(openCurlyCount)
	fmt.Println(closedCurlyCount)
	fmt.Printf("%v", jsonTaskPoss)

	return nil
}
