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

	// mkdir -p
	dir, _ := filepath.Split(elfPath)
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}

	cmd := exec.Command("cp", thisElf, elfPath)
	err = cmd.Run()
	if err != nil {
		return err
	}

	if err != nil {
		log.Fatal(err)
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
	err = os.MkdirAll(dir, os.ModeDir)
	if err != nil {
		return err
	}

	err = os.WriteFile(systemdPath, []byte(unit), os.ModePerm)
	if err != nil {
		return err
	}

	cmd := exec.Command("systemctl", "--user", "enable", "--now", "sd-ge-proton-updater.service")
	err = cmd.Run()
	if err != nil {
		return err
	}

	if err != nil {
		log.Fatal(err)
	}

	return err
}

func setup() (err error) {
	// this binary, but in the correct place
	err = setupThis()
	if err != nil {
		return err
	}

	// systemd
	err = setupSystemd()
	if err != nil {
		return err
	}

	return err
}
