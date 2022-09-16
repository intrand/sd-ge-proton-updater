package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/AllenDang/giu"
)

type Proton struct {
	Widget  giu.Widget
	Checked bool
	Name    string
	Major   int
	Minor   int
}

var (
	masterFlags giu.MasterWindowFlags = giu.MasterWindowFlagsFloating + giu.MasterWindowFlagsNotResizable
	protons     []Proton
	rows        []*giu.TableRowWidget
)

func getInstalledProtons() (err error) {
	protons = []Proton{}
	fileNames, err := os.ReadDir(protonPath)
	if err != nil {
		return err
	}

	for _, fileName := range fileNames {
		if fileName.Type().IsDir() {

			name := fileName.Name() // simplify

			if !strings.Contains(name, "GE-Proton") {
				continue
			}

			trimmed := strings.TrimPrefix(name, "GE-Proton") // leaves 7-31, for example

			split := strings.Split(trimmed, "-") // should result in []string{"7", "31"}, for example
			if len(split) < 2 {                  // skip if we can't parse
				log.Println("couldn't determine version information from: " + name)
				continue
			}

			major, err := strconv.Atoi(split[0]) // convert major number to int
			if err != nil {
				log.Println("problem converting major: " + err.Error())
				continue
			}
			minor, err := strconv.Atoi(split[1]) // same with minor
			if err != nil {
				log.Println("problem converting minor: " + err.Error())
				continue
			}

			protons = append(protons, Proton{ // create Proton object
				Name:    name,
				Checked: false,
				Major:   major,
				Minor:   minor,
			})
		}

		for i, proton := range protons {
			proton.Widget = giu.Checkbox(proton.Name, &protons[i].Checked)
		}

	}

	return err
}

func getLatestInstalled() (latest Proton) {
	for _, proton := range protons {
		if latest.Major < proton.Major {
			latest = proton
			continue
		}
		if latest.Minor < proton.Minor {
			latest = proton
			continue
		}
	}

	return latest
}

func uncheckLatest() []Proton {
	latest := getLatestInstalled()
	for i, proton := range protons {
		if proton.Major == latest.Major && proton.Minor == latest.Minor {
			protons[i].Checked = false
		}
	}

	return protons
}

func buildRows() {
	err := getInstalledProtons()
	if err != nil {
		return
	}

	protons = uncheckLatest()

	protonBoxes := []giu.Widget{}
	for i, proton := range protons {
		protonBoxes = append(protonBoxes, giu.Checkbox(proton.Name, &protons[i].Checked))
	}

	rows = make([]*giu.TableRowWidget, len(protonBoxes))
	for i := range rows {
		rows[i] = giu.TableRow(
			protonBoxes[i],
		)
	}
}

func prune() {
	for _, proton := range protons {
		if !proton.Checked {
			continue
		}

		fmt.Println("Pruning " + proton.Name + "...")
		err := os.RemoveAll(path.Join(protonPath, proton.Name))
		if err != nil {
			log.Println("error removing " + proton.Name + ": " + err.Error())
			continue
		}
	}

	buildRows()
}

func uncheckAll() {
	for i := range protons {
		protons[i].Checked = false
	}
}

func checkAllExceptLatest() {
	latest := getLatestInstalled()
	for i, proton := range protons {
		if proton == latest {
			protons[i].Checked = false
		} else {
			protons[i].Checked = true
		}
	}

}

func loop() {
	giu.SingleWindow().Layout(
		giu.Label("Please click the Prune button to PERMANENTLY DELETE the installations checked below.").Wrapped(true),
		giu.Row(
			giu.Button("Refresh").OnClick(buildRows).Size(100, 50),
			giu.Button("Uncheck All").OnClick(uncheckAll).Size(100, 50),
			giu.Button("Check All").OnClick(checkAllExceptLatest).Size(100, 50),
			giu.Label("                                     "),
			giu.Button("Prune").OnClick(prune).Size(100, 50),
		),
		giu.Table().Rows(rows...),
	)
}

func gui() (err error) {
	masterWindow := giu.NewMasterWindow("Steam Deck GE-Proton Updater", 600, 600, masterFlags)

	buildRows()

	masterWindow.Run(loop)

	return err
}
