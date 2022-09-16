package main

import (
	"log"

	"github.com/blang/semver"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
)

func doSelfUpdate() (result bool, err error) {
	v := semver.MustParse(version)
	latest, err := selfupdate.UpdateSelf(v, "intrand/sd-ge-proton-updater")
	if err != nil {
		return false, err
	}

	if latest.Version.Equals(v) {
		log.Println("Already running the latest version:", version)
		return false, err
	} else {
		log.Println("Successfully updated to latest version: ", latest.Version)
		log.Println("Release notes:\n", latest.ReleaseNotes)
		return true, err
	}

	// return true, err
}
