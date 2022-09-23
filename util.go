package main

import (
	"os"
	"os/exec"
	"path/filepath"
)

func daemonReload() (err error) {
	cmd := exec.Command("systemctl", "--user", "daemon-reload")
	err = cmd.Run()
	if err != nil {
		return err
	}

	return err
}

func isInstalled() (installed bool, err error) {
	thisElf, err := os.Executable()
	if err != nil {
		return false, err
	}

	if thisElf == elfPath { // already installed
		return true, err
	}

	return false, err // assume not installed
}

// exists returns whether the given file or directory exists
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func mkTempDir(tagName string) (dir string, err error) {
	dir = os.TempDir() // get tmp dir (usually /tmp)
	if err != nil {
		return dir, err
	}
	dir = filepath.Join(dir, "sd-ge-proton-updater", tagName) // set our custom dir

	exist, err := exists(dir) // check if it exists already
	if err != nil {
		return dir, err
	}

	if !exist { // create it if it doesn't
		err = os.MkdirAll(dir, dirMode)
		if err != nil {
			return dir, err
		}
	}

	return dir, err // return dir for future use
}
