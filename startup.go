package main

import (
	"errors"
	"log"
	"runtime"
)

func startup() (err error) {
	// operating system check - linux only
	opsys := runtime.GOOS
	switch opsys {
	case "linux":
	default:
		return errors.New(opsys + " is not supported")
	}

	// Handle updating to a new version
	log.Print("Attempting update of " + cmdname + "...")
	// updated, err := doSelfUpdate()
	_, err = doSelfUpdate()
	if err != nil {
		log.Println("Couldn't update at this time. Continuing. Here's what happened: " + err.Error())
	}

	// if updated {
	// 	log.Println()
	// }

	// ensure proper setup
	err = setup()
	if err != nil {
		return err
	}

	return err
}
