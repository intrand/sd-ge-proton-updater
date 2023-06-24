package main

import (
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/alecthomas/kingpin/v2"
)

var (
	version      string = "" // to be filled in by goreleaser
	commit       string = "" // to be filled in by goreleaser
	date         string = "" // to be filled in by goreleaser
	builtBy      string = "" // to be filled in by goreleaser
	cmdname      string = filepath.Base(os.Args[0])
	installed    bool   = false
	opsys        string = runtime.GOOS
	gamesInfoMap map[int]string
)

const (
	protonGeApiUrl        string      = "https://api.github.com/repos/GloriousEggroll/proton-ge-custom/releases"
	protonGeUrl           string      = "https://github.com/GloriousEggroll/proton-ge-custom"
	systemdPath           string      = "/home/deck/.config/systemd/user/sd-ge-proton-updater.service"
	elfPath               string      = "/home/deck/.sd-ge-proton-updater"
	protonPath            string      = "/home/deck/.local/share/Steam/compatibilitytools.d"
	vdfPathConfig         string      = "/home/deck/.local/share/Steam/config/config.vdf"
	vdfPathLibraryFolders string      = "/home/deck/.local/share/Steam/config/libraryfolders.vdf"
	regExecMode           os.FileMode = 0755
	dirMode               os.FileMode = 0755
	regMode               os.FileMode = 0644
	dirModeDeck           os.FileMode = 0775
	regModeDeck           os.FileMode = 0664
	// vdfPathAppInfo string = "/home/deck/.local/share/Steam/appcache/appinfo.vdf" // FIXME: VDF binary parser and patcher needed before this is useful
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

	gamesInfoMap, err = getAllSteamAppsInfo() // this should only happen one time per run
	if err != nil {
		log.Println(err) // don't bail because of this
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
	case cmd_uninstall.FullCommand():
		err = uninstall()
		if err != nil {
			log.Fatalln(err)
		}

	// update
	case cmd_update.FullCommand():
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
