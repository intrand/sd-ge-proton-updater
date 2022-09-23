package main

import (
	"log"
	"os"
)

func install() (err error) {
	// Handle updating to a new version
	log.Print("Attempting update of " + cmdname + "...")
	// updated, err := doSelfUpdate()
	_, err = doSelfUpdate()
	if err != nil {
		log.Println("Couldn't update at this time. Continuing. Here's what happened: " + err.Error())
	}

	// ensure proper setup
	err = setup()
	if err != nil {
		return err
	}

	return err
}

func uninstall() (err error) {
	err = os.Remove(systemdPath)
	if err != nil {
		return err
	}

	err = daemonReload()
	if err != nil {
		return err
	}

	err = os.Remove(elfPath)
	if err != nil {
		return err
	}

	return err
}
