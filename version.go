package main

import (
	"encoding/json"
	"fmt"
)

type Version struct {
	Version string `json:"version"`
	Commit  string `json:"commit"`
	Date    string `json:"date"`
	BuiltBy string `json:"builtBy"`
}

func showVersion() (err error) {
	versionOutput := Version{
		Version: version,
		Commit:  commit,
		Date:    date,
		BuiltBy: builtBy,
	}

	versionBytes, err := json.Marshal(versionOutput)
	if err != nil {
		return err
	}

	fmt.Println(string(versionBytes))

	return err
}
