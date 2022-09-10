package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/google/go-github/v47/github"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	version = "" // to be filled in by goreleaser
	commit  = "" // to be filled in by goreleaser
	date    = "" // to be filled in by goreleaser
	builtBy = "" // to be filled in by goreleaser
	cmdname = filepath.Base(os.Args[0])
)

const (
	protonGeApiUrl = "https://api.github.com/repos/GloriousEggroll/proton-ge-custom/releases"
	protonGeUrl    = "https://github.com/GloriousEggroll/proton-ge-custom"
	protonPath     = "/home/deck/.steam/root/compatibilitytools.d/"
	systemdPath    = "/home/deck/.config/systemd/user/sd-ge-proton-updater.service"
	elfPath        = "/home/deck/.sd-ge-proton-updater"
)

type Version struct {
	Version string `json:"version"`
	Commit  string `json:"commit"`
	Date    string `json:"date"`
	BuiltBy string `json:"builtBy"`
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
		err = os.MkdirAll(dir, os.ModeDir)
		if err != nil {
			return dir, err
		}
	}

	return dir, err // return dir for future use
}

func main() {
	// operating system check - linux only
	opsys := runtime.GOOS
	switch opsys {
	case "linux":
	default:
		log.Println(opsys + " is not supported. Exiting.")
		os.Exit(1)
	}

	// // Handle updating to a new version
	// log.Print("Attempting update of " + cmdname + "...")
	// update_result, err := doSelfUpdate()
	// if err != nil {
	// 	log.Fatalln("Couldn't update at this time. Please try again later. Exiting.")
	// }
	// if update_result {
	// 	log.Println("Please run " + cmdname + " again.")
	// 	os.Exit(0)
	// }

	// ensure proper setup
	err := setup()
	if err != nil {
		log.Fatalln(err)
	}

	// parse args
	args := kingpin.MustParse(app.Parse(os.Args[1:]))

	// main decision tree
	switch args { // look for operations at the root of the command

	// version
	case cmd_version.FullCommand():
		versionOutput := Version{
			Version: version,
			Commit:  commit,
			Date:    date,
			BuiltBy: builtBy,
		}
		versionBytes, err := json.Marshal(versionOutput)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println(string(versionBytes))
		// fmt.Println(
		// 	"{\"version\":\"" + version + "\",\"commit\":\"" + commit + "\",\"date\":\"" + date + "\",\"built_by\":\"" + builtBy + "\"}")

	// update
	case cmd_install.FullCommand():
		ctx := context.Background()             // create context
		releases, err := getLatestReleases(ctx) // get latest stable releases
		if err != nil {
			log.Fatalln(err)
		}

		latestRelease := releases[0]
		// latestMinusOneRelease := releases[1]

		latestPath := filepath.Join(protonPath, *latestRelease.TagName)

		exist, err := exists(latestPath)
		if err != nil {
			log.Fatalln(err)
		}

		if exist {
			log.Println("Release " + *latestRelease.TagName + " already exists on this console. Nothing to do. Exiting.")
			os.Exit(0)
		} // end exists

		var shaAsset *github.ReleaseAsset
		var tarballAsset *github.ReleaseAsset
		for _, asset := range latestRelease.Assets {
			if strings.Contains(*asset.Name, ".sha512sum") {
				shaAsset = asset
			}
			if strings.Contains(*asset.Name, ".tar.gz") {
				tarballAsset = asset
			}
		} // end list of assets

		var naked *github.ReleaseAsset
		if shaAsset == naked || tarballAsset == naked { // check we got data for both
			log.Fatalln("couldn't get enough info about releases. Did something change?")
		}

		dir, err := mkTempDir(*latestRelease.TagName)
		if err != nil {
			log.Fatalln(err)
		}

		var shaPath string = filepath.Join(dir, *shaAsset.Name)
		var tarballPath string = filepath.Join(dir, *tarballAsset.Name)

		// download SHA-512 checksum
		err = downloadAsset(ctx, shaAsset, shaPath)
		if err != nil {
			log.Fatalln(err)
		}

		// download GE-Proton tarball distribution
		err = downloadAsset(ctx, tarballAsset, tarballPath)
		if err != nil {
			log.Fatalln(err)
		}

		// verify SHA of tarball against SHA-512 checksum file (which is not signed!)
		err = verifySha(shaPath, tarballPath)
		if err != nil {
			log.Fatalln(err)
		}

		err = installTarGzAsset(tarballPath, protonPath)
		if err != nil {
			log.Fatalln(err)
		}

		log.Println("Successfully installed: " + *latestRelease.TagName)
	} // end args
}
