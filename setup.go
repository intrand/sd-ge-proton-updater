package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func setupThis() (err error) {
	thisElf, err := os.Executable()
	if err != nil {
		return err
	}

	if thisElf == elfPath { // avoid copying self over self
		return err
	}

	// mkdir -p
	dir, _ := filepath.Split(elfPath) // get directory of ELF

	exist, err := exists(dir) // check if it exists already
	if err != nil {
		return err
	}

	if !exist { // create it if it doesn't
		err = os.MkdirAll(dir, 0700)
		if err != nil {
			return err
		}
	}

	cmd := exec.Command("/usr/bin/cp", "--force", thisElf, elfPath)
	err = cmd.Run()
	if err != nil {
		return err
	}

	err = os.Chmod(elfPath, 0700)
	if err != nil {
		return err
	}

	return err
}

func setupSystemd() (err error) {
	unit := `[Unit]
Description=sd-ge-proton-updater
After=network.target

[Service]
ExecStart=/home/deck/.sd-ge-proton-updater install
RemainAfterExit=true
Type=oneshot
Restart=on-failure
RestartSec=10

[Install]
WantedBy=default.target
`

	// mkdir -p
	dir, _ := filepath.Split(systemdPath)
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}

	err = os.WriteFile(systemdPath, []byte(unit), 0644)
	if err != nil {
		return err
	}

	cmd := exec.Command("systemctl", "--user", "daemon-reload")
	err = cmd.Run()
	if err != nil {
		return err
	}

	// handle recursive nature of self-setup
	thisElf, err := os.Executable()
	if err != nil {
		return err
	}

	if thisElf == elfPath { // avoid recursive starting
		return err
	}

	cmd = exec.Command("systemctl", "--user", "enable", "--now", "sd-ge-proton-updater.service")
	err = cmd.Run()
	if err != nil {
		return err
	}

	return err
}

func setup() (err error) {
	// this binary, but in the correct place
	log.Println("setting up binary...")
	err = setupThis()
	if err != nil {
		return err
	}

	// systemd
	log.Println("setting up systemd...")
	err = setupSystemd()
	if err != nil {
		return err
	}

	return err
}
