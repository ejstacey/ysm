/*
Joyrex YSM - Manager for Youtube Subscriptions
Copyright (C) 2025 Eric Stacey <ejstacey@joyrex.net>

This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.

*/

package main

import (
	"encoding/json"
	"fmt"
	"os"

	"gitea.joyrex.net/ejstacey/ysm/channel"
	"gitea.joyrex.net/ejstacey/ysm/tag"
	"gitea.joyrex.net/ejstacey/ysm/tui"
	"gitea.joyrex.net/ejstacey/ysm/utils"
	"github.com/adhocore/jsonc"
	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

type Settings struct {
	Refresh bool   `json:"Refresh"`
	DbFile  string `json:"DbFile"`
}

// My stuff

func loadSettings() Settings {
	j := jsonc.New()

	input, err := os.ReadFile("settings.json")
	utils.HandleError(err, "Unable to open settings.json")

	input = j.Strip(input)
	var settings Settings
	err = json.Unmarshal(input, &settings)
	utils.HandleError(err, "Unable to read settings.json")

	return settings
}

func main() {
	var channels channel.Channels
	var tags tag.Tags

	var settings = loadSettings()

	utils.InitDb(settings.DbFile)

	fmt.Printf("Loading existing DB channel entries.\n")
	channels.LoadEntriesFromDb()

	if settings.Refresh || len(channels.ById()) == 0 {
		fmt.Printf("No existing DB entries found or Refresh was set to True in settings.json. Loading a fresh list from YouTube.\n\n")
		var newChannels = channel.LoadChannelsYoutube()
		channels.CompareAndUpdateChannelsDb(newChannels)
		channels.LoadEntriesFromDb()
	}

	fmt.Printf("Loading existing DB tag entries.\n")
	tags.LoadEntriesFromDb()

	genChannels := make([]channel.Channel, 0, len(channels.ByName()))
	genTags := make([]tag.Tag, 0, len(tags.ByName()))

	for _, chanInfo := range channels.ByName() {
		genChannels = append(genChannels, chanInfo)
	}

	for _, tagInfo := range tags.ByName() {
		genTags = append(genTags, tagInfo)
	}

	// dump.Print(genTags)

	tui.StartTea(channels, tags)
}

// err := os.WriteFile("debug.log", []byte(dump.Format(tag)), 0644)
// if err != nil {
// 	panic(err)
// }
