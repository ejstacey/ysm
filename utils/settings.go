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
	"os"

	"github.com/adhocore/jsonc"
)

type GeneratorSettings struct {
	Title        string `json:"Title"`
	TemplateFile string `json:"TemplateFile"`
	OutputFile   string `json:"OutputFile"`
}

type Settings struct {
	Refresh   bool              `json:"Refresh"`
	DbFile    string            `json:"DbFile"`
	Generator GeneratorSettings `json:"Generator"`
}

// My stuff

func LoadSettings() Settings {
	var settings Settings
	settings.DbFile = "ysm.db"
	settings.Refresh = false
	settings.Generator.Title = "My Youtube Subscriptions"
	settings.Generator.OutputFile = "html/index.html"
	settings.Generator.TemplateFile = "templates/default.tmpl"

	input, err := os.ReadFile("settings.json")
	HandleError(err, "Unable to open settings.json")

	j := jsonc.New()
	input = j.Strip(input)
	err = json.Unmarshal(input, &settings)
	HandleError(err, "Unable to read settings.json")

	return settings
}
