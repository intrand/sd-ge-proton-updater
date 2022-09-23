package main

import (
	"log"
	"os"
	"path/filepath"
	"runtime"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	version    string = "" // to be filled in by goreleaser
	commit     string = "" // to be filled in by goreleaser
	date       string = "" // to be filled in by goreleaser
	builtBy    string = "" // to be filled in by goreleaser
	cmdname    string = filepath.Base(os.Args[0])
	installed  bool   = false
	protonPath string = "/home/deck/.local/share/Steam/compatibilitytools.d"
	opsys      string = runtime.GOOS
)

const (
	protonGeApiUrl string      = "https://api.github.com/repos/GloriousEggroll/proton-ge-custom/releases"
	protonGeUrl    string      = "https://github.com/GloriousEggroll/proton-ge-custom"
	systemdPath    string      = "/home/deck/.config/systemd/user/sd-ge-proton-updater.service"
	elfPath        string      = "/home/deck/.sd-ge-proton-updater"
	regExecMode    os.FileMode = 0755
	dirMode        os.FileMode = 0755
	regMode        os.FileMode = 0644
	dirModeDeck    os.FileMode = 0775
	regModeDeck    os.FileMode = 0664
)

func main() {
	var err error

	// parse args
	args := kingpin.MustParse(app.Parse(os.Args[1:]))

	installed, err = isInstalled()
	if err != nil {
		log.Fatalln(err)
	}

	// operating system check - linux only
	switch opsys {
	case "linux":
	default:
		log.Fatalln(opsys + " is not supported")
	}

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

	// uninstall
	case cmd_get.FullCommand():
		err = uninstall()
		if err != nil {
			log.Fatalln(err)
		}

	// update
	case cmd_get.FullCommand():
		_, err = doSelfUpdate()
		if err != nil {
			log.Fatalln(err)
		}

	// get
	case cmd_get.FullCommand():
		err = get()
		if err != nil {
			log.Fatalln(err)
		}

	// gui
	case cmd_gui.FullCommand():
		gui()
	} // end args
}
