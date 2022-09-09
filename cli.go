package main

import (
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	app         = kingpin.New("sdpud", "Updates GE-Proton on the Steam Deck in the background").Author("intrand")
	cmd_version = app.Command("version", "prints version and exits")
	// cmd_install = app.Command("install", "Installs the latest release of GE-Proton, keeping n and n-1 distributions available.") // FIXME: uncomment this and delete the cmd_install below once the rolling removals are implemented
	cmd_install = app.Command("install", "Installs the latest release of GE-Proton.")
)
