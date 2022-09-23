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

const (
	mainWindowWidth      int                   = 600
	mainWindowHeight     int                   = 600
	mainButtonWidth      float32               = 280
	mainButtonHeight     float32               = 150
	pruneSubButtonWidth  float32               = 100
	pruneSubButtonHeight float32               = 50
	mainWindowFlags      giu.MasterWindowFlags = giu.MasterWindowFlagsNotResizable
	// mainWindowFlags      giu.MasterWindowFlags = giu.MasterWindowFlagsFloating + giu.MasterWindowFlagsNotResizable
)

var (
	protons         []Proton
	rows            []*giu.TableRowWidget
	showPruneWindow bool = false
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

func stub() {}

func togglePruneWindow() {
	giu.OpenPopup("prune")
}

func popupError(err error) {
	giu.Msgbox("Error", err.Error())
}

func doInstall() {
	err := install()
	if err != nil {
		popupError(err)
	}
}

func userMustRelaunch(result giu.DialogResult) {
	os.Exit(0)
}

func doUpdate() {
	res, err := doSelfUpdate()
	if err != nil {
		popupError(err)
	}

	if res {
		giu.Msgbox("Info", "Update successful! You MUST close Steam Deck GE-Proton Updater and run it again to get the benefits.").ResultCallback(userMustRelaunch)
	} else {
		giu.Msgbox("Info", "Already on the latest release of Steam Deck GE-Proton Updater.")
	}
}

func doUninstall(result giu.DialogResult) {
	if result == false {
		return
	}

	err := uninstall()
	if err != nil {
		popupError(err)
		return
	}
}

func promptUninstall() {
	giu.Msgbox("Info", "Are you positive you want to remove Steam Deck GE-Proton Updater? This will not remove GE-Proton from your Steam Deck.").Buttons(giu.MsgboxButtonsYesNo).ResultCallback(doUninstall)
}

func popupVersion() {
	giu.Msgbox("About", "Steam Deck GE-Proton Updater will run at every boot. In order, it will...\n\t1. try to update itself to the latest version,\n\t2. ensure it is configured correctly,\n\t3. attempt to get information about the latest GE-Proton release,\n\t4. install the latest release if you don't already have it.\n\n"+
		"Version: "+version+"\n"+
		"Commit: "+commit+"\n"+
		"Date: "+date+"\n"+
		"Built By: "+builtBy+"\n",
	)
}

func loop() {
	pruneTable := giu.Table().Size(
		float32(mainWindowWidth-36),   // fit scrollbar
		float32(mainWindowHeight-152), // fit other widgets
	).Columns(
		giu.TableColumn("Installed GE-Proton version").Flags(
			giu.TableColumnFlagsWidthStretch + giu.TableColumnFlagsNoDirectResize + giu.TableColumnFlagsNoResize,
		),
	).Rows(
		rows...,
	)

	prunePopup := giu.PopupModal(
		"prune",
	).Flags(
		giu.WindowFlagsAlwaysAutoResize+giu.WindowFlagsNoMove+giu.WindowFlagsNoTitleBar,
	).Layout(
		giu.Row(
			giu.Button("Refresh").Size(pruneSubButtonWidth, pruneSubButtonHeight).OnClick(buildRows),
			giu.Button("Uncheck All").Size(pruneSubButtonWidth, pruneSubButtonHeight).OnClick(uncheckAll),
			giu.Button("Check All").Size(pruneSubButtonWidth, pruneSubButtonHeight).OnClick(checkAllExceptLatest),
			giu.Label("                                     "),
			giu.Button("Prune").Size(pruneSubButtonWidth, pruneSubButtonHeight).OnClick(prune),
		),
		giu.Label("Please click the Prune button to PERMANENTLY DELETE the installations checked below."),
		pruneTable,
		giu.Button("Close").Size(pruneSubButtonWidth, pruneSubButtonHeight).OnClick(giu.CloseCurrentPopup),
	)

	aboutButton := giu.Button("About").Size(mainButtonWidth, mainButtonHeight).OnClick(popupVersion)
	updateButton := giu.Button("Check for updates").Size(mainButtonWidth, mainButtonHeight).OnClick(doUpdate)
	installButton := giu.Button("Install").Size(mainButtonWidth, mainButtonHeight).OnClick(doInstall)
	uninstallButton := giu.Button("Uninstall").Size(mainButtonWidth, mainButtonHeight).OnClick(promptUninstall)
	getProtonButton := giu.Button("Get latest GE-Proton release").Size(mainButtonWidth, mainButtonHeight) //.OnClick(stub)
	pruneProtonButton := giu.Button("Prune chosen GE-Proton releases").Size(mainButtonWidth, mainButtonHeight).OnClick(togglePruneWindow)

	// buttons we want disabled when we aren't running the real deal
	if !installed {
		uninstallButton.Disabled(true)
		getProtonButton.Disabled(true)
		pruneProtonButton.Disabled(true)
	}

	getProtonButton.Disabled(true) // not ready yet

	giu.SingleWindow().Layout(
		giu.PrepareMsgbox(),
		prunePopup,
		giu.Row(
			aboutButton,
			updateButton,
		),
		giu.Row(
			installButton,
			uninstallButton,
		),
		giu.Row(
			getProtonButton,
			pruneProtonButton,
		),
	)
}

func gui() {
	mainWindow := giu.NewMasterWindow("Steam Deck GE-Proton Updater "+version+" "+commit, mainWindowWidth, mainWindowHeight, mainWindowFlags)

	buildRows()

	mainWindow.Run(loop)
}
