package repos

import (
	"os"
	"fmt"
)

type NumberRepo struct {}

func NewNumberRepo() NumberRepo {
	//initialize our database :P
	createLog()

	return NumberRepo{}
}

const path = "numbers.log"

/**
 Read from our values stream and write unique values to our numbers.log
 @todo validate the value is numeric and strip any leading 0s
 */
func (r NumberRepo) Save(value string) error {
	var file, err = os.OpenFile(path, os.O_RDWR, 0644)

	if isError(err) { return err }
	defer file.Close()

	// write some text line-by-line to file
	_, err = file.WriteString(value)
	if isError(err) { return err }

	// save changes
	serr := file.Sync()
	if isError(err) { return err}

	return serr
}

func createLog() error{
	// detect if file exists
	var _, err = os.Stat(path)

	if os.IsExist(err) {
		deleteLog()
	}

	var file, cerr = os.Create(path)
	if isError(cerr) { return cerr }
	return file.Close()
}

func deleteLog() {
	// delete file
	var err = os.Remove(path)
	if isError(err) { return }
}


func isError(err error) bool {
	if err != nil {
		fmt.Println(err.Error())
	}

	return err != nil
}
