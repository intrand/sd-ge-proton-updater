package main

import (
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	app = kingpin.New("Steam Deck GE-Proton Updater", "Automatically updates GE-Proton on the Steam Deck in the background").Author("intrand")

	cmd_gui = app.Command("gui", "opens the graphical user interface of this tool").Hidden().Default()

	// tool commands
	cmd_version   = app.Command("version", "Prints version and exits")
	cmd_install   = app.Command("install", "Performs installation of sd-ge-proton-updater, and set it to run automatically on Steam Deck boot")
	cmd_uninstall = app.Command("uninstall", "Removes sd-ge-proton-updater from your Steam Deck")
	cmd_update    = app.Command("update", "Updates sd-ge-proton-updater to the latest stable release on GitHub")

	// ge-proton commands
	cmd_get   = app.Command("get", "Gets the latest release of GE-Proton")
	cmd_prune = app.Command("prune", "Removes chosen GE-Proton versions from your Steam Deck")
)
