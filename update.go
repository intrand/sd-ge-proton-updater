package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/blang/semver"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
)

func doSelfUpdate() (bool, error) {
	v := semver.MustParse(version)
	latest, found, err := selfupdate.DetectLatest("intrand/sd-ge-proton-updater")
	if err != nil {
		return false, err
	}

	if !found {
		return false, errors.New("couldn't find latest release")
	}

	if latest.Version.Equals(v) {
		fmt.Println("already at the latest version.")
		return false, nil
	}

	if latest.Version.LT(v) {
		fmt.Println("skipped.")
		fmt.Println("Your binary seems to be newer than latest stable release.")
		return false, nil
	}

	exe, err := os.Executable()
	if err != nil {
		return false, err
	}

	err = selfupdate.UpdateTo(latest.AssetURL, exe)
	if err != nil {
		return false, err
	}

	fmt.Println("success!")
	return true, nil
}
