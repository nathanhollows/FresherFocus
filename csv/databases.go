package csv

import (
	"io/ioutil"
	"strings"
)

// Return a list of CSV files in the static/databases directory
func ListDatabaseFiles() ([]string, error) {
	files, err := ioutil.ReadDir("static/databases")
	if err != nil {
		return nil, err
	}

	var databaseFiles []string
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".csv") {
			databaseFiles = append(databaseFiles, file.Name())
		}
	}

	return databaseFiles, nil
}
