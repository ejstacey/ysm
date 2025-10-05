/*
Joyrex YSM - Manager for Youtube Subscriptions
Copyright (C) 2025 Eric Stacey <ejstacey@joyrex.net>

This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.

*/

package utils

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/adhocore/jsonc"
	gap "github.com/muesli/go-app-paths"
)

type GeneratorSettings struct {
	Title        string `json:"Title"`
	TemplateFile string `json:"TemplateFile"`
	OutputFile   string `json:"OutputFile"`
}

type Settings struct {
	Refresh      bool              `json:"Refresh"`
	DbFile       string            `json:"DbFile"`
	Generator    GeneratorSettings `json:"Generator"`
	BackupCopies int               `json:"BackupCopies"`
}

// My stuff

func LoadSettings() Settings {
	userScope := gap.NewScope(gap.User, "ysm")

	var settings Settings
	dbFile, err := userScope.DataPath("ysm.db")
	HandleError(err, "Could not determine user data file!")
	settings.DbFile = dbFile

	settings.Refresh = true
	settings.Generator.Title = "My Youtube Subscriptions"

	outputFile, err := userScope.DataPath("html/index.html")
	HandleError(err, "Could not determine user data path for output file!")
	settings.Generator.OutputFile = outputFile

	templateFile, err := userScope.DataPath("templates/default.tmpl")
	HandleError(err, "Could not determine system data path for template file!")
	settings.Generator.TemplateFile = templateFile

	settings.BackupCopies = 7

	settingsFile, err := userScope.ConfigPath("settings.json")
	HandleError(err, "Could not determine user config file!")
	result, err := FileDirExists(settingsFile)
	HandleError(err, "Checking for existence of user settings file.")
	if !result {
		fmt.Println("Settings file " + settingsFile + " not found. Do you need to run --install?")
		os.Exit(1)
	}

	fmt.Printf("Loading %s\n", settingsFile)
	input, err := os.ReadFile(settingsFile)
	HandleError(err, "Unable to open settings.json")

	j := jsonc.New()
	input = j.Strip(input)
	err = json.Unmarshal(input, &settings)
	HandleError(err, "Unable to read settings.json")

	return settings
}
