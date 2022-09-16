package main

import (
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	app         = kingpin.New("sdpud", "Updates GE-Proton on the Steam Deck in the background").Author("intrand")
	cmd_version = app.Command("version", "Prints version and exits")
	cmd_install = app.Command("install", "Installs the latest release of GE-Proton")
	cmd_prune   = app.Command("prune", "Removes versions older than latest from this machine")
)
