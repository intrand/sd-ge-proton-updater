package main

import (
	"log"
	"os"
	"path/filepath"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	version string = "" // to be filled in by goreleaser
	commit  string = "" // to be filled in by goreleaser
	date    string = "" // to be filled in by goreleaser
	builtBy string = "" // to be filled in by goreleaser
	cmdname string = filepath.Base(os.Args[0])
)

const (
	protonGeApiUrl string      = "https://api.github.com/repos/GloriousEggroll/proton-ge-custom/releases"
	protonGeUrl    string      = "https://github.com/GloriousEggroll/proton-ge-custom"
	protonPath     string      = "/home/deck/.local/share/Steam/compatibilitytools.d"
	systemdPath    string      = "/home/deck/.config/systemd/user/sd-ge-proton-updater.service"
	elfPath        string      = "/home/deck/.sd-ge-proton-updater"
	regExecMode    os.FileMode = 0755
	dirMode        os.FileMode = 0755
	regMode        os.FileMode = 0644
	dirModeDeck    os.FileMode = 0775
	regModeDeck    os.FileMode = 0664
)

func main() {
	// perform startup tasks
	err := startup()
	if err != nil {
		log.Fatalln(err)
	}

	// parse args
	args := kingpin.MustParse(app.Parse(os.Args[1:]))

	// main decision tree
	switch args { // look for operations at the root of the command

	// version
	case cmd_version.FullCommand():
		err = showVersion()
		if err != nil {
			log.Fatalln(err)
		}

	// install
	case cmd_install.FullCommand():
		err = install()
		if err != nil {
			log.Fatalln(err)
		}

	// prune
	case cmd_prune.FullCommand():
		err = gui()
		if err != nil {
			log.Fatalln(err)
		}
	} // end args
}
