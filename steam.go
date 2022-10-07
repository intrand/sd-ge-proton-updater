package main

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path"
	"sort"
	"strconv"

	"github.com/andygrunwald/vdf"
)

type game struct {
	Id                int
	Name              string
	CompatibilityTool string
	Path              string
}

type folder struct {
	Path string
	Apps []int
}

type allGameInfo struct {
	Applist struct {
		Apps []struct {
			Appid int    `json:"appid"`
			Name  string `json:"name"`
		} `json:"apps"`
	} `json:"applist"`
}

// getAllSteamAppsInfo() returns a map of {appId:"gameName"}
func getAllSteamAppsInfo() (gamesInfo map[int]string, err error) {
	gamesInfo = map[int]string{} // golang is so frustrating sometimes

	resp, err := http.Get("https://api.steampowered.com/ISteamApps/GetAppList/v2")
	if err != nil {
		return gamesInfo, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return gamesInfo, err
	}

	// // start file testing
	// file, err := os.Open("./dev-data/all.json")
	// if err != nil {
	// 	return gamesInfo, err
	// }

	// body, err := io.ReadAll(file)
	// if err != nil {
	// 	return gamesInfo, err
	// }
	// // end file testing

	var infos allGameInfo
	err = json.Unmarshal(body, &infos)
	if err != nil {
		return gamesInfo, err
	}

	for _, app := range infos.Applist.Apps {
		if app.Name == "" {
			continue
		}
		gamesInfo[app.Appid] = app.Name
	}

	return gamesInfo, err
}

// getSteamGameNames() takes a list of games and adds names to them based on the map from getAllSteamAppsInfo()
func getSteamGameNames(games []game) ([]game, error) {
	var err error

	for i, game := range games {
		games[i].Name = gamesInfoMap[game.Id]
	}

	return games, err
}

// getFolders() parses libraryfolders to figure out which appIds are installed where
func getFolders(games []game) ([]game, error) {
	var err error

	// installed games
	file, err := os.Open(vdfPathLibraryFolders)
	if err != nil {
		return games, err
	}

	parser := vdf.NewParser(file)
	contents, err := parser.Parse()
	if err != nil {
		return games, err
	}

	libraryfolders := contents["libraryfolders"].(map[string]interface{}) // .libraryfolders
	for _, v := range libraryfolders {
		folderId := v.(map[string]interface{})
		apps := folderId["apps"].(map[string]interface{}) // .libraryfolders.folder_id.apps
		libraryPath := folderId["path"].(string)          // .libraryfolders.folder_id.path

		for app := range apps { // .libraryfolders.folder_id.apps
			// fmt.Println(path, app)
			appId, err := strconv.Atoi(app) // .libraryfolders.folder_id.apps.key_name (not the value)
			if err != nil {
				continue
			}

			for i, game := range games { // we're an app that's in a folder. the folder has a path
				if game.Id == appId { // do we already have a game with this id?
					games[i].Path = path.Join(libraryPath, "steamapps", "common") // set the path to the app
				}
			}
		}
	}

	return games, err
}

func getConfiguredGames() (games []game, err error) {
	// compat tool
	file, err := os.Open(vdfPathConfig)
	if err != nil {
		return games, err
	}

	parser := vdf.NewParser(file)
	contents, err := parser.Parse()
	if err != nil {
		return games, err
	}

	// fmt.Println(contents)

	// c = vdf.load(open(config_vdf_file)).get('InstallConfigStore').get('Software').get('Valve').get('Steam').get('CompatToolMapping')
	installConfigStore := contents["InstallConfigStore"].(map[string]interface{})
	software := installConfigStore["Software"].(map[string]interface{})
	valve := software["Valve"].(map[string]interface{})
	steam := valve["Steam"].(map[string]interface{})
	compatToolMapping := steam["CompatToolMapping"].(map[string]interface{})

	for k, v := range compatToolMapping {
		gameId := v.(map[string]interface{})
		compatTool := gameId["name"].(string)
		if compatTool != "" {
			id, err := strconv.Atoi(k)
			if err != nil {
				continue
			}

			games = append(games, game{
				Id:                id,
				CompatibilityTool: compatTool,
			})
		}
	}

	return games, err
}

func sortGames(games []game) []game {
	sort.Slice(games, func(i, j int) bool {
		return games[i].Id < games[j].Id
	})

	return games
}

func getSteamGames() (games []game, err error) {
	// get a list of games and their compatibility tools
	games, err = getConfiguredGames()
	if err != nil {
		return games, err
	}

	games = sortGames(games) // consistent table order

	// get the folders the games are installed into
	games, err = getFolders(games)
	if err != nil {
		return games, err
	}

	// get the names of the games
	games, err = getSteamGameNames(games)
	if err != nil {
		return games, err
	}

	return games, err
}
